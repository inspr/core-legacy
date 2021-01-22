package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type ChannelTypeMemoryManager struct {
	root *meta.App
}

func (tmm *TreeMemoryManager) ChannelTypes() memory.ChannelTypeMemory {
	return &ChannelTypeMemoryManager{
		root: tmm.root,
	}
}

func (ctm *ChannelTypeMemoryManager) CreateChannelType(ct *meta.ChannelType, context string) error {

	cur_ct, err := ctm.GetChannelType(context, ct.Meta.Name)
	if cur_ct != nil || err == nil {
		return ierrors.NewError().AlreadyExists().
			Message("Target app already has a '" + context + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return err
	}

	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}

func (ctm *ChannelTypeMemoryManager) GetChannelType(context string, ctName string) (*meta.ChannelType, error) {
	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return nil, err

	}

	err = ierrors.NewError().NotFound().Message("No ChannelType found for query.").Build()

	ct := parentApp.Spec.ChannelTypes[ctName]
	if ct != nil {
		return ct, nil
	}

	return nil, err
}

func (ctm *ChannelTypeMemoryManager) DeleteChannelType(context string, ctName string) error {
	curCt, err := ctm.GetChannelType(context, ctName)
	if curCt == nil || err != nil {
		return ierrors.NewError().BadRequest().
			Message("Target app doesn't contain a '" + context + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return err
	}

	delete(parentApp.Spec.ChannelTypes, ctName)

	curCt, err = ctm.GetChannelType(context, ctName)
	if curCt != nil {
		return ierrors.NewError().InternalServer().
			Message("Couldn't delete '" + context + "' ChannelType from target app").Build()
	}
	if err != nil {
		return err
	}
	return nil
}

func (ctm *ChannelTypeMemoryManager) UpdateChannelType(ct *meta.ChannelType, context string) error {

	cur_ct, err := ctm.GetChannelType(context, ct.Meta.Name)
	if cur_ct == nil {
		return ierrors.NewError().BadRequest().
			Message("Target app doesn't contain a '" + context + "' ChannelType").Build()
	}
	if err != nil {
		return err
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return err
	}

	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}
