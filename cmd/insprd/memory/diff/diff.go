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

type apps map[string]*meta.App
type channels map[string]*meta.Channel
type types map[string]*meta.ChannelType

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
	if ctx == "" {
		ctx = "->"
	}
	change := Change{
		Context: ctx,
	}
	_, err := change.diffMetadata(appOrgn.Meta, appCurr.Meta, "")
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
	set := make(map[string]void)
	for k := range appOrgn.Spec.Apps {
		set[k] = exists
	}
	for k := range appCurr.Spec.Apps {
		set[k] = exists
	}
	for k := range set {
		newOrgn, orgn := appOrgn.Spec.Apps[k]
		newCurr, curr := appCurr.Spec.Apps[k]

		if orgn && curr {
			cl, err = cl.diff(newOrgn, newCurr, change.Context+"Spec.Apps")
			if err != nil {
				return Changelog{}, err
			}
		}
	}
	return cl, nil
}

func (change *Change) diffAppSpec(specOrgn meta.AppSpec, specCurr meta.AppSpec) (bool, error) {
	specsDiffer := false
	_, err := change.diffApps(specOrgn.Apps, specCurr.Apps)
	if err != nil {
		return false, err
	}
	_, err = change.diffChannels(specOrgn.Channels, specCurr.Channels)
	if err != nil {
		return false, err
	}
	_, err = change.diffChannelTypes(specOrgn.ChannelTypes, specCurr.ChannelTypes)
	if err != nil {
		return false, err
	}
	return specsDiffer, nil
}

func (change *Change) diffApps(appsOrgn apps, appsCurr apps) (bool, error) {
	appsDiffer := false

	set := make(map[string]void)
	for k := range appsOrgn {
		set[k] = exists
	}
	for k := range appsCurr {
		set[k] = exists
	}
	for k := range set {
		_, orgn := appsOrgn[k]
		_, curr := appsCurr[k]
		orgnAppStatus := "<nil>"
		currAppStatus := "<nil>"
		if orgn {
			orgnAppStatus = "{...}"
		}

		if curr {
			currAppStatus = "{...}"
		}

		if orgn != curr {
			appsDiffer = true
			change.Diff = append(change.Diff, Difference{
				Field: "Spec.Apps[" + k + "]",
				From:  orgnAppStatus,
				Curr:  currAppStatus,
			})
		}
	}
	return appsDiffer, nil

}

func (change *Change) diffChannels(chOrgn channels, chCurr channels) (bool, error) {
	channelsDiffer := false

	set := make(map[string]void)
	for k := range chOrgn {
		set[k] = exists
	}
	for k := range chCurr {
		set[k] = exists
	}
	for k := range set {
		orgnCh, orgn := chOrgn[k]
		currCh, curr := chCurr[k]

		if orgn != curr {
			orgnChStatus := "<nil>"
			currChStatus := "<nil>"
			if orgn {
				orgnChStatus = "{...}"
			}
			if curr {
				currChStatus = "{...}"
			}
			channelsDiffer = true
			change.Diff = append(change.Diff, Difference{
				Field: "Spec.Channels[" + k + "]",
				From:  orgnChStatus,
				Curr:  currChStatus,
			})
		} else if orgn && curr {
			if orgnCh.Spec.Type != currCh.Spec.Type {
				channelsDiffer = true
				change.Diff = append(change.Diff, Difference{
					Field: "Spec.Channels[" + k + "].Spec.Type",
					From:  orgnCh.Spec.Type,
					Curr:  currCh.Spec.Type,
				})
			}

			_, err := change.diffMetadata(orgnCh.Meta, currCh.Meta, "Spec.Channels["+k+"].")
			if err != nil {
				return channelsDiffer, err
			}

		}
	}
	return channelsDiffer, nil
}

func (change *Change) diffChannelTypes(chtOrgn types, chtCurr types) (bool, error) {
	channelsDiffer := false

	set := make(map[string]void)
	for k := range chtOrgn {
		set[k] = exists
	}
	for k := range chtCurr {
		set[k] = exists
	}
	for k := range set {
		orgnCht, orgn := chtOrgn[k]
		currCht, curr := chtCurr[k]

		if orgn != curr {
			orgnChtStatus := "<nil>"
			currChtStatus := "<nil>"
			if orgn {
				orgnChtStatus = "{...}"
			}
			if curr {
				currChtStatus = "{...}"
			}
			channelsDiffer = true
			change.Diff = append(change.Diff, Difference{
				Field: "Spec.ChannelTypes[" + k + "]",
				From:  orgnChtStatus,
				Curr:  currChtStatus,
			})
		} else if orgn && curr {
			if string(orgnCht.Schema) != string(currCht.Schema) {
				channelsDiffer = true
				change.Diff = append(change.Diff, Difference{
					Field: "Spec.ChannelTypes[" + k + "].Spec.Schema",
					From:  string(orgnCht.Schema),
					Curr:  string(currCht.Schema),
				})
			}

			_, err := change.diffMetadata(orgnCht.Meta, currCht.Meta, "Spec.ChannelTypes["+k+"].")
			if err != nil {
				return channelsDiffer, err
			}

		}
	}
	return channelsDiffer, nil
}

func (change *Change) diffMetadata(metaOrgn meta.Metadata, metaCurr meta.Metadata, ctx string) (bool, error) {
	var err error
	differs := false
	err = nil
	if metaOrgn.Name != metaCurr.Name {
		differs = true
		err = errors.New("Diffrent name")
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Name",
			From:  metaOrgn.Name,
			Curr:  metaCurr.Name,
		})
	}

	if metaOrgn.Reference != metaCurr.Reference {
		differs = true
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Reference",
			From:  metaOrgn.Reference,
			Curr:  metaCurr.Reference,
		})
	}

	if metaOrgn.Parent != metaCurr.Parent {
		differs = true
		err = errors.New("Diffrent parent")
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.Parent",
			From:  metaOrgn.Parent,
			Curr:  metaCurr.Parent,
		})
	}

	if metaOrgn.SHA256 != metaCurr.SHA256 {
		differs = true
		change.Diff = append(change.Diff, Difference{
			Field: ctx + "Meta.SHA256",
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
				Field: ctx + "Meta.Annotations[" + k + "]",
				From:  orgVal,
				Curr:  currVal,
			})
		}
	}

	return differs, err
}
