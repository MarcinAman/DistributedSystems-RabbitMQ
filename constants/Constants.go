package constants

/*
docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
*/
const (
	ExchangeName      = "exchange1"
	MqUrl             = "amqp://guest:guest@localhost:5672/"
	AdminLoggingQueue = "app.logging"
	AdminRoutingKey   = "app.#"
	ElbowRoutingKey   = "app.technicians.elbow"
	KneeRoutingKey    = "app.technicians.knee"
	HipRoutingKey     = "app.technicians.hip"
)
