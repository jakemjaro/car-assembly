package logging

import (
	"fmt"
	"sync"
	"time"
)

type PowertrainType int
type ChassisRobotType int 

const (
	GasEngine PowertrainType = iota
	HybridEngine
	PowertrainTypeN
)

const (
	Titano ChassisRobotType = iota
	MegaForce
	ChassisRobotTypeN
)

type PowertrainAdded struct {
	powertrainType PowertrainType
	produced []uint
	inAssemblyQueue []uint
}

type PowertrainRemoved struct {
	consumer ChassisRobotType
	powertrainType PowertrainType
	consumed []uint
	inAssemblyQueue []uint
}

var mutex sync.Mutex

var powertrainProducers []string = []string{"GasEngine powertrain", "HybridEngine powertrain"}
var powertrainProducerNames []string = []string{"GAS", "HYBRID"}

var powertrainConsumerNames []string = []string{"Titano", "MegaForce"}
var poweredchassisConsumerName string = "RoboMount"

var firsttime = 1
var start time.Time

func elapsedSeconds() float64 {
	t := time.Now()

	if firsttime == 1 {
		firsttime = 0
		start = t
	}
	
	s := t.Sub(start)
	return s.Seconds()
}

// Show that a powertrain has been added to the powertrain queue
// and print the current status
func LogAddedPowertrain(powertrainAdded PowertrainAdded) {
	var idx int
	var total uint

	mutex.Lock()

	// Show what is in the powertrain queue
	fmt.Printf("Powertrain_queue:")
	total = 0
	for idx = 0; idx < int(PowertrainTypeN); idx++ {
		if idx > 0 {
			fmt.Printf(" + ")
		}
		fmt.Printf("%d %s", powertrainAdded.inAssemblyQueue[idx], powertrainProducerNames[idx])
		total += powertrainAdded.inAssemblyQueue[idx]
	}

	fmt.Printf(" = %d. ", total)
	fmt.Printf("Added %s.", powertrainProducers[powertrainAdded.powertrainType])

	total = 0
	fmt.Printf(" Produced: ")
	for idx = 0; idx < int(PowertrainTypeN); idx++ {
		total += powertrainAdded.produced[idx]
		if idx > 0 {
			fmt.Printf(" + ")
		}
		fmt.Printf("%d %s", powertrainAdded.produced[idx], powertrainProducerNames[idx])
	}
	fmt.Printf(" = %d in %.3f s.\n", total, elapsedSeconds());

	mutex.Unlock()
}

func LogRemovedPowertrain(powertrainRemoved PowertrainRemoved) {
	var idx int
	var total uint

	mutex.Lock()

	total = 0
	fmt.Printf("Powertrain_queue: ")
	for idx = 0; idx < int(PowertrainTypeN); idx++ {
		if idx > 0 {
			fmt.Printf(" + ")
		}
		fmt.Printf("%d %s", powertrainRemoved.inAssemblyQueue[idx], powertrainProducerNames[idx])
		total += powertrainRemoved.inAssemblyQueue[idx]
	}
	fmt.Printf(" = %d. ", total)

	fmt.Printf("%s consumed %s. %s totals: ", 
		powertrainProducerNames[powertrainRemoved.consumer],
		powertrainProducers[powertrainRemoved.powertrainType],
		powertrainConsumerNames[powertrainRemoved.consumer])
	
	total = 0
	for idx = 0; idx < int(PowertrainTypeN); idx++ {
		if idx > 0 {
			fmt.Printf(" + ")
		}
		total += powertrainRemoved.consumed[idx]
		fmt.Printf("%d %s", powertrainRemoved.consumed[idx], powertrainProducerNames[idx])
	}

	fmt.Printf(" = %d consumed in %.3f s.\n", total, elapsedSeconds())

	mutex.Unlock()
}

func LogAddedPoweredChassis(poweredChassisAdded string, poweredChassisQueueSize uint) {
	mutex.Lock()

	fmt.Printf("Poweredchassis_queue: produced and added %s in %.3f s, queue size: %d\n",
		poweredChassisAdded, elapsedSeconds(), poweredChassisQueueSize)

	mutex.Unlock()
}

func LogRemovedPoweredChassis(poweredChassisRemoved string, poweredChassisQueueSize uint, totalConsumed uint) {
	mutex.Lock()

	fmt.Printf("Poweredchassis_queue: removed and consumed %s in %.3f s, queue size: %d, total consumed: %d\n", 
		poweredChassisRemoved, elapsedSeconds(), poweredChassisQueueSize, totalConsumed)

	mutex.Unlock()
}

func LogPowertrainHistory(produced []uint, consumed [][]int) {
	var p, c int
	var total int
  	fmt.Printf("\nREQUEST REPORT\n----------------------------------------\n");

	for p = 0; p < int(PowertrainTypeN); p++ {
		fmt.Printf("%s producer generated %d requests\n", powertrainProducers[p], produced[p])
	}

	for c = 0; c < int(ChassisRobotTypeN); c++ {
		fmt.Printf("%s consumed ", powertrainConsumerNames[c])
		total = 0
		for p = 0; p < int(PowertrainTypeN); p++ {
			if p > 0 {
				fmt.Printf(" + ")
			}

			total += consumed[c][p]
			fmt.Printf("%d %s", consumed[c][p], powertrainProducerNames[p])
		}
		fmt.Printf(" = %d total\n", total)
	}

	fmt.Printf("Elapsed time %.3f s\n", elapsedSeconds())
}
