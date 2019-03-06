package tracer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

// Msg is a struct for roll upload
type Msg struct {
	TableName string                 `json:"table_name"`
	Index     map[string]string      `json:"index"`
	Fields    map[string]interface{} `json:"fields"`
}

type config struct {
	KafkaAddress []string
	KafkaTopic   string
	URL          string
}

var rollConfig config

var (
	producer sarama.AsyncProducer
	enable   bool
)

// Init init config
func Init(address, topic, url string) {
	var err error
	rollConfig.KafkaAddress = strings.Split(address, ",")
	rollConfig.KafkaTopic = topic
	rollConfig.URL = url
	if producer, err = newProducer(rollConfig.KafkaAddress); err != nil {
		panic("roll 初始化异常")
	}
}

func newProducer(brokerList []string) (sarama.AsyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	enable = true
	if err != nil {
		return nil, err
	}
	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				enable = false
			}
		}()
		for err := range producer.Errors() {
			fmt.Println(err)
			producer.Close()
			enable = false
		}
	}()

	return producer, nil
}

func produce(topic string, key string, content string) {
	// if kafka is err, pass
	if enable != true {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}
	// because defer when producer new，if kafka is error can not jump here
	producer.Input() <- msg

}

// TracerWithKafka upload message with kafka and the main thread will be health when upload error
func (m *Msg) TracerWithKafka() {
	if content, err := json.Marshal(m); err != nil {
		fmt.Println(err)
	} else {
		produce(rollConfig.KafkaTopic, "", string(content))
	}
}

// TracerWithURL upload message with http api and the main thread will be health when upload error
func (m *Msg) TracerWithURL() {
	if msg, err := json.Marshal(m); err != nil {
		fmt.Println(err)
	} else {
		body := bytes.NewReader(msg)
		if resp, err := http.Post(rollConfig.URL, "application/json", body); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp.StatusCode)
		}
	}
}
