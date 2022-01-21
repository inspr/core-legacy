package utils

import (
	"fmt"
	"io"
	"strconv"

	"github.com/disiqueira/gotree"
	"inspr.dev/inspr/pkg/meta"
)

// PrintAppTree prints the app tree
func PrintAppTree(app *meta.App, out io.Writer) {
	tree := gotree.New(app.Meta.Name)
	meta := tree.Add("Meta")

	spec := tree.Add("Spec")

	populateMeta(meta, &app.Meta)

	addAppsTree(spec, app)

	addChannelsTree(spec, app)

	addTypesTree(spec, app)

	addAliasesTree(spec, app)

	addRoutesTree(spec, app)

	addNodesTree(spec, app)

	addBoundarysTree(spec, app)

	auth := spec.Add("Auth")
	auth.Add("Scope: " + app.Spec.Auth.Scope)

	addPermissionsTree(spec, app)

	fmt.Fprintln(out, tree.Print())

}

// PrintChannelTree prints the channel structure
func PrintChannelTree(ch *meta.Channel, out io.Writer) {
	channel := gotree.New(ch.Meta.Name)
	meta := channel.Add("Meta")

	populateMeta(meta, &ch.Meta)

	spec := channel.Add("Spec")
	spec.Add("Type: " + ch.Spec.Type)

	if len(ch.Spec.BrokerPriorityList) > 0 {
		brokerList := spec.Add("BrokerPriorityList")
		for _, broker := range ch.Spec.BrokerPriorityList {
			brokerList.Add(broker)
		}
	}

	spec.Add("SelectedBroker: " + ch.Spec.SelectedBroker)

	if len(ch.ConnectedApps) > 0 {
		conApps := channel.Add("ConnectedApps")
		for _, appName := range ch.ConnectedApps {
			conApps.Add(appName)
		}
	}

	if len(ch.ConnectedAliases) > 0 {
		conAliases := channel.Add("ConnectedAliases")
		for _, alias := range ch.ConnectedAliases {
			conAliases.Add(alias)
		}
	}

	fmt.Fprintln(out, channel.Print())
}

// PrintTypeTree prints the type structure
func PrintTypeTree(t *meta.Type, out io.Writer) {
	insprType := gotree.New(t.Meta.Name)
	meta := insprType.Add("Meta")

	populateMeta(meta, &t.Meta)

	spec := insprType.Add("Spec")
	spec.Add("Schema: " + string(t.Schema))

	if len(t.ConnectedChannels) > 0 {
		conChannels := insprType.Add("ConnectedChannels")
		for _, appName := range t.ConnectedChannels {
			conChannels.Add(appName)
		}
	}

	fmt.Fprintln(out, insprType.Print())
}

// PrintAliasTree prints the alias structure
func PrintAliasTree(al *meta.Alias, out io.Writer) {
	alias := gotree.New(al.Meta.Name)
	meta := alias.Add("Meta")

	populateMeta(meta, &al.Meta)

	alias.Add("Resource: " + al.Resource)
	alias.Add("Source: " + al.Source)
	alias.Add("Destination: " + al.Destination)

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

func addAppsTree(spec gotree.Tree, app *meta.App) {
	if len(app.Spec.Apps) > 0 {
		apps := spec.Add("Apps")
		for appName := range app.Spec.Apps {
			apps.Add(appName)
		}
	}
}

func addChannelsTree(spec gotree.Tree, app *meta.App) {
	if len(app.Spec.Channels) > 0 {
		channels := spec.Add("Channels")
		for chName := range app.Spec.Channels {
			channels.Add(chName)
		}
	}
}

func addTypesTree(spec gotree.Tree, app *meta.App) {
	if len(app.Spec.Types) > 0 {
		insprTypes := spec.Add("Types")
		for typeName := range app.Spec.Types {
			insprTypes.Add(typeName)
		}
	}
}

func addAliasesTree(spec gotree.Tree, app *meta.App) {
	if len(app.Spec.Aliases) > 0 {
		aliases := spec.Add("Aliases")
		for aliasKey, alias := range app.Spec.Aliases {
			aliasTree := aliases.Add(aliasKey)
			aliasTree.Add("Resource: " + alias.Resource)
			aliasTree.Add("Source: " + alias.Source)
			aliasTree.Add("Destination: " + alias.Destination)
		}
	}
}

func addRoutesTree(spec gotree.Tree, app *meta.App) {
	if len(app.Spec.Routes) > 0 {
		routes := spec.Add("Routes")
		for routeName, routeConnection := range app.Spec.Routes {
			routeEndpoint := routes.Add(routeName)
			for _, endpoint := range routeConnection.Endpoints {
				routeEndpoint.Add(endpoint)
			}
		}
	}
}

func addNodesTree(spec gotree.Tree, app *meta.App) {
	if app.Spec.Node.Spec.Image != "" {
		node := spec.Add("Node")

		nodeMeta := node.Add("Meta")
		populateMeta(nodeMeta, &app.Spec.Node.Meta)

		nodeSpec := node.Add("Spec")

		nodeSpec.Add("Image: " + app.Spec.Node.Spec.Image)

		if len(app.Spec.Node.Spec.Environment) > 0 {
			env := spec.Add("Environment")
			for name, value := range app.Spec.Types {
				env.Add(fmt.Sprintf("%s: %s", name, value))
			}
		}
		nodeSpec.Add(fmt.Sprintf("Replicas: %d", app.Spec.Node.Spec.Replicas))

		sidecarPort := nodeSpec.Add("SidecarPort")
		sidecarPort.Add(fmt.Sprintf("LBRead: %d", app.Spec.Node.Spec.SidecarPort.LBRead))
		sidecarPort.Add(fmt.Sprintf("LBWrite: %d", app.Spec.Node.Spec.SidecarPort.LBWrite))

		if len(app.Spec.Node.Spec.Ports) > 0 {
			ports := spec.Add("Ports")
			for index, nodePort := range app.Spec.Node.Spec.Ports {
				npIndex := ports.Add("Port " + strconv.Itoa(index+1))
				npIndex.Add(fmt.Sprintf("Port: %d", nodePort.Port))
				npIndex.Add(fmt.Sprintf("TargetPort: %d", nodePort.TargetPort))
			}
		}

	}
}

func addBoundarysTree(spec gotree.Tree, app *meta.App) {
	if len(app.Spec.Boundary.Channels.Input.Union(app.Spec.Boundary.Channels.Output)) > 0 {
		boundary := spec.Add("Boundary")
		if len(app.Spec.Boundary.Channels.Input) > 0 {
			input := boundary.Add("Input")
			for _, ch := range app.Spec.Boundary.Channels.Input {
				input.Add(ch)
			}
		}
		if len(app.Spec.Boundary.Channels.Output) > 0 {
			output := boundary.Add("Output")
			for _, ch := range app.Spec.Boundary.Channels.Output {
				output.Add(ch)
			}
		}
	}
}

func addPermissionsTree(auth gotree.Tree, app *meta.App) {
	if len(app.Spec.Auth.Permissions) > 0 {
		permissions := auth.Add("Permissions")
		for _, permission := range app.Spec.Auth.Permissions {
			permissions.Add(permission)
		}
	}
}
