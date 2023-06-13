package Commons

import (
	"ContainerManager/Actuator"
	"ContainerManager/ContainersFunc"
	"fmt"
	"github.com/docker/docker/client"
)

type ContainersStats struct {
	UpdatedNumberReplicas float64
	CurrentNumberReplicas int
}

func NewContainerStats() *ContainersStats {
	return &ContainersStats{
		UpdatedNumberReplicas: 1.0,
		CurrentNumberReplicas: 1,
	}
}

var ContainersStatsRepository = NewContainerStats()

func (stats *ContainersStats) UpdateCurrentContainers(cli *client.Client, imageName string) {
	currentReplicas, err := ContainersFunc.GetContainerCount(cli, imageName)
	stats.CurrentNumberReplicas = currentReplicas
	if err != nil {
		panic(err)
	}
}

// func Hysteresis(cli *client.Client, imageName string, updatedNumberReplicas *float64, currentNumberReplicas int, lastInputControl float64, inputControl float64, setPoint float64, measured float64) {
func Hysteresis(cli *client.Client, imageName string, lastInputControl float64, inputControl float64, setPoint float64, measured float64) {
	// Implementa uma deadzone ou hysteresis (range 10% superior ou inferior para evitar mudanças frequentes)
	rangeSetPoint := setPoint * 0.10
	upperBound := measured - setPoint
	lowerBound := setPoint - measured

	// Obtem a quantidade atual de réplicas
	ContainersStatsRepository.UpdateCurrentContainers(cli, imageName)

	if upperBound >= rangeSetPoint {
		// Calcula o número de réplicas, acrescentando uma porcentagem baseada no input control, para acelarar a medida que o input control sobe
		if inputControl > 5 {
			if inputControl > lastInputControl {
				//Guarda o valor em float do número de rélicas, e.g., 1.2
				if ContainersStatsRepository.UpdatedNumberReplicas < float64(ContainersStatsRepository.CurrentNumberReplicas) {
					ContainersStatsRepository.UpdatedNumberReplicas = float64(ContainersStatsRepository.CurrentNumberReplicas) * 1.1
				} else {
					ContainersStatsRepository.UpdatedNumberReplicas = ContainersStatsRepository.UpdatedNumberReplicas * 1.1
				}

				err := Actuator.ScaleOut(cli, imageName, int(ContainersStatsRepository.UpdatedNumberReplicas))
				if err != nil {
					fmt.Println("Erro ao realizar o scale-out:", err)
				}
			} else {
				ContainersStatsRepository.UpdatedNumberReplicas = float64(ContainersStatsRepository.CurrentNumberReplicas)
			}
		}
		fmt.Println("Current Number of Replicas: ", ContainersStatsRepository.CurrentNumberReplicas)
		fmt.Printf("Estimaded Number of Replicas: %.2f\n", ContainersStatsRepository.UpdatedNumberReplicas)
	}
	if lowerBound >= rangeSetPoint {
		if inputControl < 5 {
			if inputControl <= lastInputControl {
				ContainersStatsRepository.UpdatedNumberReplicas = ContainersStatsRepository.UpdatedNumberReplicas * 0.9
				if ContainersStatsRepository.UpdatedNumberReplicas < 1 {
					ContainersStatsRepository.UpdatedNumberReplicas = 1
				}
				err := Actuator.ScaleIn(cli, imageName, int(ContainersStatsRepository.UpdatedNumberReplicas))
				if err != nil {
					fmt.Println("Erro ao realizar o scale-in:", err)
				}

			} else {
				ContainersStatsRepository.UpdatedNumberReplicas = float64(ContainersStatsRepository.CurrentNumberReplicas)
			}
		}
		fmt.Println("Current Number of Replicas: : ", ContainersStatsRepository.CurrentNumberReplicas)
		fmt.Printf("Estimaded Number of Replicas: %.2f\n", ContainersStatsRepository.UpdatedNumberReplicas)
	}
}
