package tracer

const CONFIG = "KAFKA_ROLL_ADDRESS,KAFKA_ROLL_TOPIC,ROLL_URL"

type Msg struct {
	TableName string                 `json:"table_name"`
	Index     map[string]string      `json:"index"`
	Fields    map[string]interface{} `json:"fields"`
}

type Config struct {
	KafkaRollAddress map[string]string
	KafkaRollTopic   string
	RollUrl          string
}

var config *Config

func init() {

}
