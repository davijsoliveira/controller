package main

import (
	"ContainerManager/Commons"
	"ContainerManager/PID"
	"ContainerManager/Sensor"
	"fmt"
	"github.com/docker/docker/client"
	"os"
	"time"
)

func main() {
	// Especifica a versão da API do Docker
	os.Setenv("DOCKER_API_VERSION", "1.42")

	// Conexão com o Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Imagem da aplicação
	imageName := "alpine"

	cpuUsage, err := Sensor.GetCPUUsageByImage(cli, imageName)
	if err != nil {
		fmt.Println("Erro ao obter uso de CPU:", err)
	} else {
		avgCPU := Sensor.CalculateAverageCPU(cpuUsage)
		fmt.Printf("A média de utilização de CPU dos contêineres da imagem %s é: %.2f%%\n", imageName, avgCPU)
	}

	// Instanciar controlador
	controller := PID.NewPIDController(-0.7, 0.005, 0.0)
	measured := 70.0
	//updatedNumberReplicas := 1.0
	stop := false

	for {
		// Total de réplicas atualmente
		//currentNumberReplicas, err := Actuator.GetContainerCount(cli, imageName)

		fmt.Println("                                          ")
		//fmt.Println("Number of Replicas: ", currentNumberReplicas)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Output Measured: %.2f\n", measured)
		inputControl := controller.Update(measured)

		fmt.Printf("Input Control: %.2f\n", inputControl)

		//Commons.Hysteresis(cli, imageName, &updatedNumberReplicas, currentNumberReplicas, controller.LastInputControl, inputControl, controller.Setpoint, measured)
		Commons.Hysteresis(cli, imageName, controller.LastInputControl, inputControl, controller.Setpoint, measured)

		// Implementa uma deadzone ou hysteresis (range 10% superior ou inferior para evitar mudanças frequentes)
		//rangeSetPoint := controller.Setpoint * 0.10
		//upperBound := measured - controller.Setpoint
		//lowerBound := controller.Setpoint - measured
		//if upperBound >= rangeSetPoint {
		//	// Calcula o número de réplicas, acrescentando uma porcentagem baseada no input control, para acelarar a medida que o input control sobe
		//	if inputControl > 5 {
		//		if inputControl > lastInputControl {
		//			//Guarda o valor em float do número de rélicas, e.g., 1.2
		//			if updateNumberReplica < float64(replicas) {
		//				updateNumberReplica = float64(replicas) * 1.1
		//			} else {
		//				updateNumberReplica = updateNumberReplica * 1.1
		//			}
		//
		//			err = Actuator.ScaleOut(cli, imageName, int(updateNumberReplica))
		//			if err != nil {
		//				fmt.Println("Erro ao realizar o scale-out:", err)
		//			}
		//		} else {
		//			updateNumberReplica = float64(replicas)
		//		}
		//	}
		//	fmt.Printf("Estimaded Number of Replicas: %.2f\n", updateNumberReplica)
		//}
		//if lowerBound >= rangeSetPoint {
		//	if inputControl < 5 {
		//		if inputControl <= lastInputControl {
		//			updateNumberReplica -= float64(replicas) * 0.2
		//			if updateNumberReplica < 1 {
		//				updateNumberReplica = 1
		//			}
		//			err := Actuator.ScaleIn(cli, imageName, int(updateNumberReplica))
		//			if err != nil {
		//				fmt.Println("Erro ao realizar o scale-in:", err)
		//			}
		//
		//		} else {
		//			updateNumberReplica = float64(replicas)
		//		}
		//	}
		//	fmt.Printf("Estimaded Number of Replicas: %.2f\n", updateNumberReplica)
		//}

		// Simular mudanças nos valores medidos e de controle
		if measured < 90 && stop == false {
			measured += 1.0
		} else {
			measured -= 1.0
			stop = true
		}

		time.Sleep(time.Second)
		controller.LastInputControl = inputControl
	}
}
