package diff

import (
	"fmt"
	"os"
	"text/tabwriter"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

/*
Difference is the most basic diff structure, it represents a difference between two apps.
The object carries information abaout what field differs from one app to another,
the value of that field on the original app and the value of that field on the current app.
*/
type Difference struct {
	Field string `json:"field"`
	From  string `json:"from"`
	To    string `json:"curr"`
}

/*
Change encapsulates all differences between two apps and carries the
information about the context those apps exist in the app tree.
*/
type Change struct {
	Context string       `json:"context"`
	Diff    []Difference `json:"diff"`
}

//Changelog log of all changes between two app trees.
type Changelog []Change

//Diff returns the changelog betwen two app trees.
func Diff(appOrig *meta.App, appCurr *meta.App) (Changelog, error) {
	var err error
	cl := Changelog{}
	cl, err = cl.diff(appOrig, appCurr, "")
	return cl, err
}

//Print is an auxiliar method used for displaying a Changelog
func (cl Changelog) Print() {
	var w *tabwriter.Writer

	for _, change := range cl {
		fmt.Println("On: ", change.Context)
		w = tabwriter.NewWriter(os.Stdout, 12, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintf(w, "Field\t From\t To\n")
		for _, diff := range change.Diff {
			fmt.Fprintf(w, "%s\t %s\t %s\n",
				diff.Field, diff.From, diff.To)
		}
		w.Flush()
	}
}

func (cl Changelog) diff(appOrig *meta.App, appCurr *meta.App, ctx string) (Changelog, error) {
	if ctx == "" {
		ctx = "*"
	}

	change := Change{
		Context: ctx,
	}

	err := change.diffMetadata(appOrig.Meta, appCurr.Meta, "")
	if err != nil {
		return Changelog{}, err
	}

	if appCurr.Meta.Name != "" {
		change.Context = fmt.Sprintf("%s.%s", change.Context, appCurr.Meta.Name)
	}

	err = change.diffAppSpec(appOrig.Spec, appCurr.Spec)
	if err != nil {
		return Changelog{}, err
	}

	if len(change.Diff) > 0 {
		cl = append(cl, change)
	}

	set := utils.AppIntersecSet(appOrig.Spec.Apps, appCurr.Spec.Apps)
	for k := range set {
		newOrig := appOrig.Spec.Apps[k]
		newCurr := appCurr.Spec.Apps[k]

		cl, err = cl.diff(newOrig, newCurr, change.Context+".Spec.Apps")
		if err != nil {
			return Changelog{}, err
		}
	}

	return cl, nil
}

func (change *Change) diffAppSpec(specOrig meta.AppSpec, specCurr meta.AppSpec) error {
	change.diffApps(specOrig.Apps, specCurr.Apps)

	err := change.diffChannels(specOrig.Channels, specCurr.Channels)
	if err != nil {
		return err
	}

	err = change.diffChannelTypes(specOrig.ChannelTypes, specCurr.ChannelTypes)
	if err != nil {
		return err
	}

	change.diffBoudaries(specOrig.Boundary, specCurr.Boundary)

	return nil
}

func (change *Change) diffBoudaries(boundOrig meta.AppBoundary, boundCurr meta.AppBoundary) {
	var orig string
	var curr string
	inputSet := utils.ArrDisjuncSet(boundOrig.Input, boundCurr.Input)
	inputOrig := utils.ArrMakeSet(boundOrig.Input)
	for k := range inputSet {
		orig = "<nil>"
		curr = "<nil>"

		if inputOrig[k] {
			orig = k
		} else {
			curr = k
		}

		change.Diff = append(change.Diff, Difference{
			Field: "Spec.Boundary.Input",
			From:  orig,
			To:    curr,
		})
	}

	outputSet := utils.ArrDisjuncSet(boundOrig.Output, boundCurr.Output)
	outputOrig := utils.ArrMakeSet(boundOrig.Output)
	for k := range outputSet {
		orig = "<nil>"
		curr = "<nil>"

		if outputOrig[k] {
			orig = k
		} else {
			curr = k
		}

		change.Diff = append(change.Diff, Difference{
			Field: "Spec.Boundary.Output",
			From:  orig,
			To:    curr,
		})
	}

}

func (change *Change) diffApps(appsOrig utils.Apps, appsCurr utils.Apps) {
	set := utils.AppDisjuncSet(appsOrig, appsCurr)

	for k := range set {
		_, orig := appsOrig[k]

		origAppStatus := "<nil>"
		currAppStatus := "<nil>"

		if orig {
			origAppStatus = "{...}"
		} else {

			currAppStatus = "{...}"
		}

		change.Diff = append(change.Diff, Difference{
			Field: fmt.Sprintf("Spec.Apps[%s]", k),
			From:  origAppStatus,
			To:    currAppStatus,
		})
	}
}

func (change *Change) diffChannels(chOrig utils.Channels, chCurr utils.Channels) error {
	disjunction := utils.ChsDisjuncSet(chOrig, chCurr)

	for k := range disjunction {
		_, orig := chOrig[k]
		origChStatus := "<nil>"
		currChStatus := "<nil>"

		if orig {
			origChStatus = "{...}"
		} else {
			currChStatus = "{...}"
		}

		change.Diff = append(change.Diff, Difference{
			Field: fmt.Sprintf("Spec.Channels[%s]", k),
			From:  origChStatus,
			To:    currChStatus,
		})
	}

	intersection := utils.ChsIntersecSet(chOrig, chCurr)

	for k := range intersection {
		origCh := chOrig[k]
		currCh := chCurr[k]
		if origCh.Spec.Type != currCh.Spec.Type {
			change.Diff = append(change.Diff, Difference{
				Field: fmt.Sprintf("Spec.Channels[%s].Spec.Type", k),
				From:  origCh.Spec.Type,
				To:    currCh.Spec.Type,
			})
		}

		err := change.diffMetadata(origCh.Meta, currCh.Meta, "Spec.Channels["+k+"].")
		if err != nil {
			return err
		}
	}

	return nil
}

func (change *Change) diffChannelTypes(chtOrig utils.Types, chtCurr utils.Types) error {
	disjunction := utils.TypesDisjuncSet(chtOrig, chtCurr)

	for k := range disjunction {
		_, orig := chtOrig[k]

		origChtStatus := "<nil>"
		currChtStatus := "<nil>"

		if orig {
			origChtStatus = "{...}"
		} else {
			currChtStatus = "{...}"
		}

		change.Diff = append(change.Diff, Difference{
			Field: fmt.Sprintf("Spec.ChannelTypes[%s]", k),
			From:  origChtStatus,
			To:    currChtStatus,
		})
	}

	intersection := utils.TypesIntersecSet(chtOrig, chtCurr)

	for k := range intersection {
		origCht := chtOrig[k]
		currCht := chtCurr[k]

		if string(origCht.Schema) != string(currCht.Schema) {
			change.Diff = append(change.Diff, Difference{
				Field: fmt.Sprintf("Spec.ChannelTypes[%s].Spec.Schema", k),
				From:  string(origCht.Schema),
				To:    string(currCht.Schema),
			})
		}

		err := change.diffMetadata(origCht.Meta, currCht.Meta, "Spec.ChannelTypes["+k+"].")
		if err != nil {
			return err
		}
	}

	return nil
}

func (change *Change) diffMetadata(metaOrig meta.Metadata, metaCurr meta.Metadata, ctx string) error {
	var err error
	err = nil

	if metaOrig.Name != metaCurr.Name {
		err = fmt.Errorf("on %s Metadata: Different name", ctx)
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Name",
			From:  metaOrig.Name,
			To:    metaCurr.Name,
		})
	}

	if metaOrig.Reference != metaCurr.Reference {
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Reference",
			From:  metaOrig.Reference,
			To:    metaCurr.Reference,
		})
	}

	if metaOrig.Parent != metaCurr.Parent {
		err = fmt.Errorf("on %s Metadata: Different parent", ctx)
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Parent",
			From:  metaOrig.Parent,
			To:    metaCurr.Parent,
		})
	}

	if metaOrig.SHA256 != metaCurr.SHA256 {
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.SHA256",
			From:  metaOrig.SHA256,
			To:    metaCurr.SHA256,
		})
	}

	set := utils.StrDisjuncSet(metaOrig.Annotations, metaCurr.Annotations)

	for k := range set {
		origVal := metaOrig.Annotations[k]
		currVal := metaCurr.Annotations[k]

		if origVal == "" {
			origVal = "<nil>"
		}

		if currVal == "" {
			currVal = "<nil>"
		}

		change.Diff = append(change.Diff, Difference{
			Field: fmt.Sprintf("Meta.Annotations[%s]", k),
			From:  origVal,
			To:    currVal,
		})
	}

	return err
}
