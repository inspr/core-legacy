package tree

import (
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	metautils "gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

// Auxiliar dapp unexported functions
func (amm *AppMemoryManager) recursiveCheckAndRefineApp(app *meta.App, parentApp *meta.App) string {
	structureErrors := validAppStructure(app, parentApp)
	for _, childApp := range app.Spec.Apps {
		return amm.recursiveCheckAndRefineApp(childApp, app)
	}
	return structureErrors
}

func validAppStructure(app, parentApp *meta.App) string {
	errDescription := ""
	var validSubstructure, parentWithoutNode bool

	nameErr := metautils.StructureNameIsValid(app.Meta.Name)
	parentWithoutNode = nodeIsEmpty(parentApp.Spec.Node)
	validSubstructure = nodeIsEmpty(app.Spec.Node) || (len(app.Spec.Apps) == 0)
	validChannels, msg := checkAndUpdateChannels(app)

	boundariesExist := len(app.Spec.Boundary.Input) > 0 || len(app.Spec.Boundary.Output) > 0
	if boundariesExist {
		errDescription = errDescription + validBoundaries(app.Meta.Name, app.Spec.Boundary, parentApp.Spec.Channels)
	}

	if nameErr != nil {
		errDescription = errDescription + "invalid dApp name;"
	}
	if !validSubstructure {
		errDescription = errDescription + "invalid substructure;"
	}
	if !parentWithoutNode {
		errDescription = errDescription + "parent has Node;"
	}
	if !validChannels {
		errDescription = errDescription + msg
	}

	return errDescription
}

func (amm *AppMemoryManager) checkApp(app, parentApp *meta.App) error {
	structureErrors := amm.recursiveCheckAndRefineApp(app, parentApp)
	if structureErrors == "" {
		return nil
	}
	return ierrors.NewError().InvalidApp().Message(structureErrors).Build()
}

func (amm *AppMemoryManager) addAppInTree(app, parentApp *meta.App) {
	updateAppBoundary(app, parentApp)
	if app.Spec.Apps == nil {
		app.Spec.Apps = map[string]*meta.App{}
	}
	app.Meta.Parent = parentApp.Meta.Name
	parentApp.Spec.Apps[app.Meta.Name] = app

	if !nodeIsEmpty(app.Spec.Node) {
		app.Spec.Node.Meta.Parent = app.Meta.Name
		if app.Spec.Node.Meta.Annotations == nil {
			app.Spec.Node.Meta.Annotations = map[string]string{}
		}
	}
}

func checkAndUpdateChannels(app *meta.App) (bool, string) {
	channels := app.Spec.Channels
	chTypes := app.Spec.ChannelTypes
	for ctName := range chTypes {
		nameErr := metautils.StructureNameIsValid(ctName)
		if nameErr != nil {
			return false, "invalid channelType name: " + ctName
		}
	}
	for channelName, channel := range channels {
		nameErr := metautils.StructureNameIsValid(channelName)
		if nameErr != nil {
			return false, "invalid channel name: " + channelName
		}
		if channel.Spec.Type != "" {
			if _, ok := chTypes[channel.Spec.Type]; !ok {
				return false, "invalid channel: using non-existent channel type;"
			}

			for _, appName := range channel.ConnectedApps {
				if _, ok := app.Spec.Apps[appName]; !ok {
					app.Spec.Channels[channelName].ConnectedApps = utils.Remove(channel.ConnectedApps, appName)
				}
				appBoundary := utils.StringSliceUnion(app.Spec.Apps[appName].Spec.Boundary.Input, app.Spec.Apps[appName].Spec.Boundary.Output)
				if !utils.Includes(appBoundary, channelName) {
					app.Spec.Channels[channelName].ConnectedApps = utils.Remove(channel.ConnectedApps, appName)
				}
			}

			connectedChannels := chTypes[channel.Spec.Type].ConnectedChannels
			if !utils.Includes(connectedChannels, channelName) {
				chTypes[channel.Spec.Type].ConnectedChannels = append(connectedChannels, channelName)
			}

		}
	}
	return true, ""
}

func nodeIsEmpty(node meta.Node) bool {
	noAnnotations := node.Meta.Annotations == nil
	noName := node.Meta.Name == ""
	noParent := node.Meta.Parent == ""
	noImage := node.Spec.Image == ""

	return noAnnotations && noName && noParent && noImage
}

func validBoundaries(appName string, bound meta.AppBoundary, parentChannels map[string]*meta.Channel) string {
	appBoundary := utils.StringSliceUnion(bound.Input, bound.Output)

	for _, chName := range appBoundary {
		if parentChannels[chName] == nil {
			return "invalid app boundary - channel '" + chName + "' doesnt exist in parent app;"
		}
	}

	return ""
}

func updateAppBoundary(app *meta.App, parentApp *meta.App) {
	for _, childApp := range app.Spec.Apps {
		updateAppBoundary(childApp, app)
	}
	updateSingleBoundary(app.Meta.Name, app.Spec.Boundary, parentApp.Spec.Channels)
}

func updateSingleBoundary(appName string, bound meta.AppBoundary, parentChannels map[string]*meta.Channel) {
	appBoundary := utils.StringSliceUnion(bound.Input, bound.Output)
	for _, chName := range appBoundary {
		if !utils.Includes(parentChannels[chName].ConnectedApps, appName) {
			parentChannels[chName].ConnectedApps = append(parentChannels[chName].ConnectedApps, appName)
		}
	}
}

func getParentApp(sonQuery string) (*meta.App, error) {
	sonRef := strings.Split(sonQuery, ".")
	parentQuery := strings.Join(sonRef[:len(sonRef)-1], ".")

	parentApp, err := GetTreeMemory().Apps().GetApp(parentQuery)

	return parentApp, err
}
