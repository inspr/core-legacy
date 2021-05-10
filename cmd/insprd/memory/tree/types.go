package tree

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils"
	"go.uber.org/zap"
)

// TypeMemoryManager implements the Type interface
// and provides methos for operating on Types
type TypeMemoryManager struct {
	*MemoryManager
}

// Types is a MemoryManager method that provides an access point for Types
func (tmm *MemoryManager) Types() memory.TypeMemory {
	return &TypeMemoryManager{
		MemoryManager: tmm,
	}
}

// Create creates, if it doesn't already exist, a new Type for a given app.
// ct: ChannetType to be created.
// context: Path to reference app (x.y.z...)
func (ctm *TypeMemoryManager) Create(context string, ct *meta.Type) error {
	logger.Info("trying to create a Type",
		zap.String("Type", ct.Meta.Name),
		zap.String("context", context))

	logger.Debug("validating Type structure")
	nameErr := utils.StructureNameIsValid(ct.Meta.Name)
	if nameErr != nil {
		logger.Error("invalid Type name",
			zap.String("ctype", ct.Meta.Name))
		return ierrors.NewError().InnerError(nameErr).Message(nameErr.Error()).Build()
	}

	logger.Debug("checking if Type already exists",
		zap.String("channel", ct.Meta.Name),
		zap.String("context", context))

	_, err := ctm.Get(context, ct.Meta.Name)
	if err == nil {
		logger.Error("Type already exists")
		return ierrors.NewError().AlreadyExists().
			Message("target app already has a '" + ct.Meta.Name + "' Type").Build()
	}

	logger.Debug("getting Type parent dApp")
	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidChannel().
			Message("couldn't create Type " + ct.Meta.Name + "\n" + err.Error()).
			Build()
		return newError
	}

	logger.Debug("adding Type to dApp",
		zap.String("Type", ct.Meta.Name),
		zap.String("context", parentApp.Meta.Name))
	if parentApp.Spec.Types == nil {
		parentApp.Spec.Types = map[string]*meta.Type{}
	}
	ct.Meta = utils.InjectUUID(ct.Meta)
	parentApp.Spec.Types[ct.Meta.Name] = ct
	return nil
}

// Get returns, if it exists, a specific Type from a given app.
// ctName: Name of desired Type.
// context: Path to reference app (x.y.z...)
func (ctm *TypeMemoryManager) Get(context string, ctName string) (*meta.Type, error) {
	logger.Info("trying to get a Type",
		zap.String("Type", ctName),
		zap.String("context", context))

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return nil, ierrors.NewError().BadRequest().InnerError(err).
			Message("target dApp doesn't exist").Build()
	}

	if parentApp.Spec.Types != nil {
		if ct, ok := parentApp.Spec.Types[ctName]; ok {
			return ct, nil
		}
	}

	logger.Debug("unable to get Type in given context",
		zap.String("ctype", ctName),
		zap.String("context", context))

	return nil, ierrors.NewError().NotFound().
		Message("Type not found for given query").
		Build()
}

// Delete deletes, if it exists, a Type from a given app.
// ctName: Name of desired Type.
// context: Path to reference app (x.y.z...)
func (ctm *TypeMemoryManager) Delete(context string, ctName string) error {
	logger.Info("trying to delete a Type",
		zap.String("Type", ctName),
		zap.String("context", context))

	curCt, err := ctm.Get(context, ctName)
	if curCt == nil || err != nil {
		return ierrors.NewError().BadRequest().
			Message("target app doesn't contain a '" + context + "' Type").Build()
	}

	logger.Debug("checking if Type can be deleted")
	if len(curCt.ConnectedChannels) > 0 {
		logger.Error("unable to delete Type for it's being used",
			zap.Any("connected channels", curCt.ConnectedChannels))

		return ierrors.NewError().
			BadRequest().
			Message("Type cannot be deleted as it is being used by other channels").
			Build()
	}

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("target app doesn't exist").Build()
	}

	logger.Debug("removing Type from its parents 'Types' structure",
		zap.String("Type", ctName),
		zap.String("dApp", parentApp.Meta.Name))

	delete(parentApp.Spec.Types, ctName)

	return nil
}

// Update updates, if it exists, a Type of a given app.
// ct: Updated ChannetType to be updated on app
// context: Path to reference app (x.y.z...)
func (ctm *TypeMemoryManager) Update(context string, ct *meta.Type) error {
	logger.Info("trying to update a Type",
		zap.String("Type", ct.Meta.Name),
		zap.String("context", context))

	oldChType, err := ctm.Get(context, ct.Meta.Name)
	if err != nil {
		return ierrors.NewError().BadRequest().
			Message("target app doesn't contain a '" + context + "' Type").Build()
	}

	ct.ConnectedChannels = oldChType.ConnectedChannels
	ct.Meta.UUID = oldChType.Meta.UUID

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("target app doesn't exist").Build()
	}

	logger.Debug("replacing old Type with the new one in dApps 'Types",
		zap.String("channel", ct.Meta.Name),
		zap.String("dApp", parentApp.Meta.Name))

	parentApp.Spec.Types[ct.Meta.Name] = ct
	return nil
}

// TypeRootGetter returns a getter that gets Types from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type TypeRootGetter struct{}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the Type which name is equal to the last query element.
// If the specified Type is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (amm *TypeRootGetter) Get(context string, ctName string) (*meta.Type, error) {
	logger.Info("trying to get a Type (Root Getter)",
		zap.String("Type", ctName),
		zap.String("context", context))

	parentApp, err := GetTreeMemory().Root().Apps().Get(context)
	if err != nil {
		return nil, ierrors.
			NewError().
			BadRequest().
			InnerError(err).
			Message("target dApp does not exist on root").
			Build()
	}

	if parentApp.Spec.Types != nil {
		if ch, ok := parentApp.Spec.Types[ctName]; ok {
			return ch, nil
		}
	}

	logger.Error("unable to get Type in given context (Root Getter)",
		zap.String("ctype", ctName),
		zap.String("context", context))

	return nil, ierrors.
		NewError().
		NotFound().
		Message("Type not found for given query on root").
		Build()
}
