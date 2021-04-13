package mocks

import "github.com/inspr/inspr/pkg/controller"

//ClientMock test asset for mocking a controller
type ClientMock struct {
	err error
}

//NewClientMock mocks a controller client
func NewClientMock(err error) controller.Interface {
	return &ClientMock{
		err: err,
	}
}

//Apps mocks a app controller
func (cm *ClientMock) Apps() controller.AppInterface {
	return NewAppMock(cm.err)
}

//Channels mocks a channel controller
func (cm *ClientMock) Channels() controller.ChannelInterface {
	return NewChannelMock(cm.err)
}

//ChannelTypes mocks a chanl types controller
func (cm *ClientMock) ChannelTypes() controller.ChannelTypeInterface {
	return NewChannelTypeMock(cm.err)
}

//Authorization mocks a app controller
func (cm *ClientMock) Authorization() controller.AuthorizationInterface {
	return NewAuthMock(cm.err)
}
