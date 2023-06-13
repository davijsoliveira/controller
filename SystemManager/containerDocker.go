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

	avgCPU := Sensor.CalculateAverageCPU(cli, imageName)
	fmt.Printf("A média de utilização de CPU dos contêineres da imagem %s é: %.2f%%\n", imageName, avgCPU)

	// Instanciar controlador
	controller := PID.NewPIDController(-0.7, 0.005, 0.0)
	measured := 70.0

	//updatedNumberReplicas := 1.0
	stop := false

	for {
		fmt.Println("                                          ")
		fmt.Printf("Output Measured: %.2f\n", measured)
		inputControl := controller.Update(measured)

		fmt.Printf("Input Control: %.2f\n", inputControl)

		// Filtro de hysteresis
		Commons.Hysteresis(cli, imageName, controller.LastInputControl, inputControl, controller.Setpoint, measured)

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
