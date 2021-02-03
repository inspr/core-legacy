package diff

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type Difference struct {
	Field string
	From  string
	Curr  string
}

type Change struct {
	Context string
	Diff    []Difference
}

type Changelog []Change

type void struct{}

var exists void

func Diff(appOrgn *meta.App, appCurr *meta.App) (Changelog, error) {
	var err error
	cl := Changelog{}
	cl, err = cl.diff(appOrgn, appCurr, "")
	return cl, err
}

func (cl Changelog) Print() {
	var w *tabwriter.Writer

	for _, change := range cl {
		fmt.Println("On: ", change.Context)
		w = tabwriter.NewWriter(os.Stdout, 12, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintf(w, "Field\t From\t To\n")
		for _, diff := range change.Diff {
			fmt.Fprintf(w, "%s\t %s\t %s\n",
				diff.Field, diff.From, diff.Curr)
		}
		w.Flush()
	}
}

func (cl Changelog) diff(appOrgn *meta.App, appCurr *meta.App, ctx string) (Changelog, error) {
	change := Change{
		Context: ctx,
	}
	_, err := change.diffMetadata(appOrgn.Meta, appCurr.Meta)
	if err != nil {
		return Changelog{}, err
	}
	change.Context = change.Context + "." + appCurr.Meta.Name
	_, err = change.diffAppSpec(appOrgn.Spec, appCurr.Spec)
	if err != nil {
		return Changelog{}, err
	}
	if len(change.Diff) > 0 {
		cl = append(cl, change)
	}
	return cl, nil
}

func (change *Change) diffAppSpec(specOrgn meta.AppSpec, specCurr meta.AppSpec) (bool, error) {
	appsDiffer := false

	set := make(map[string]void)
	for k := range specOrgn.Apps {
		set[k] = exists
	}
	for k := range specCurr.Apps {
		set[k] = exists
	}
	for k := range set {
		_, orgn := specOrgn.Apps[k]
		_, curr := specCurr.Apps[k]
		orgnApp := "<nil>"
		currApp := "<nil>"
		if orgn {
			orgnApp = "{...}"
		}

		if curr {
			currApp = "{...}"
		}

		if orgn != curr {
			appsDiffer = true
			change.Diff = append(change.Diff, Difference{
				Field: "Spec.Apps[" + k + "]",
				From:  orgnApp,
				Curr:  currApp,
			})
		}
	}
	return appsDiffer, nil
}

func (change *Change) diffMetadata(metaOrgn meta.Metadata, metaCurr meta.Metadata) (bool, error) {
	var err error
	differs := false
	err = nil
	if metaOrgn.Name != metaCurr.Name {
		differs = true
		err = errors.New("Diffrent name")
		change.Diff = append(change.Diff, Difference{
			Field: "Meta.Name",
			From:  metaOrgn.Name,
			Curr:  metaCurr.Name,
		})
	}

	if metaOrgn.Reference != metaCurr.Reference {
		differs = true
		change.Diff = append(change.Diff, Difference{
			Field: "Meta.Reference",
			From:  metaOrgn.Reference,
			Curr:  metaCurr.Reference,
		})
	}

	if metaOrgn.Parent != metaCurr.Parent {
		differs = true
		err = errors.New("Diffrent parent")
		change.Diff = append(change.Diff, Difference{
			Field: "Meta.Parent",
			From:  metaOrgn.Parent,
			Curr:  metaCurr.Parent,
		})
	}

	if metaOrgn.SHA256 != metaCurr.SHA256 {
		differs = true
		change.Diff = append(change.Diff, Difference{
			Field: "Meta.SHA256",
			From:  metaOrgn.SHA256,
			Curr:  metaCurr.SHA256,
		})
	}

	set := make(map[string]void)
	for k := range metaOrgn.Annotations {
		set[k] = exists
	}
	for k := range metaCurr.Annotations {
		set[k] = exists
	}
	for k := range set {
		orgVal, org := metaOrgn.Annotations[k]
		currVal, curr := metaCurr.Annotations[k]
		if orgVal == "" {
			orgVal = "<nil>"
		}
		if currVal == "" {
			currVal = "<nil>"
		}
		if org != curr {
			differs = true
			change.Diff = append(change.Diff, Difference{
				Field: "Meta.Annotations[" + k + "]",
				From:  orgVal,
				Curr:  currVal,
			})
		}
	}

	return differs, err
}
