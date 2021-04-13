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
	if len(app.Spec.Apps) > 0 {
		apps := spec.Add("Apps")
		for appName := range app.Spec.Apps {
			apps.Add(appName)
		}
	}
	if len(app.Spec.Channels) > 0 {
		channels := spec.Add("Channels")
		for chName := range app.Spec.Channels {
			channels.Add(chName)
		}
	}
	if len(app.Spec.ChannelTypes) > 0 {
		channelTypes := spec.Add("ChannelTypes")
		for ctName := range app.Spec.ChannelTypes {
			channelTypes.Add(ctName)
		}
	}
	if len(app.Spec.Aliases) > 0 {
		aliases := spec.Add("Aliases")
		for aliasKey := range app.Spec.Aliases {
			aliases.Add(aliasKey)
		}
	}
	if app.Spec.Node.Spec.Image != "" {
		node := spec.Add("Node")
		nodeSpec := node.Add("Spec")
		nodeSpec.Add("Image: " + app.Spec.Node.Spec.Image)
		if len(app.Spec.Node.Spec.Environment) > 0 {
			env := spec.Add("Environment")
			for name, value := range app.Spec.ChannelTypes {
				env.Add(fmt.Sprintf("%s: %s", name, value))
			}
		}
		nodeSpec.Add(fmt.Sprintf("Replicas: %d", app.Spec.Node.Spec.Replicas))

	}
	if len(app.Spec.Boundary.Input.Union(app.Spec.Boundary.Output)) > 0 {
		boundary := spec.Add("Boundary")
		if len(app.Spec.Boundary.Input) > 0 {
			input := boundary.Add("Input")
			for _, ch := range app.Spec.Boundary.Input {
				input.Add(ch)
			}
		}
		if len(app.Spec.Boundary.Output) > 0 {
			output := boundary.Add("Output")
			for _, ch := range app.Spec.Boundary.Output {
				output.Add(ch)
			}
		}
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

	if len(ch.ConnectedApps) > 0 {
		conApps := channel.Add("ConnectedApps")
		for _, appName := range ch.ConnectedApps {
			conApps.Add(appName)
		}
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

	if len(ct.ConnectedChannels) > 0 {
		conChannels := channelType.Add("ConnectedChannels")
		for _, appName := range ct.ConnectedChannels {
			conChannels.Add(appName)
		}
	}

	fmt.Fprintln(out, channelType.Print())
}

// PrintAliasTree prints the alias structure
func PrintAliasTree(al *meta.Alias, out io.Writer) {
	alias := gotree.New(al.Meta.Name)
	meta := alias.Add("Meta")

	populateMeta(meta, &al.Meta)

	alias.Add("Target: " + al.Target)

	fmt.Fprintln(out, alias.Print())
}

func populateMeta(metaTree gotree.Tree, meta *meta.Metadata) {
	metaTree.Add("Name: " + meta.Name)
	if meta.Parent != "" {
		metaTree.Add("Parent: " + meta.Parent)
	}
	if meta.Reference != "" {
		metaTree.Add("Reference: " + meta.Reference)
	}
	if meta.UUID != "" {
		metaTree.Add("UUID: " + meta.UUID)
	}
	var annotations gotree.Tree
	if len(meta.Annotations) > 0 {
		annotations = metaTree.Add("Annotations")
		for noteName, note := range meta.Annotations {
			annotations.Add(noteName + ": " + note)
		}
	}
}
