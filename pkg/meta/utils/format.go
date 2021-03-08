package utils

import (
	"fmt"
	"io"

	"github.com/disiqueira/gotree"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// PrintAppTree prints the app tree
func PrintAppTree(app *meta.App, out io.Writer) {
	tree := gotree.New(app.Meta.Name)
	meta := tree.Add("Meta")

	populateMeta(meta, &app.Meta)

	spec := tree.Add("Spec")
	apps := spec.Add("Apps")
	for appName := range app.Spec.Apps {
		apps.Add(appName)
	}

	channels := spec.Add("Channels")
	for chName := range app.Spec.Channels {
		channels.Add(chName)
	}

	channelTypes := spec.Add("ChannelTypes")
	for ctName := range app.Spec.ChannelTypes {
		channelTypes.Add(ctName)
	}

	spec.Add("Node: " + app.Spec.Node.Meta.Name)

	boundary := spec.Add("Boundary")
	input := boundary.Add("Input")
	for _, ch := range app.Spec.Boundary.Input {
		input.Add(ch)
	}

	output := boundary.Add("Output")
	for _, ch := range app.Spec.Boundary.Output {
		output.Add(ch)
	}

	fmt.Fprintln(out, tree.Print())

}

// PrintChannelTree prints the channel structure
func PrintChannelTree(ch *meta.Channel, out io.Writer) {
	channel := gotree.New(ch.Meta.Name)
	meta := channel.Add("Meta")

	populateMeta(meta, &ch.Meta)

	spec := channel.Add("Spec")
	spec.Add("Type: " + ch.Spec.Type)

	conApps := channel.Add("ConnectedApps")
	for _, appName := range ch.ConnectedApps {
		conApps.Add(appName)
	}

	fmt.Fprintln(out, channel.Print())
}

// PrintChannelTypeTree prints the channel structure
func PrintChannelTypeTree(ct *meta.ChannelType, out io.Writer) {
	channelType := gotree.New(ct.Meta.Name)
	meta := channelType.Add("Meta")

	populateMeta(meta, &ct.Meta)

	spec := channelType.Add("Spec")
	spec.Add("Schema: " + string(ct.Schema))

	conChannels := channelType.Add("ConnectedChannels")
	for _, appName := range ct.ConnectedChannels {
		conChannels.Add(appName)
	}

	fmt.Println(out, channelType.Print())
}

func populateMeta(metaTree gotree.Tree, meta *meta.Metadata) {
	metaTree.Add("Name: " + meta.Name)
	metaTree.Add("Parent: " + meta.Parent)
	metaTree.Add("Reference: " + meta.Reference)
	metaTree.Add("SHA256: " + meta.SHA256)
	annotations := metaTree.Add("Annotations")
	for noteName, note := range meta.Annotations {
		annotations.Add(noteName + ": " + note)
	}
}
