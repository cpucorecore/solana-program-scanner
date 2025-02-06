package config

type PublisherConf struct {
	BufferSize              int
	MqMessagesMin           int
	MqMessagesMax           int
	InspectIntervalBySecond int
	WaitUnitOnErrByMs       int
}

type DeliverConf struct {
	BufferSize int
}

type MsgBrokerMqConf struct {
	MqUrl     string
	QueueName string
	Publisher *PublisherConf
	Deliver   *DeliverConf
}

type MsgBrokerConf struct {
	SlotMqOn            bool
	SlotMq              *MsgBrokerMqConf
	SlotLocalCapacity   int
	BlockCapacity       int
	BlockErrCapacity    int
	BlockParsedCapacity int
	MarketCapacity      int
	PairCapacity        int
	TokenCapacity       int
	TokenErrCapacity    int
}

var defaultMsgBrokerConf = &MsgBrokerConf{
	SlotMqOn: false,
	SlotMq: &MsgBrokerMqConf{
		MqUrl:     "amqp://guest:guest@localhost:5672/",
		QueueName: "slot_queue",
		Publisher: &PublisherConf{
			BufferSize:              1000,
			MqMessagesMin:           100000,
			MqMessagesMax:           1000000,
			InspectIntervalBySecond: 10,
			WaitUnitOnErrByMs:       100,
		},
		Deliver: &DeliverConf{
			BufferSize: 10,
		},
	},
	SlotLocalCapacity:   100,
	BlockCapacity:       10,
	BlockErrCapacity:    100,
	BlockParsedCapacity: 10,
	MarketCapacity:      1000,
	PairCapacity:        100,
	TokenCapacity:       100,
	TokenErrCapacity:    100,
}
