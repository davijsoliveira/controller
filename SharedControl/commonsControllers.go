package SharedControl

import (
	"ContainerManager/Actuator"
	"ContainerManager/ContainersFunc"
	"fmt"
)

type MovingAverageFilter struct {
	CpuUtilization []float64
}

func NewMovingAverageFilter() *MovingAverageFilter {
	return &MovingAverageFilter{
		CpuUtilization: []float64{},
	}
}

func (maf *MovingAverageFilter) MovingAveragesFilter(measured float64) float64 {
	maf.CpuUtilization = append(maf.CpuUtilization, measured)
	if len(maf.CpuUtilization) > 3 {
		for i := range maf.CpuUtilization {
			if i == 0 {
				maf.CpuUtilization[i] = measured
			} else {
				maf.CpuUtilization[i-1] = maf.CpuUtilization[i]
			}
		}
		maf.CpuUtilization = maf.CpuUtilization[:3]
	}
	fmt.Println("Valores da janela", maf.CpuUtilization)
	sum := 0.0
	for _, value := range maf.CpuUtilization {
		sum += value
	}
	return sum / float64(len(maf.CpuUtilization))
}

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
				if ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas < float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas) {
					ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas) * 1.1
				} else {
					ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas * 1.1
				}

				err := Actuator.ScaleOut(int(ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas))
				if err != nil {
					fmt.Println("Erro ao realizar o scale-out:", err)
				}
			} else {
				ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = float64(ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
			}
		}
		fmt.Println("Current Number of Replicas: ", ContainersFunc.ContainersStatsRepository.CurrentNumberReplicas)
		fmt.Printf("Scale Number of Replicas: %.2f\n", ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas)
	}
	if lowerBound >= rangeSetPoint {
		if inputControl < 5 {
			if inputControl <= lastInputControl {
				ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas * 0.9
				if ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas < 1 {
					ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas = 1
				}
				err := Actuator.ScaleIn(int(ContainersFunc.ContainersStatsRepository.ScaleNumberReplicas))
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
