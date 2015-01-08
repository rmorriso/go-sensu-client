package checks

import (
	"fmt"
	"github.com/rmorriso/sensu-client/sensu"
	"log"
	"time"
)

// CPU Status for Linux based machines
//
// DESCRIPTION
//  This plugin gets the load average and reports it in graphite line format
//
// OUTPUT
//   Graphite plain-text format (name value timestamp\n)
//
// PLATFORMS
//   Linux

type LoadStats struct {
	q      sensu.MessageQueuer
	config *sensu.Config
	close  chan bool

	frequency map[int]int
	cpu_count int
}

var loadAvgInterval = 30 * time.Second

func (load *LoadStats) Init(q sensu.MessageQueuer, config *sensu.Config) error {
	if err := q.ExchangeDeclare(
		RESULTS_QUEUE,
		"direct",
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	load.q = q
	load.config = config
	load.close = make(chan bool)

	return nil
}

func (load *LoadStats) Start() {
	clientConfig := load.config.Data().Get("client")

	reset := make(chan bool)
	timer := time.AfterFunc(0, func() {
		var err error
		result := NewResult(clientConfig)
		result.Output, err = load.createLoadAveragePayload(result.Executed)
		if nil != err {
			result.Status = 1
			result.Output = fmt.Sprintf("Error: %s", err)
			load.Stop() // no point in continually reporting the same error.
		}
		load.publish(result)
		reset <- true
	})
	defer timer.Stop()

	for {
		select {
		case <-reset:
			timer.Reset(loadAvgInterval)
		case <-load.close:
			return
		}
	}
}

func (load *LoadStats) Stop() {
	load.close <- true
}

func (load *LoadStats) publish(result *Result) {
	if err := load.q.Publish(
		RESULTS_QUEUE,
		"",
		result.GetPayload(),
	); err != nil {
		log.Printf("LoadAvg.publish: %v", err)
		return
	}
	log.Print("Load Avg stats published")
}
