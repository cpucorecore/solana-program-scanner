package msg_broker_test

import (
	"fmt"
	"solana-program-scanner/config"
	"solana-program-scanner/log"
	"solana-program-scanner/msg_broker"
	"sync"
	"testing"
	"time"
)

const (
	defaultMqUrl = "amqp://guest:guest@localhost:5672/"
)

func TestMq(t *testing.T) {
	log.InitLoggerForTest()
	conf := &config.MsgBrokerMqConf{
		Mq: config.MqConf{
			MqUrl: defaultMqUrl,
		},
		QueueName:         "slot_queue",
		PublisherCapacity: 1000,
		DeliverCapacity:   100,
	}
	mq := msg_broker.NewMsgBrokerMq[uint64](conf)

	for i := uint64(0); i < 100000; i++ {
		mq.Produce(i)
	}
	mq.Close()

	time.Sleep(time.Second * 110)
}

func TestA(t *testing.T) {
	ch := make(chan uint64, 10)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			v, ok := <-ch
			if !ok {
				break
			}
			fmt.Println(v)
		}
	}(wg)

	go func() {
		for i := 0; i < 10; i++ {
			ch <- uint64(i)
		}
		close(ch)
	}()

	wg.Wait()
}
