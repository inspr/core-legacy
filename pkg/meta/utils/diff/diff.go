package diff

import (
	"fmt"
	"os"
	"text/tabwriter"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
)

// Kind represents a kind of difference between two structures
type Kind int

// Kinds of diff
const (
	AppKind Kind = 1 << iota
	NodeKind
	MetaKind
	ChannelKind
	ChannelTypeKind
	BoundaryKind
	FieldKind
	AnnotationKind
)

// Operation represents an operation that has been applied in a diff
type Operation int

// The kinds of operation
const (
	Delete Operation = 1 << iota
	Update
	Create
)

/*
Difference is the most basic diff structure, it represents a difference between two apps.
The object carries information abaout what field differs from one app to another,
the value of that field on the original app and the value of that field on the current app.
*/
type Difference struct {
	Field     string `json:"field"`
	From      string `json:"from"`
	To        string `json:"to"`
	Kind      Kind
	Name      string
	Operation Operation
}

/*
Change encapsulates all differences between two apps and carries the
information about the context those apps exist in the app tree.
*/
type Change struct {
	Context   string       `json:"context"`
	Diff      []Difference `json:"diff"`
	Kind      Kind
	Operation Operation
}

//Changelog log of all changes between two app trees.
type Changelog []Change

//Diff returns the changelog between two app trees.
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

	err := change.diffMetadata(appOrig.Meta.Name, AppKind, appOrig.Meta, appCurr.Meta, "")
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
	err := change.diffNodes(specOrig.Node, specCurr.Node)
	if err != nil {
		return err
	}

	change.diffApps(specOrig.Apps, specCurr.Apps)

	err = change.diffChannels(specOrig.Channels, specCurr.Channels)
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

func (change *Change) diffNodes(nodeOrig meta.Node, nodeCurr meta.Node) error {
	err := change.diffMetadata(nodeOrig.Meta.Name, NodeKind, nodeOrig.Meta, nodeCurr.Meta, "Spec.Node.")
	if err != nil {
		return err
	}

	if nodeOrig.Spec.Image != nodeCurr.Spec.Image {
		change.Diff = append(change.Diff, Difference{
			Field:     "Spec.Node.Spec.Image",
			From:      nodeOrig.Spec.Image,
			To:        nodeCurr.Spec.Image,
			Kind:      NodeKind,
			Operation: Update,
		})
		change.Kind |= NodeKind
		change.Operation |= Update
	}
	return nil
}

func (change *Change) diffBoudaries(boundOrig meta.AppBoundary, boundCurr meta.AppBoundary) {
	var orig string
	var curr string
	inputSet := utils.ArrDisjuncSet(boundOrig.Input, boundCurr.Input)
	inputOrig := utils.ArrMakeSet(boundOrig.Input)
	for k := range inputSet {
		var op Operation
		orig = "<nil>"
		curr = "<nil>"

		if inputOrig[k] {
			orig = k
			op = Delete
		} else {
			curr = k
			op = Create
		}

		change.Diff = append(change.Diff, Difference{
			Field:     "Spec.Boundary.Input",
			From:      orig,
			To:        curr,
			Kind:      BoundaryKind,
			Name:      k,
			Operation: op,
		})

		change.Kind |= BoundaryKind
		change.Operation |= op
	}

	outputSet := utils.ArrDisjuncSet(boundOrig.Output, boundCurr.Output)
	outputOrig := utils.ArrMakeSet(boundOrig.Output)
	for k := range outputSet {
		var op Operation
		orig = "<nil>"
		curr = "<nil>"

		if outputOrig[k] {
			orig = k
			op = Delete
		} else {
			curr = k
			op = Create
		}

		change.Diff = append(change.Diff, Difference{
			Field:     "Spec.Boundary.Output",
			From:      orig,
			To:        curr,
			Kind:      BoundaryKind,
			Operation: op,
			Name:      k,
		})
		change.Kind |= BoundaryKind
		change.Operation |= op
	}

}

func (change *Change) diffApps(appsOrig utils.MApps, appsCurr utils.MApps) {
	set := utils.AppDisjuncSet(appsOrig, appsCurr)

	for k := range set {
		var op Operation
		_, orig := appsOrig[k]

		from := "<nil>"
		to := "<nil>"

		if orig {
			from = "{...}"
			op = Delete
		} else {
			to = "{...}"
			op = Create
		}

		change.Diff = append(change.Diff, Difference{
			Field:     fmt.Sprintf("Spec.Apps[%s]", k),
			From:      from,
			To:        to,
			Kind:      AppKind,
			Operation: op,
			Name:      k,
		})
		change.Kind |= AppKind
		change.Operation |= op
	}
}

