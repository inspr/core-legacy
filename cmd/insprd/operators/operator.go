package operators

// OperatorInterface is an interface for inspr runtime operators
//
// To implement the interface you need to create two implementations,
// a node implementation, that creates nodes from inspr in the given runtime
// and a channel implementation, that creates channels from inspr in the given
// runtime.
type OperatorInterface interface {
	Nodes() NodeOperatorInterface
	Channels() ChannelOperatorInterface
}
