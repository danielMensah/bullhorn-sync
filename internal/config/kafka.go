package config

// Kafka defines methods to retrieve kafka configuration
type Kafka interface {
	KafkaAddress() string
}

// KafkaAddress returns the kafka address
func (c Config) KafkaAddress() string {
	return c.GetString("kafka.address")
}
