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

//Types mocks a chanl types controller
func (cm *ClientMock) Types() controller.TypeInterface {
	return NewTypeMock(cm.err)
}

//Authorization mocks a app controller
func (cm *ClientMock) Authorization() controller.AuthorizationInterface {
	return NewAuthMock(cm.err)
}

//Alias mocks a alias controller
func (cm *ClientMock) Alias() controller.AliasInterface {
	return NewAliasMock(cm.err)
}

func (cm *ClientMock) Brokers() controller.BrokersInterface {
	return NewBrokersMock(cm.err)
}
