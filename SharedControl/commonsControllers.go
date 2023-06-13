package SharedControl

import (
	"ContainerManager/Actuator"
	"ContainerManager/ContainersFunc"
	"fmt"
)

func Hysteresis(lastInputControl float64, inputControl float64, setPoint float64, measured float64) {
	// Implementa uma deadzone ou hysteresis (range 10% superior ou inferior para evitar mudanças frequentes)
	rangeSetPoint := setPoint * 0.10
	upperBound := measured - setPoint
	lowerBound := setPoint - measured

	// Obtem a quantidade atual de réplicas
	ContainersFunc.ContainersStatsRepository.CurrentNumberContainers()

	if upperBound >= rangeSetPoint {
		// Calcula o número de réplicas, acrescentando uma porcentagem baseada no input control, para acelarar a medida que o input control sobe
		if inputControl > 5 {
			if inputControl > lastInputControl {
				//Guarda o valor em float do número de rélicas, e.g., 1.2
				if ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas < float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas) {
					ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas = float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas) * 1.1
				} else {
					ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas = ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas * 1.1
				}

				err := Actuator.ScaleOut(int(ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas))
				if err != nil {
					fmt.Println("Erro ao realizar o scale-out:", err)
				}
			} else {
				ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas = float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
			}
		}
		fmt.Println("Current Number of Replicas: ", ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
		fmt.Printf("Estimaded Number of Replicas: %.2f\n", ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas)
	}
	if lowerBound >= rangeSetPoint {
		if inputControl < 5 {
			if inputControl <= lastInputControl {
				ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas = ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas * 0.9
				if ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas < 1 {
					ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas = 1
				}
				err := Actuator.ScaleIn(int(ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas))
				if err != nil {
					fmt.Println("Erro ao realizar o scale-in:", err)
				}

			} else {
				ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas = float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
			}
		}
		fmt.Println("Current Number of Replicas: : ", ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
		fmt.Printf("Estimaded Number of Replicas: %.2f\n", ContainersFunc.ContainersStatsRepository.UpdatedNumberReplicas)
	}
}
