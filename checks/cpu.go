package checks

import (
	"fmt"
	"log"
	"time"

	"github.com/rmorriso/sensu-client/sensu"
)

// CPU Status for Linux based machines
//
// DESCRIPTION
//  This plugin gets the CPU stats from linux machines and puts them on the wire without prompting for sensu
//
// OUTPUT
//   Graphite plain-text format (name value timestamp\n)
//
// PLATFORMS
//   Linux

type CpuStats struct {
	q      sensu.MessageQueuer
	config *sensu.Config
	close  chan bool

	frequency map[int]int
	cpu_count int
}

var cpuStatInterval = 30 * time.Second

func (cpu *CpuStats) Init(q sensu.MessageQueuer, config *sensu.Config) error {
	err := q.ExchangeDeclare(RESULTS_QUEUE, "direct")
        if err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	cpu.q = q
	cpu.config = config
	cpu.close = make(chan bool)

	return cpu.setup()
}

func (cpu *CpuStats) Start() {
	clientConfig := cpu.config.Data().Get("client")

	reset := make(chan bool)
	timer := time.AfterFunc(0, func() {
		var err error
		result := NewResult(clientConfig)
		result.Output, err = cpu.createCpuFreqPayload(result.Executed)
		if err  != nil {
			result.Status = 1
			result.Output = fmt.Sprintf("Error: %s", err)
			cpu.Stop() // no point in continually reporting the same error.
		}
		cpu.publish(result)
		reset <- true
	})
	defer timer.Stop()

	for {
		select {
		case <-reset:
			timer.Reset(cpuStatInterval)
		case <-cpu.close:
			return
		}
	}
}

func (cpu *CpuStats) Stop() {
	cpu.close <- true
}

func (cpu *CpuStats) publish(result *Result) {
	err := cpu.q.Publish(RESULTS_QUEUE, "", result.GetPayload())
	if err != nil {
		log.Printf("CpuStats.publish: %v", err)
		return
	}
	log.Print("CPU stats published")
}