func (change *Change) diffChannels(chOrig utils.MChannels, chCurr utils.MChannels) error {
	disjunction := utils.ChsDisjuncSet(chOrig, chCurr)

	for ch := range disjunction {
		_, orig := chOrig[ch]
		from := "<nil>"
		to := "<nil>"
		var op Operation
		if orig {
			from = "{...}"
			op = Delete
		} else {
			to = "{...}"
			op = Create
		}

		change.Diff = append(change.Diff, Difference{
			Field:     fmt.Sprintf("Spec.Channels[%s]", ch),
			From:      from,
			To:        to,
			Kind:      ChannelKind,
			Operation: op,
			Name:      ch,
		})
		change.Kind |= ChannelKind
		change.Operation |= op
	}

	intersection := utils.ChsIntersecSet(chOrig, chCurr)

	for ch := range intersection {
		origCh := chOrig[ch]
		currCh := chCurr[ch]
		if origCh.Spec.Type != currCh.Spec.Type {
			change.Diff = append(change.Diff, Difference{
				Field:     fmt.Sprintf("Spec.Channels[%s].Spec.Type", ch),
				From:      origCh.Spec.Type,
				To:        currCh.Spec.Type,
				Kind:      ChannelKind,
				Operation: Update,
				Name:      ch,
			})
			change.Kind |= ChannelKind
			change.Operation |= Update
		}

		err := change.diffMetadata(ch, ChannelKind, origCh.Meta, currCh.Meta, "Spec.Channels["+ch+"].")
		if err != nil {
			return err
		}
	}

	return nil
}

func (change *Change) diffChannelTypes(chtOrig utils.MTypes, chtCurr utils.MTypes) error {
	disjunction := utils.TypesDisjuncSet(chtOrig, chtCurr)

	for ct := range disjunction {
		_, orig := chtOrig[ct]

		from := "<nil>"
		to := "<nil>"
		var op Operation
		if orig {
			from = "{...}"
			op = Delete

		} else {
			to = "{...}"
			op = Create
		}

		change.Diff = append(change.Diff, Difference{
			Field:     fmt.Sprintf("Spec.ChannelTypes[%s]", ct),
			From:      from,
			To:        to,
			Kind:      ChannelTypeKind,
			Operation: op,
			Name:      ct,
		})
		change.Kind |= ChannelTypeKind
		change.Operation |= op
	}

	intersection := utils.TypesIntersecSet(chtOrig, chtCurr)

	for ct := range intersection {
		from := chtOrig[ct]
		to := chtCurr[ct]

		if string(from.Schema) != string(to.Schema) {
			change.Diff = append(change.Diff, Difference{
				Field:     fmt.Sprintf("Spec.ChannelTypes[%s].Spec.Schema", ct),
				From:      string(from.Schema),
				To:        string(to.Schema),
				Kind:      ChannelTypeKind,
				Operation: Update,
				Name:      ct,
			})
			change.Kind |= ChannelTypeKind
			change.Operation |= Update
		}

		err := change.diffMetadata(ct, ChannelTypeKind, from.Meta, to.Meta, fmt.Sprintf("Spec.ChannelTypes[%s].", ct))
		if err != nil {
			return err
		}

	}

	return nil
}

func (change *Change) diffMetadata(parentElement string, parentKind Kind, metaOrig meta.Metadata, metaCurr meta.Metadata, ctx string) error {
	var errs string

	if metaOrig.Name != metaCurr.Name {
		errs += fmt.Sprintf("on %s Metadata: Different name", ctx)
		change.Diff = append(change.Diff, Difference{
			Field:     ctx + "Meta.Name",
			From:      metaOrig.Name,
			To:        metaCurr.Name,
			Kind:      MetaKind | parentKind,
			Operation: Update,
			Name:      parentElement,
		})
		change.Kind |= MetaKind | parentKind
		change.Operation |= Update
	}

	if metaOrig.Reference != metaCurr.Reference {
		change.Diff = append(change.Diff, Difference{
			Field:     ctx + "Meta.Reference",
			From:      metaOrig.Reference,
			To:        metaCurr.Reference,
			Kind:      MetaKind | parentKind,
			Operation: Update,
			Name:      parentElement,
		})
		change.Kind |= MetaKind | parentKind
		change.Operation |= Update
	}

	if metaOrig.Parent != metaCurr.Parent {
		errs += fmt.Sprintf("on %s Metadata: Different parent", ctx)
		change.Diff = append(change.Diff, Difference{
			Field:     ctx + "Meta.Parent",
			From:      metaOrig.Parent,
			To:        metaCurr.Parent,
			Kind:      MetaKind | parentKind,
			Operation: Update,
			Name:      parentElement,
		})
		change.Kind |= MetaKind | parentKind
		change.Operation |= Update
	}

	if metaOrig.SHA256 != metaCurr.SHA256 {
		change.Diff = append(change.Diff, Difference{
			Field:     ctx + "Meta.SHA256",
			From:      metaOrig.SHA256,
			To:        metaCurr.SHA256,
			Kind:      MetaKind | parentKind,
			Operation: Update,
			Name:      parentElement,
		})
		change.Operation |= Update
		change.Kind |= MetaKind | parentKind
	}

	set := utils.StrDisjuncSet(metaOrig.Annotations, metaCurr.Annotations)

	for k := range set {
		var op Operation
		origVal := metaOrig.Annotations[k]
		currVal := metaCurr.Annotations[k]

		if origVal == "" {
			origVal = "<nil>"
			op = Create
		}

		if currVal == "" {
			currVal = "<nil>"
			op = Delete
		}

		change.Diff = append(change.Diff, Difference{
			Field:     fmt.Sprintf("Meta.Annotations[%s]", k),
			From:      origVal,
			To:        currVal,
			Kind:      MetaKind | parentKind | AnnotationKind,
			Name:      k,
			Operation: op,
		})
		change.Kind |= MetaKind | parentKind | AnnotationKind
		change.Operation |= op
	}

	if errs != "" {
		return ierrors.NewError().InvalidApp().Message(errs).Build()
	}
	return nil
}
