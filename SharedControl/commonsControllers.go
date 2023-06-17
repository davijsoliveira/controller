package SharedControl

import (
	"ContainerManager/Actuator"
	"ContainerManager/ContainersFunc"
	"fmt"
)

type MovingAverageFilter struct {
	NumberRequests []int
}

type HysteresisFilter struct{}

func NewMovingAverageFilter() *MovingAverageFilter {
	return &MovingAverageFilter{
		NumberRequests: []int{},
	}
}

func NewHysteresisFilter() *HysteresisFilter {
	return &HysteresisFilter{}
}

func (maf *MovingAverageFilter) MovingAveragesFilter(fromSensor chan int, toController chan int) {
	for {
		measured := <-fromSensor
		maf.NumberRequests = append(maf.NumberRequests, measured)
		if len(maf.NumberRequests) > 5 {
			for i := range maf.NumberRequests {
				if i == 0 {
					maf.NumberRequests[i] = measured
				} else {
					maf.NumberRequests[i-1] = maf.NumberRequests[i]
				}
			}
			maf.NumberRequests = maf.NumberRequests[:5]
		}
		fmt.Println("Valores da janela", maf.NumberRequests)
		sum := 0
		for _, value := range maf.NumberRequests {
			sum += value
		}
		filtered := sum / len(maf.NumberRequests)
		fmt.Println("Valor da média", filtered)
		toController <- filtered
	}
}

// func Hysteresis(lastInputControl float64, inputControl float64, setPoint float64, measured float64) {
func (histeresisfilter *HysteresisFilter) Hysteresis(fromController chan []float64) {
	for {
		hysteresisFilter := <-fromController
		lastInputControl := hysteresisFilter[0]
		inputControl := hysteresisFilter[1]
		setPoint := hysteresisFilter[2]
		measured := hysteresisFilter[3]

		// Implementa uma deadzone ou hysteresis (range 10% superior ou inferior para evitar mudanças frequentes)
		rangeSetPoint := setPoint * 0.10
		upperBound := measured - setPoint
		lowerBound := setPoint - measured

		// Obtem a quantidade atual de réplicas
		ContainersFunc.ContainersStatsRepository.GetReplicaCount()

		if upperBound >= rangeSetPoint {
			if inputControl > 0.5 {
				if inputControl > lastInputControl {
					//Guarda o valor em float do número de rélicas, e.g., 1.2
					if ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas < float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas) {
						ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas) * 1.1
					} else {
						ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas * 1.1
					}

					err := Actuator.ScaleDeployment(int32(ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas))
					if err != nil {
						fmt.Println("Erro ao realizar o scale-out:", err)
					}
				} else {
					fmt.Println("Scale-out ELSE.........")
					ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
				}
			}
			fmt.Println("Current Number of Replicas: ", ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
			fmt.Printf("Scale Number of Replicas: %.2f\n", ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas)
		}
		if lowerBound >= rangeSetPoint {
			if inputControl < -0.5 {
				if inputControl <= lastInputControl {
					ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas * 0.9
					if ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas < 1 {
						ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = 1
					}
					err := Actuator.ScaleDeployment(int32(ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas))
					if err != nil {
						fmt.Println("Erro ao realizar o scale-in:", err)
					}

				} else {
					ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
				}
			}
			fmt.Println("Current Number of Replicas: : ", ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
			fmt.Printf("Scale Number of Replicas: %.2f\n", ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas)
		}

	}
}
