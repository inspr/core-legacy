package meta

// KafkaProvider is a provider for a kafka service
type KafkaProvider struct {
	BootstrapServer string
	AutoOffsetReset string
}
