package main

import (
	"ContainerManager/Actuator"
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

	Sensor.VerifyCPU(cli, imageName)

	// Instanciar controlador
	controller := PID.NewPIDController(-0.7, 0.005, 0.0)
	measured := 70.0
	updateReplicas := 1.0
	stop := false
	lastInputControl := 0.0

	for {
		// Total de réplicas atualmente
		replicas, err := Actuator.GetContainerCount(cli, imageName)
		fmt.Println("                                          ")
		fmt.Println("Number of Replicas: ", replicas)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Output Measured: %.2f\n", measured)
		inputControl := controller.Update(measured)

		// Nivela o input control para o mínimo de réplicas
		if inputControl < 1 {
			inputControl = 1
		}

		fmt.Printf("Input Control: %.2f\n", inputControl)

		// Implementa uma deadzone ou hysteresis (range 10% superior ou inferior para evitar mudanças frequentes)
		rangeSetPoint := controller.Setpoint * 0.10
		upperBound := measured - controller.Setpoint
		lowerBound := controller.Setpoint - measured
		if upperBound >= rangeSetPoint {
			// Calcula o número de réplicas, acrescentando uma porcentagem baseada no input control, para acelarar a medida que o input control sobe
			if inputControl > 5 {
				if inputControl > lastInputControl {
					updateReplicas += float64(replicas) * 0.1
					err = Actuator.ScaleOut(cli, imageName, int(updateReplicas))
					if err != nil {
						fmt.Println("Erro ao realizar o scale-out:", err)
					}
				} else {
					updateReplicas = float64(replicas)
				}
			}
			fmt.Printf("Estimaded Number of Replicas: %.2f\n", updateReplicas)
		}
		if lowerBound >= rangeSetPoint {
			if inputControl < 5 {
				if inputControl <= lastInputControl {
					updateReplicas -= float64(replicas) * 0.2
					if updateReplicas < 1 {
						updateReplicas = 1
					}
					err := Actuator.ScaleIn(cli, imageName, int(updateReplicas))
					if err != nil {
						fmt.Println("Erro ao realizar o scale-in:", err)
					}

				} else {
					updateReplicas = float64(replicas)
				}
			}
			fmt.Printf("Estimaded Number of Replicas: %.2f\n", updateReplicas)
		}

		// Simular mudanças nos valores medidos e de controle
		if measured < 90 && stop == false {
			measured += 1.0
		} else {
			measured -= 1.0
			stop = true
		}

		time.Sleep(time.Second)
		lastInputControl = inputControl
	}
}
