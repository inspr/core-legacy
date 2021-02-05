package diff

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
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

type void struct{}

var exists void

type apps map[string]*meta.App
type channels map[string]*meta.Channel
type types map[string]*meta.ChannelType

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

	_, err := change.diffMetadata(appOrig.Meta, appCurr.Meta, "")
	if err != nil {
		return Changelog{}, err
	}

	if appCurr.Meta.Name != "" {
		change.Context = fmt.Sprintf("%s.%s", change.Context, appCurr.Meta.Name)
	}

	_, err = change.diffAppSpec(appOrig.Spec, appCurr.Spec)
	if err != nil {
		return Changelog{}, err
	}

	if len(change.Diff) > 0 {
		cl = append(cl, change)
	}

	set := make(map[string]void)

	for k := range appOrig.Spec.Apps {
		set[k] = exists
	}

	for k := range appCurr.Spec.Apps {
		set[k] = exists
	}

	for k := range set {
		newOrig, orig := appOrig.Spec.Apps[k]
		newCurr, curr := appCurr.Spec.Apps[k]

		if orig && curr {
			cl, err = cl.diff(newOrig, newCurr, change.Context+"Spec.Apps")
			if err != nil {
				return Changelog{}, err
			}
		}
	}

	return cl, nil
}

func (change *Change) diffAppSpec(specOrig meta.AppSpec, specCurr meta.AppSpec) (bool, error) {
	specsDiffer := false
	_, err := change.diffApps(specOrig.Apps, specCurr.Apps)
	if err != nil {
		return false, err
	}

	_, err = change.diffChannels(specOrig.Channels, specCurr.Channels)
	if err != nil {
		return false, err
	}

	_, err = change.diffChannelTypes(specOrig.ChannelTypes, specCurr.ChannelTypes)
	if err != nil {
		return false, err
	}

	return specsDiffer, nil
}

func (change *Change) diffApps(appsOrig apps, appsCurr apps) (bool, error) {
	appsDiffer := false

	set := make(map[string]void)

	for k := range appsOrig {
		set[k] = exists
	}

	for k := range appsCurr {
		set[k] = exists
	}

	for k := range set {
		_, orig := appsOrig[k]
		_, curr := appsCurr[k]

		origAppStatus := "<nil>"
		currAppStatus := "<nil>"

		if orig {
			origAppStatus = "{...}"
		}

		if curr {
			currAppStatus = "{...}"
		}

		if orig != curr {
			appsDiffer = true
			change.Diff = append(change.Diff, Difference{
				Field: "Spec.Apps[" + k + "]",
				From:  origAppStatus,
				To:    currAppStatus,
			})
		}
	}

	return appsDiffer, nil
}

func (change *Change) diffChannels(chOrig channels, chCurr channels) (bool, error) {
	channelsDiffer := false

	set := make(map[string]void)

	for k := range chOrig {
		set[k] = exists
	}

	for k := range chCurr {
		set[k] = exists
	}

	for k := range set {
		origCh, orig := chOrig[k]
		currCh, curr := chCurr[k]

		if orig != curr {
			origChStatus := "<nil>"
			currChStatus := "<nil>"

			if orig {
				origChStatus = "{...}"
			}

			if curr {
				currChStatus = "{...}"
			}

			channelsDiffer = true
			change.Diff = append(change.Diff, Difference{
				Field: "Spec.Channels[" + k + "]",
				From:  origChStatus,
				To:    currChStatus,
			})

		} else if orig && curr {

			if origCh.Spec.Type != currCh.Spec.Type {
				channelsDiffer = true
				change.Diff = append(change.Diff, Difference{
					Field: "Spec.Channels[" + k + "].Spec.Type",
					From:  origCh.Spec.Type,
					To:    currCh.Spec.Type,
				})
			}

			_, err := change.diffMetadata(origCh.Meta, currCh.Meta, "Spec.Channels["+k+"].")
			if err != nil {
				return channelsDiffer, err
			}
		}
	}

	return channelsDiffer, nil
}

func (change *Change) diffChannelTypes(chtOrig types, chtCurr types) (bool, error) {
	channelsDiffer := false

	set := make(map[string]void)

	for k := range chtOrig {
		set[k] = exists
	}

	for k := range chtCurr {
		set[k] = exists
	}

	for k := range set {
		origCht, orig := chtOrig[k]
		currCht, curr := chtCurr[k]

		if orig != curr {
			origChtStatus := "<nil>"
			currChtStatus := "<nil>"

			if orig {
				origChtStatus = "{...}"
			}

			if curr {
				currChtStatus = "{...}"
			}

			channelsDiffer = true
			change.Diff = append(change.Diff, Difference{
				Field: "Spec.ChannelTypes[" + k + "]",
				From:  origChtStatus,
				To:    currChtStatus,
			})
		} else if orig && curr {
			if string(origCht.Schema) != string(currCht.Schema) {
				channelsDiffer = true
				change.Diff = append(change.Diff, Difference{
					Field: "Spec.ChannelTypes[" + k + "].Spec.Schema",
					From:  string(origCht.Schema),
					To:    string(currCht.Schema),
				})
			}

			_, err := change.diffMetadata(origCht.Meta, currCht.Meta, "Spec.ChannelTypes["+k+"].")

			if err != nil {
				return channelsDiffer, err
			}
		}
	}
	return channelsDiffer, nil
}

func (change *Change) diffMetadata(metaOrig meta.Metadata, metaCurr meta.Metadata, ctx string) (bool, error) {
	var err error
	differs := false
	err = nil
	if metaOrig.Name != metaCurr.Name {
		differs = true
		err = errors.New("Diffrent name")
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Name",
			From:  metaOrig.Name,
			To:    metaCurr.Name,
		})
	}

	if metaOrig.Reference != metaCurr.Reference {
		differs = true
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Reference",
			From:  metaOrig.Reference,
			To:    metaCurr.Reference,
		})
	}

	if metaOrig.Parent != metaCurr.Parent {
		differs = true
		err = errors.New("Diffrent parent")
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Parent",
			From:  metaOrig.Parent,
			To:    metaCurr.Parent,
		})
	}

	if metaOrig.SHA256 != metaCurr.SHA256 {
		differs = true
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.SHA256",
			From:  metaOrig.SHA256,
			To:    metaCurr.SHA256,
		})
	}

	set := make(map[string]void)

	for k := range metaOrig.Annotations {
		set[k] = exists
	}

	for k := range metaCurr.Annotations {
		set[k] = exists
	}

	for k := range set {
		origVal, orig := metaOrig.Annotations[k]
		currVal, curr := metaCurr.Annotations[k]

		if origVal == "" {
			origVal = "<nil>"
		}

		if currVal == "" {
			currVal = "<nil>"
		}

		if orig != curr {
			differs = true
			change.Diff = append(change.Diff, Difference{
				Field: ctx + "Meta.Annotations[" + k + "]",
				From:  origVal,
				To:    currVal,
			})
		}
	}

	return differs, err
}
