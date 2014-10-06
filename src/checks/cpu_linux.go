package checks

import (
	"io/ioutil"
	"log"
	"strings"
	"strconv"
	"sensu"
	"fmt"
	"time"
)

// PLATFORMS
//   Linux

var cpuStatInterval = 30 * time.Second

func (cpu *CpuStats) Init(q sensu.MessageQueuer, config *sensu.Config) error {
	if err := q.ExchangeDeclare(
		"RESULTS_QUEUE",
		"direct",
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	cpu.q = q
	cpu.config = config
	cpu.close = make(chan bool)

	// get the number of CPUs from: /sys/devices/system/cpu/
	online, err := ioutil.ReadFile("/sys/devices/system/cpu/present")
	cpu.cpu_count = 1
	if nil != err {
		log.Printf("Unable to determine number of CPUs. Intialising only 1 CPU")
	} else {
		online_bits := strings.Split(string(online), "-")
		if len(online_bits) != 2 {
			log.Printf("/sys/devices/system/cpu/present CPU count file malformed. Initialising only 1 CPU")
		} else {
			cpu.cpu_count, err = strconv.Atoi(strings.Trim(online_bits[1], "\n"))
			if nil != err {
				log.Printf("Failed converting CPU count. Initialising on 1 CPU. %s", err)
				cpu.cpu_count = 1
			} else {
				// /sys/devices/system/cpu/present is 0 based
				cpu.cpu_count++
			}
		}
	}
	cpu.frequency = make(map[int]int, cpu.cpu_count)

	return nil
}

func (cpu *CpuStats) Gather() {
	cpu.cpu_freq();
}

func (cpu *CpuStats) cpu_freq() {

}


