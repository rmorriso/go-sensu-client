package checks

import (
	simplejson "github.com/bitly/go-simplejson"
	amqp "github.com/streadway/amqp"
	//	"fmt"
	"encoding/json"
	"log"
	"time"
)

const RESULTS_QUEUE = "results"

type Result struct {
	Name     string        `json:"name"`
	Address  string        `json:"address"`
	Command  string        `json:"command"`
	Executed uint          `json:"executed"`
	Status   int           `json:"status"`
	Output   string        `json:"output"`
	Duration time.Duration `json:"duration"`
	Timeout  int           `json:"timeout"`
	started  time.Time
}

func NewResult(clientConfig *simplejson.Json) *Result {
	result := new(Result)
	result.Name, _ = clientConfig.Get("name").String()
	result.Address, _ = clientConfig.Get("address").String()
	result.Executed = uint(time.Now().Unix())
	result.started = time.Now()

	return result
}

func (result *Result) toJson() []byte {
	json, _ := json.Marshal(result)
	log.Printf(string(json))
	return json
}

func (result *Result) GetPayload() amqp.Publishing {
	return amqp.Publishing{
		ContentType:  "application/octet-stream",
		Body:         result.toJson(),
		DeliveryMode: amqp.Transient,
	}
}
