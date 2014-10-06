package checks

import (
	"sensu"
	"log"
	"fmt"
	"time"
	"io/ioutil"
	"strconv"
	"strings"
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

func (cpu *CpuStats) Start() {
	clientConfig := cpu.config.Data().Get("client")
	reset := make(chan bool)
	timer := time.AfterFunc(0, func() {
			result := NewResult(clientConfig)
			result.Output = cpu.createCpuFreqPayload(result.Executed)
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

func (cpu *CpuStats) createCpuFreqPayload(timestamp uint) string {



	var payload string
	// now inject our data
	for i := 0; i < cpu.cpu_count; i++ {
		cpu.frequency[i] = 0
		// attempt to load the file
		content, err := ioutil.ReadFile(fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/cpuinfo_cur_freq", i))
		if nil == err {
			// we have content!
			cpu.frequency[i], err = strconv.Atoi(strings.Trim(string(content), "\n"))
			if nil != err {
				log.Printf("Failed to convert '%s' to an int", string(content))
			}
		} else {
			log.Printf("Could not get CPU Freq for CPU %d: %s",i, err)
		}

		payload += fmt.Sprintf("cpu.frequency.current.cpu%d %d %d\n", i, cpu.frequency[i], timestamp)
	}

	return payload
}

func (cpu *CpuStats) publish(result *Result) {
	if err := cpu.q.Publish(
		"cpu",
		"",
		result.GetPayload(),
	); err != nil {
		log.Printf("CpuStats.publish: %v", err)
		return
	}
	log.Print("CPU stats published")
}
