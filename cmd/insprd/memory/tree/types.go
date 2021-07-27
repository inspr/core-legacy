package tree

import (
	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
)

// TypeMemoryManager implements the Type interface
// and provides methos for operating on Types
type TypeMemoryManager struct {
	logger *zap.Logger
	*treeMemoryManager
}

// Types is a MemoryManager method that provides an access point for Types
func (mm *treeMemoryManager) Types() TypeMemory {
	logger.Debug("recovering type manager on the memory tree")
	return &TypeMemoryManager{
		treeMemoryManager: mm,
		logger:            logger.With(zap.String("subSection", "types")),
	}
}

// Create creates, if it doesn't already exist, a new Type for a given app.
// insprType: Type to be created.
// scope: Path to reference app (x.y.z...)
func (tmm *TypeMemoryManager) Create(scope string, insprType *meta.Type) error {
	l := tmm.logger.With(
		zap.String("operation", "create"),
		zap.String("type", insprType.Meta.Name),
		zap.String("scope", scope),
	)
	l.Info("received type creation request")

	l.Debug("validating Type structure")
	nameErr := utils.StructureNameIsValid(insprType.Meta.Name)
	if nameErr != nil {
		l.Error("invalid Type name")
		return ierrors.From(nameErr)
	}

	l.Debug("checking if Type already exists")

	_, err := tmm.Get(scope, insprType.Meta.Name)
	if err == nil {
		l.Info("type already exists")
		return ierrors.New(
			"target app already has a '%v' Type", insprType.Meta.Name,
		).AlreadyExists()
	}

	l.Debug("getting Type parent dApp")
	parentApp, err := tmm.Apps().Get(scope)
	if err != nil {
		newError := ierrors.Wrap(
			ierrors.From(err).InvalidType(),
			"couldn't create Type %v", insprType.Meta.Name,
		)
		return newError
	}

	l.Debug("adding Type to dApp")
	if parentApp.Spec.Types == nil {
		parentApp.Spec.Types = map[string]*meta.Type{}
	}
	insprType.Meta = utils.InjectUUID(insprType.Meta)
	parentApp.Spec.Types[insprType.Meta.Name] = insprType
	l.Debug("type created")
	return nil
}

// Get returns, if it exists, a specific Type from a given app.
// name: Name of desired Type.
// scope: Path to reference app (x.y.z...)
func (tmm *TypeMemoryManager) Get(scope, name string) (*meta.Type, error) {
	l := tmm.logger.With(
		zap.String("operation", "get"),
		zap.String("type", name),
		zap.String("scope", scope),
	)
	l.Debug("received type recovery request")

	parentApp, err := tmm.Apps().Get(scope)
	if err != nil {
		l.Debug("parent app does not exist, returning error")
		return nil, ierrors.Wrap(
			ierrors.From(err).BadRequest(),
			"target dApp doesn't exist",
		)
	}

	if parentApp.Spec.Types != nil {
		if insprType, ok := parentApp.Spec.Types[name]; ok {
			l.Debug("recovered type, returning value")
			return insprType, nil
		}
	}

	l.Debug("unable to get Type in given scope")

	return nil, ierrors.New("Type not found for given query").NotFound()
}

// Delete deletes, if it exists, a Type from a given app.
// name: Name of desired Type.
// scope: Path to reference app (x.y.z...)
func (tmm *TypeMemoryManager) Delete(scope, name string) error {
	l := tmm.logger.With(
		zap.String("operation", "delete"),
		zap.String("type", name),
		zap.String("scope", scope))
	l.Info("received type deletion request")

	currType, err := tmm.Get(scope, name)
	if currType == nil || err != nil {
		l.Debug("unable to find type in tree")
		return ierrors.New(
			"target app doesn't contain a '%v' Type", name,
		).BadRequest()
	}

	l.Debug("checking if Type can be deleted")
	if len(currType.ConnectedChannels) > 0 {
		l.Info("unable to delete Type for it's being used",
			zap.Any("connected channels", currType.ConnectedChannels))

		return ierrors.New(
			"Type cannot be deleted as it is being used by other structures",
		).BadRequest()
	}

	parentApp, err := tmm.Apps().Get(scope)
	if err != nil {
		l.Info("unable to get dApp from memory tree")
		return ierrors.Wrap(
			ierrors.From(err).NotFound(),
			"target app doesn't exist",
		)
	}

	l.Info("removing Type from its parents 'Types' structure")

	delete(parentApp.Spec.Types, name)

	return nil
}

// Update updates, if it exists, a Type of a given app.
// insprType: Updated ChannetType to be updated on app
// scope: Path to reference app (x.y.z...)
func (tmm *TypeMemoryManager) Update(scope string, insprType *meta.Type) error {
	l := logger.With(
		zap.String("operation", "update"),
		zap.String("type", insprType.Meta.Name),
		zap.String("scope", scope),
	)
	l.Info("received request for type update")

	oldChType, err := tmm.Get(scope, insprType.Meta.Name)
	if err != nil {
		l.Debug("unable to find type in the tree")
		return ierrors.New(
			"target app doesn't contain a '%v' Type", insprType.Meta.Name,
		).BadRequest()
	}

	insprType.ConnectedChannels = oldChType.ConnectedChannels
	insprType.Meta.UUID = oldChType.Meta.UUID

	parentApp, err := tmm.Apps().Get(scope)
	if err != nil {
		return ierrors.Wrap(
			ierrors.From(err).InternalServer(),
			"target app doesn't exist",
		)
	}

	l.Info("replacing old Type with the new one")

	parentApp.Spec.Types[insprType.Meta.Name] = insprType
	return nil
}

// TypePermTreeGetter returns a getter that gets Types from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type TypePermTreeGetter struct {
	*PermTreeGetter
	logs *zap.Logger
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the Type which name is equal to the last query element.
// If the specified Type is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (trg *TypePermTreeGetter) Get(scope, name string) (*meta.Type, error) {
	l := trg.logs.With(
		zap.String("operation", "get-root"),
		zap.String("type", name),
		zap.String("scope", scope),
	)
	l.Info("received request for type recovery")

	parentApp, err := trg.Apps().Get(scope)
	if err != nil {
		l.Info("unable to find parent dapp")
		return nil, ierrors.Wrap(
			ierrors.From(err).BadRequest(),
			"target dApp does not exist on root",
		)
	}

	if parentApp.Spec.Types != nil {
		if ch, ok := parentApp.Spec.Types[name]; ok {
			return ch, nil
		}
	}

	l.Info("unable to get Type in given scope (root-tree)")

	return nil, ierrors.New(
		"Type not found for given query on root",
	).NotFound()
}
