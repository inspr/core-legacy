package tree

import (
	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
)

// TypeMemoryManager implements the Type interface
// and provides methos for operating on Types
type TypeMemoryManager struct {
	*MemoryManager
}

// Types is a MemoryManager method that provides an access point for Types
func (mm *MemoryManager) Types() memory.TypeMemory {
	return &TypeMemoryManager{
		MemoryManager: mm,
	}
}

// Create creates, if it doesn't already exist, a new Type for a given app.
// insprType: Type to be created.
// scope: Path to reference app (x.y.z...)
func (tmm *TypeMemoryManager) Create(scope string, insprType *meta.Type) error {
	logger.Info("trying to create a Type",
		zap.String("type", insprType.Meta.Name),
		zap.String("scope", scope))

	logger.Debug("validating Type structure")
	nameErr := utils.StructureNameIsValid(insprType.Meta.Name)
	if nameErr != nil {
		logger.Error("invalid Type name",
			zap.String("type", insprType.Meta.Name))
		return ierrors.NewError().InnerError(nameErr).Message(nameErr.Error()).Build()
	}

	logger.Debug("checking if Type already exists",
		zap.String("type", insprType.Meta.Name),
		zap.String("scope", scope))

	_, err := tmm.Get(scope, insprType.Meta.Name)
	if err == nil {
		logger.Error("Type already exists")
		return ierrors.NewError().AlreadyExists().
			Message("target app already has a '%v' Type", insprType.Meta.Name).Build()
	}

	logger.Debug("getting Type parent dApp")
	parentApp, err := GetTreeMemory().Apps().Get(scope)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidType().
			Message("couldn't create Type %v : %v", insprType.Meta.Name, err.Error()).
			Build()
		return newError
	}

	logger.Debug("adding Type to dApp",
		zap.String("type", insprType.Meta.Name),
		zap.String("parent dApp", parentApp.Meta.Name))
	if parentApp.Spec.Types == nil {
		parentApp.Spec.Types = map[string]*meta.Type{}
	}
	insprType.Meta = utils.InjectUUID(insprType.Meta)
	parentApp.Spec.Types[insprType.Meta.Name] = insprType
	return nil
}

// Get returns, if it exists, a specific Type from a given app.
// name: Name of desired Type.
// scope: Path to reference app (x.y.z...)
func (tmm *TypeMemoryManager) Get(scope, name string) (*meta.Type, error) {
	logger.Info("trying to get a Type",
		zap.String("type", name),
		zap.String("scope", scope))

	parentApp, err := GetTreeMemory().Apps().Get(scope)
	if err != nil {
		return nil, ierrors.NewError().BadRequest().InnerError(err).
			Message("target dApp doesn't exist").Build()
	}

	if parentApp.Spec.Types != nil {
		if insprType, ok := parentApp.Spec.Types[name]; ok {
			return insprType, nil
		}
	}

	logger.Debug("unable to get Type in given scope",
		zap.String("type", name),
		zap.String("scope", scope))

	return nil, ierrors.NewError().NotFound().
		Message("Type not found for given query").
		Build()
}

// Delete deletes, if it exists, a Type from a given app.
// name: Name of desired Type.
// scope: Path to reference app (x.y.z...)
func (tmm *TypeMemoryManager) Delete(scope, name string) error {
	logger.Info("trying to delete a Type",
		zap.String("type", name),
		zap.String("scope", scope))

	currType, err := tmm.Get(scope, name)
	if currType == nil || err != nil {
		return ierrors.NewError().BadRequest().
			Message("target app doesn't contain a '%v' Type", name).Build()
	}

	logger.Debug("checking if Type can be deleted")
	if len(currType.ConnectedChannels) > 0 {
		logger.Error("unable to delete Type for it's being used",
			zap.Any("connected channels", currType.ConnectedChannels))

		return ierrors.NewError().
			BadRequest().
			Message("Type cannot be deleted as it is being used by other structures").
			Build()
	}

	parentApp, err := GetTreeMemory().Apps().Get(scope)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("target app doesn't exist").Build()
	}

	logger.Debug("removing Type from its parents 'Types' structure",
		zap.String("Type", name),
		zap.String("dApp", parentApp.Meta.Name))

	delete(parentApp.Spec.Types, name)

	return nil
}

// Update updates, if it exists, a Type of a given app.
// insprType: Updated ChannetType to be updated on app
// scope: Path to reference app (x.y.z...)
func (tmm *TypeMemoryManager) Update(scope string, insprType *meta.Type) error {
	logger.Info("trying to update a Type",
		zap.String("type", insprType.Meta.Name),
		zap.String("scope", scope))

	oldChType, err := tmm.Get(scope, insprType.Meta.Name)
	if err != nil {
		return ierrors.NewError().BadRequest().
			Message("target app doesn't contain a '%v' Type", insprType.Meta.Name).Build()
	}

	insprType.ConnectedChannels = oldChType.ConnectedChannels
	insprType.Meta.UUID = oldChType.Meta.UUID

	parentApp, err := GetTreeMemory().Apps().Get(scope)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("target app doesn't exist").Build()
	}

	logger.Debug("replacing old Type with the new one in dApps 'Types'",
		zap.String("inspr-type", insprType.Meta.Name),
		zap.String("dApp", parentApp.Meta.Name))

	parentApp.Spec.Types[insprType.Meta.Name] = insprType
	return nil
}

// TypeRootGetter returns a getter that gets Types from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type TypeRootGetter struct{}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the Type which name is equal to the last query element.
// If the specified Type is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (trg *TypeRootGetter) Get(scope, name string) (*meta.Type, error) {
	logger.Info("trying to get a Type (Root Getter)",
		zap.String("type", name),
		zap.String("scope", scope))

	parentApp, err := GetTreeMemory().Root().Apps().Get(scope)
	if err != nil {
		return nil, ierrors.
			NewError().
			BadRequest().
			InnerError(err).
			Message("target dApp does not exist on root").
			Build()
	}

	if parentApp.Spec.Types != nil {
		if ch, ok := parentApp.Spec.Types[name]; ok {
			return ch, nil
		}
	}

	logger.Error("unable to get Type in given scope (Root Getter)",
		zap.String("type", name),
		zap.String("scope", scope))

	return nil, ierrors.
		NewError().
		NotFound().
		Message("Type not found for given query on root").
		Build()
}
