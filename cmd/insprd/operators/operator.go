package operators

type OperatorInterface interface {
	Nodes() NodeOperatorInterface
	Channels() ChannelOperatorInterface
}
