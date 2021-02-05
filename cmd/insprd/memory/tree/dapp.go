package tree

import (
	"strings"

	"github.com/r3labs/diff/v2"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

// AppMemoryManager implements the App interface
// and provides methos for operating on dApps
type AppMemoryManager struct {
	root *meta.App
}

// Apps is a MemoryManager method that provides an access point for Apps
func (tmm *MemoryManager) Apps() memory.AppMemory {
	return &AppMemoryManager{
		root: tmm.root,
	}
}

// Set defines a set structure
type Set map[string]bool

// GetApp recieves a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the dApp which name is equal to the last query element.
// The root app is returned if the query string is an empty string.
// If the specified dApp is found, it is returned. Otherwise, returns an error.
func (amm *AppMemoryManager) GetApp(query string) (*meta.App, error) {
	if query == "" {
		return amm.root, nil
	}

	reference := strings.Split(query, ".")
	err := ierrors.NewError().NotFound().Message("dApp not found for given query").Build()

	nxtApp := amm.root
	if nxtApp != nil {
		for _, element := range reference {
			nxtApp = nxtApp.Spec.Apps[element]
			if nxtApp == nil {
				return nil, err
			}
		}
		return nxtApp, nil
	}

	return nil, err
}

// CreateApp instantiates a new dApp in the given context.
// If the dApp's information is invalid, returns an error. The same goes for an invalid context.
// In case of context being an empty string, the dApp is created inside the root dApp.
func (amm *AppMemoryManager) CreateApp(app *meta.App, context string) error {
	if strings.Contains(app.Meta.Name, ".") {
		ierrors.NewError().InvalidName().Message("invalid character '.' in dApp's name").Build()
	}

	parentApp, err := amm.GetApp(context)
	if err != nil {
		return err
	}

	structureErrors := validAppStructure(*app, *parentApp)
	if structureErrors == "" {
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

		for _, chName := range app.Spec.Boundary.Input {
			connectedApps := parentApp.Spec.Channels[chName].ConnectedApps
			if !utils.Include(connectedApps, app.Meta.Name) {
				parentApp.Spec.Channels[chName].ConnectedApps = append(connectedApps, app.Meta.Name)
			}
		}

		for _, chName := range app.Spec.Boundary.Output {
			connectedApps := parentApp.Spec.Channels[chName].ConnectedApps
			if !utils.Include(connectedApps, app.Meta.Name) {
				parentApp.Spec.Channels[chName].ConnectedApps = append(connectedApps, app.Meta.Name)
			}
		}

		return nil
	}

	return ierrors.NewError().InvalidApp().Message(structureErrors).Build()
}

// DeleteApp receives a query and searches for the specified dApp through the tree.
// If the dApp is found and it doesn't have any dApps insite of it, it's deleted.
// If it has other dApps inside of itself, those dApps are deleted recursively.
// Channels and Channel Types inside the dApps to be deleted are also deleted
// dApp's reference inside of it's parent is also deleted.
// In case of dApp not found an error is returned.
func (amm *AppMemoryManager) DeleteApp(query string) error {
	if query == "" {
		return ierrors.NewError().BadRequest().Message("can't delete root dApp").Build()
	}

	app, err := amm.GetApp(query)
	if err != nil {
		return err
	}
	parent, errParent := getParentApp(query)
	if errParent != nil {
		return errParent
	}

	for _, ch := range app.Spec.Channels {
		if len(ch.ConnectedApps) > 0 {
			return ierrors.NewError().
				BadRequest().
				Message("cannot delete app: it contain some channel(s) that are been used by other apps.").
				Build()
		}
	}

	for _, chName := range app.Spec.Boundary.Input {
		parent.Spec.Channels[chName].ConnectedApps = utils.
			Removes(parent.Spec.Channels[chName].ConnectedApps, app.Meta.Name)
	}

	for _, chName := range app.Spec.Boundary.Output {
		parent.Spec.Channels[chName].ConnectedApps = utils.
			Removes(parent.Spec.Channels[chName].ConnectedApps, app.Meta.Name)
	}

	delete(parent.Spec.Apps, app.Meta.Name)

	return nil
}

// UpdateApp receives a pointer to a dApp and the path to where this dApp is inside the memory tree.
// If the current dApp is found and the new structure is valid, it's updated.
// Otherwise, returns an error.
func (amm *AppMemoryManager) UpdateApp(app *meta.App, query string) error {
	currentApp, err := amm.GetApp(query)
	if err != nil {
		return err
	}

	if currentApp.Meta.Name != app.Meta.Name {
		return ierrors.NewError().InvalidName().Message("dApp's name mustn't change when updating").Build()
	}
	if !nodeIsEmpty(app.Spec.Node) && !(len(app.Spec.Apps) == 0) {
		return ierrors.NewError().InvalidApp().Message("dApp mustn't have a Node and other dApps at the same time").Build()
	}

	structureError := validUpdateChanges(*currentApp, *app, query)
	if structureError != nil {
		return structureError
	}

	amm.DeleteApp(query)
	sonRef := strings.Split(query, ".")
	parentQuery := strings.Join(sonRef[:len(sonRef)-1], ".")
	amm.CreateApp(app, parentQuery)

	return nil
}

// Auxiliar unexported functions
func validAppStructure(app, parentApp meta.App) string {
	errDescription := ""
	var validName, validSubstructure, parentWithoutNode, validChannels bool
	_, inParentRef := parentApp.Spec.Apps[app.Meta.Name]

	validName = (app.Meta.Name != "") && !inParentRef
	parentWithoutNode = nodeIsEmpty(parentApp.Spec.Node)
	validSubstructure = nodeIsEmpty(app.Spec.Node) || (len(app.Spec.Apps) == 0)
	validChannels = checkChannels(app.Spec.Channels, app.Spec.ChannelTypes)
	boundariesExist := len(app.Spec.Boundary.Input) > 0 || len(app.Spec.Boundary.Output) > 0
	if boundariesExist {
		errDescription = errDescription + validBoundaries(app.Meta.Name, app.Spec.Boundary, parentApp.Spec.Channels)
	}

	if !validName {
		errDescription = errDescription + "invalid dApp name;"
	}
	if !validSubstructure {
		errDescription = errDescription + "invalid substructure;"
	}
	if !parentWithoutNode {
		errDescription = errDescription + "parent has Node;"
	}
	if !validChannels {
		errDescription = errDescription + "invalid channel: using non-existent channel type;"
	}

	return errDescription
}

func checkChannels(channels map[string]*meta.Channel, chTypes map[string]*meta.ChannelType) bool {
	for _, channel := range channels {
		if channel.Spec.Type != "" {
			if _, ok := chTypes[channel.Spec.Type]; !ok {
				return false
			}

			connectedChannels := chTypes[channel.Spec.Type].ConnectedChannels
			if !utils.Include(connectedChannels, channel.Meta.Name) {
				chTypes[channel.Spec.Type].ConnectedChannels = append(connectedChannels, channel.Meta.Name)
			}

		}
	}
	return true
}

func nodeIsEmpty(node meta.Node) bool {
	noAnnotations := node.Meta.Annotations == nil
	noName := node.Meta.Name == ""
	noParent := node.Meta.Parent == ""
	noImage := node.Spec.Image == ""

	return noAnnotations && noName && noParent && noImage
}

func validBoundaries(appName string, bound meta.AppBoundary, parentChannels map[string]*meta.Channel) string {
	boundaryErrors := ""
	if len(parentChannels) == 0 {
		boundaryErrors = boundaryErrors + "parent doesn't have Channels;"
	} else {
		if len(bound.Input) > 0 {
			for _, input := range bound.Input {
				if parentChannels[input] == nil {
					boundaryErrors = boundaryErrors + "invalid input boundary;"
					break
				}

				if !utils.Include(parentChannels[input].ConnectedApps, appName) {
					parentChannels[input].ConnectedApps = append(parentChannels[input].ConnectedApps, appName)
					// boundaryErrors = boundaryErrors + "invalid input boundary - channel doesnt have this app in connectedApps list;"
					// break
				}
			}
		}

		if len(bound.Output) > 0 {
			for _, output := range bound.Output {
				if parentChannels[output] == nil {
					boundaryErrors = boundaryErrors + "invalid output boundary;"
					break
				}

				if !utils.Include(parentChannels[output].ConnectedApps, appName) {
					parentChannels[output].ConnectedApps = append(parentChannels[output].ConnectedApps, appName)
					// boundaryErrors = boundaryErrors + "invalid output boundary - channel doesnt have this app in connectedApps list;"
					// break
				}
			}
		}
	}

	return boundaryErrors
}

func getParentApp(sonQuery string) (*meta.App, error) {
	sonRef := strings.Split(sonQuery, ".")
	parentQuery := strings.Join(sonRef[:len(sonRef)-1], ".")

	parentApp, err := GetTreeMemory().Apps().GetApp(parentQuery)

	return parentApp, err
}

func validUpdateChanges(currentApp, newApp meta.App, query string) error {
	boundChangelog, err := diff.Diff(currentApp.Spec.Boundary, newApp.Spec.Boundary)
	if err != nil {
		return diffError(err)
	}

	if len(boundChangelog) != 0 {
		parent, errParent := getParentApp(query)
		if errParent != nil {
			return errParent
		}

		boundError := validBoundaries(newApp.Meta.Name, newApp.Spec.Boundary, parent.Spec.Channels)
		if boundError != "" {
			return ierrors.NewError().InvalidApp().Message(boundError).Build()
		}
	}

	structuresChangelog, err := checkForChildStructureChanges(currentApp.Spec, newApp.Spec)
	if err != nil {
		return diffError(err)
	}

	if len(structuresChangelog["channel"]) > 0 && invalidChannelChanges(structuresChangelog["channel"], &newApp) {
		return ierrors.NewError().InvalidChannel().Message("channel's parent dApp doesn't contain specified channel type").Build()
	}

	if len(structuresChangelog["app"]) > 0 {
		for changedApp := range structuresChangelog["app"] {
			currApp := currentApp.Spec.Apps[changedApp]
			modifiedApp := newApp.Spec.Apps[changedApp]
			if currApp != nil {
				newQuery := query + "." + changedApp
				structureError := validUpdateChanges(*currApp, *modifiedApp, newQuery)
				if structureError != nil {
					return structureError
				}
			} else {
				delete(newApp.Spec.Apps, modifiedApp.Meta.Name)
				childAppErr := validAppStructure(*modifiedApp, newApp)
				if childAppErr != "" {
					return ierrors.NewError().InvalidApp().Message("invalid child dApp: " + childAppErr).Build()
				}
			}
		}
	}

	return nil
}

func checkForChildStructureChanges(currentStruct, newStruct meta.AppSpec) (map[string]Set, error) {
	changedStructures := map[string]Set{
		"app":     {},
		"channel": {},
		"ctype":   {},
	}

	appChangelog, err := diff.Diff(currentStruct.Apps, newStruct.Apps)
	if err != nil {
		return nil, diffError(err)
	}
	if len(appChangelog) != 0 {
		for _, change := range appChangelog {
			if change.Type != "delete" && !changedStructures["app"][change.Path[0]] {
				changedStructures["app"][change.Path[0]] = true
			}
		}
	}

	channelChangelog, err := diff.Diff(currentStruct.Channels, newStruct.Channels)
	if err != nil {
		return nil, diffError(err)
	}
	if len(channelChangelog) != 0 {
		for _, change := range channelChangelog {
			if change.Type != "delete" && !changedStructures["channel"][change.Path[0]] {
				changedStructures["channel"][change.Path[0]] = true
			}
		}
	}

	ctypeChangelog, err := diff.Diff(currentStruct.ChannelTypes, newStruct.ChannelTypes)
	if err != nil {
		return nil, diffError(err)
	}
	if len(ctypeChangelog) != 0 {
		for _, change := range ctypeChangelog {
			if change.Type != "delete" && !changedStructures["ctype"][change.Path[0]] {
				changedStructures["ctype"][change.Path[0]] = true
			}
		}
	}

	return changedStructures, nil
}

func diffError(err error) error {
	return ierrors.NewError().InnerError(err).InternalServer().Message("couldn't create diff to update dApp").Build()
}

// invalidChannelChanges checks if the channels to be updated have the app as their parent. If so,
// the app must contain the channel types used by the channels, or the channel's Type fiel must be empty.
// Returns true if these conditions are not met. Returns false otherwise
func invalidChannelChanges(changedChannels Set, newApp *meta.App) bool {
	channels := newApp.Spec.Channels
	ctypes := newApp.Spec.ChannelTypes

	if len(ctypes) > 0 {
		for change := range changedChannels {
			ct, ctypeExists := ctypes[channels[change].Spec.Type]
			if channels[change].Meta.Parent == newApp.Meta.Name && channels[change].Spec.Type != "" {
				if !ctypeExists {
					return true
				}
				if !utils.Include(ct.ConnectedChannels, change) {
					ct.ConnectedChannels = append(ct.ConnectedChannels, change)
				}
			}

		}
		return false
	}
	return true
}

/*
ESTRUTURAS PARA CHECAGEM DE DIFFS VÁLIDOS:
dAPP:
	diff:"appmeta":
		METADATA
	diff:"appspec":
		diff:"node":
			diff:"nodemeta":
				METADATA
			diff:"nodespec":
				diff:"image"
		diff:"boundary":
			diff:"input"
			diff:"output"
		diff:"apps"
		diff:"channels"
		diff:"channeltypes"

CHANNEL:
	diff:"channelmeta":
		METADATA
	diff:"channelspec":
		diff:"type"

CHANNELTYPES:
	diff:"ctypemeta":
		METADATA
	diff:"schema"

METADATA:
	diff:"name"
	diff:"reference"
	diff:"annotations"
	diff:"parent"
*/
