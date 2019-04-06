package constants

const (
	ExchangeName      = "exchange1"
	MqUrl             = "amqp://guest:guest@localhost:5672/"
	AdminLoggingQueue = "app.logging"
	AdminRoutingKey   = "app.#"
	ElbowRoutingKey   = "app.technicians.elbow"
	KneeRoutingKey    = "app.technicians.knee"
	HipRoutingKey     = "app.technicians.hip"
)
