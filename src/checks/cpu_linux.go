package checks

import (
	"io/ioutil"
	"log"
	"strings"
	"strconv"
	"fmt"
)

// PLATFORMS
//   Linux


func (cpu *CpuStats) setup() error {

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

func (cpu *CpuStats) createCpuFreqPayload(timestamp uint) (string, error) {
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
			return payload, err
			log.Printf("Could not get CPU Freq for CPU %d: %s",i, err)
		}

		payload += fmt.Sprintf("cpu.frequency.current.cpu%d %d %d\n", i, cpu.frequency[i], timestamp)
	}

	return payload, nil
}


