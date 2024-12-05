package stats

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/apache/pulsar-client-go/pulsar"
)

func Monitor_topic(client pulsar.Client, topic string, n_partion uint, subscription string, consumerName string) {
	var wg sync.WaitGroup

	for partition := range n_partion {
		wg.Add(1)

		consume_chan := make(chan pulsar.ConsumerMessage, 2000)

		partitionName := fmt.Sprintf("%s-partition-%d", topic, partition)

		consumer, err := client.Subscribe(pulsar.ConsumerOptions{
			Topic:                       partitionName,
			SubscriptionName:            subscription,
			Name:                        consumerName,
			Type:                        pulsar.Shared,
			SubscriptionInitialPosition: pulsar.SubscriptionPositionLatest,
			MessageChannel:              consume_chan,
			ReceiverQueueSize:           2000,
		})
		if err != nil {
			logrus.Fatalf("Failed create consumer topic: %s partition %d. Reason: %+v", topic, partition, err)
		}

		go monitor_partition(consumer, consume_chan, partitionName)
	}

	wg.Wait()
}

func monitor_partition(consumer pulsar.Consumer, consume_chan <-chan pulsar.ConsumerMessage, partition string) {
	defer consumer.Close()

	var n_read float64 = 0

	last_instant := time.Now()
	last_publish_time := time.Unix(0, 0)
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()

	buckets := make([]*stats, 0)

	buckets = append(buckets, getStats(100))
	buckets = append(buckets, getStats(250))
	buckets = append(buckets, getStats(500))
	buckets = append(buckets, getStats(1000))
	buckets = append(buckets, getStats(2500))
	buckets = append(buckets, getStats(5000))
	buckets = append(buckets, getStats(10000))

	for {
		select {
		case msg := <-consume_chan:
			n_read += 1

			last_publish_time = msg.PublishTime()

			for _, b := range buckets {
				b.addKey(msg.Key())
				min, max, mean, median := b.check()
				if max > 0 {
					fmt.Printf("%s;%d;%d;%d;%.1f;%.1f\n", partition, b.group_size, min, max, mean, median)
				}
			}

			if err := consumer.Ack(msg); err != nil {
				logrus.Warnf("consumer.Acks err: %+v", err)
			}

		case <-tick.C:
			since := time.Since(last_instant)
			last_instant = time.Now()
			logrus.Infof("Topic: %s read rate: %.3f msg/s; (last pulsar time %v)", partition, n_read/float64(since/time.Second), last_publish_time)
			n_read = 0
		}
	}
}
