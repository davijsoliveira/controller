package Commons

import (
	"ContainerManager/Actuator"
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
		CurrentNumberReplicas: 0,
	}
}

func (stats *ContainersStats) updateNumberContainers(cli *client.Client, imageName string) {
	currentReplicas, err := Actuator.GetContainerCount(cli, imageName)
	stats.CurrentNumberReplicas = currentReplicas
	if err != nil {
		panic(err)
	}
}

// func Hysteresis(cli *client.Client, imageName string, updatedNumberReplicas *float64, currentNumberReplicas int, lastInputControl float64, inputControl float64, setPoint float64, measured float64) {
func Hysteresis(cli *client.Client, imageName string, lastInputControl float64, inputControl float64, setPoint float64, measured float64) {
	// Cria um struct para gerenciar os numeros de containers
	numReplicas := NewContainerStats()
	numReplicas.updateNumberContainers(cli, imageName)

	// Implementa uma deadzone ou hysteresis (range 10% superior ou inferior para evitar mudanças frequentes)
	rangeSetPoint := setPoint * 0.10
	upperBound := measured - setPoint
	lowerBound := setPoint - measured
	if upperBound >= rangeSetPoint {
		// Calcula o número de réplicas, acrescentando uma porcentagem baseada no input control, para acelarar a medida que o input control sobe
		if inputControl > 5 {
			if inputControl > lastInputControl {
				//Guarda o valor em float do número de rélicas, e.g., 1.2
				if numReplicas.UpdatedNumberReplicas < float64(numReplicas.CurrentNumberReplicas) {
					numReplicas.UpdatedNumberReplicas = float64(numReplicas.CurrentNumberReplicas) * 1.1
				} else {
					numReplicas.UpdatedNumberReplicas = numReplicas.UpdatedNumberReplicas * 1.1
				}

				err := Actuator.ScaleOut(cli, imageName, int(numReplicas.UpdatedNumberReplicas))
				if err != nil {
					fmt.Println("Erro ao realizar o scale-out:", err)
				}
			} else {
				numReplicas.UpdatedNumberReplicas = float64(numReplicas.CurrentNumberReplicas)
			}
		}
		fmt.Printf("Estimaded Number of Replicas: %.2f\n", numReplicas.UpdatedNumberReplicas)
	}
	if lowerBound >= rangeSetPoint {
		if inputControl < 5 {

			if inputControl <= lastInputControl {
				numReplicas.UpdatedNumberReplicas = numReplicas.UpdatedNumberReplicas * 0.9
				if numReplicas.UpdatedNumberReplicas < 1 {
					numReplicas.UpdatedNumberReplicas = 1
				}
				err := Actuator.ScaleIn(cli, imageName, int(numReplicas.UpdatedNumberReplicas))
				if err != nil {
					fmt.Println("Erro ao realizar o scale-in:", err)
				}

			} else {
				numReplicas.UpdatedNumberReplicas = float64(numReplicas.CurrentNumberReplicas)
			}
		}
		fmt.Printf("Estimaded Number of Replicas: %.2f\n", numReplicas.UpdatedNumberReplicas)
	}
}
