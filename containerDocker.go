package main

import (
	"ContainerManager/PID"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"log"
	"os"
	"time"
)

func runContainer(client *client.Client, imagename string, containername string, port string, inputEnv []string) error {

	// Configured hostConfig:
	// https://godoc.org/github.com/docker/docker/api/types/container#HostConfig
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		LogConfig: container.LogConfig{
			Type:   "json-file",
			Config: map[string]string{},
		},
	}

	// Define Network config (why isn't PORT in here...?:
	// https://godoc.org/github.com/docker/docker/api/types/network#NetworkingConfig
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	gatewayConfig := &network.EndpointSettings{
		Gateway: "gatewayname",
	}
	networkConfig.EndpointsConfig["bridge"] = gatewayConfig

	// Configuration
	// https://godoc.org/github.com/docker/docker/api/types/container#Config
	config := &container.Config{
		Image:    imagename,
		Cmd:      []string{"tail", "-f", "/dev/null"},
		Env:      inputEnv,
		Hostname: fmt.Sprintf("%s-hostnameexample", imagename),
	}

	// Creating the actual container. This is "nil,nil,nil" in every example.
	cont, err := client.ContainerCreate(context.Background(), config, hostConfig, nil, nil, containername)

	if err != nil {
		log.Println(err)
		return err
	}

	// Run the actual container
	client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	log.Printf("Container %s is created", cont.ID)

	return nil
}

func alpine(containerName string) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containerConfig := &container.Config{
		Image: "alpine",
		//Cmd:   []string{"tail", "-f", "/dev/null"},
		Cmd: []string{"sh", "-c", "while true; do yes > /dev/null; done"},
	}
	hostConfig := &container.HostConfig{}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}

func calculateCPUPercentUnix(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	cpuPercent := (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	return cpuPercent
}

func verifyCPU() {
	// Conecte-se ao daemon do Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Especifique o ID ou nome do container
	containerID := "MyAlpine"

	// Crie um contexto para a chamada de API
	ctx := context.Background()
	for {
		time.Sleep(5 * time.Second)
		// Obtenha as estatísticas do container
		stats, err := cli.ContainerStats(ctx, containerID, false)
		if err != nil {
			panic(err)
		}
		defer stats.Body.Close()

		// Decodifique as estatísticas em JSON
		var statJSON types.StatsJSON
		if err := json.NewDecoder(stats.Body).Decode(&statJSON); err != nil {
			panic(err)
		}

		// Obtenha o uso de CPU do container
		cpuPercent := calculateCPUPercentUnix(&statJSON)

		fmt.Printf("Uso de CPU do container %s: %.2f%%\n", containerID, cpuPercent)
	}
}

func GetContainerCount(cli *client.Client, imageName string) (int, error) {
	// Crie um contexto para a chamada de API
	ctx := context.Background()

	// Liste todos os contêineres em execução
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return 0, err
	}

	// Conte os contêineres com base na imagem
	count := 0
	for _, container := range containers {
		if container.Image == imageName {
			count++
		}
	}

	return count, nil
}

// scaleOut verifica e atualiza o número de contêineres com base no valor fornecido
func scaleOut(cli *client.Client, imageName string, newValue int) error {
	// Obtenha o total de contêineres atual
	currentValue, err := GetContainerCount(cli, imageName)
	if err != nil {
		return err
	}

	// Verifique se o novo valor é maior que o valor atual
	if newValue > currentValue {
		// Crie um contexto para a chamada de API
		ctx := context.Background()

		// Calcule a diferença entre o novo valor e o valor atual
		diff := newValue - currentValue

		containerConfig := &container.Config{
			Image: "alpine",
			//Cmd:   []string{"tail", "-f", "/dev/null"},
			Cmd: []string{"sh", "-c", "while true; do yes > /dev/null; done"},
		}

		// Inicie novos contêineres para atingir o novo valor
		for i := 0; i < diff; i++ {
			currentTime := time.Now().Format("20060102150405")
			name := imageName + currentTime
			resp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, name)
			if err != nil {
				return err
			}
			if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				panic(err)
			}
			//time.Sleep(1000 * time.Millisecond)
		}
	}

	return nil
}
func scaleIn(imageName string, numContainersNeeded int) error {
	// Crie um cliente Docker
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// Obtenha a lista de todos os contêineres
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	// Obtenha o número atual de contêineres da imagem
	currentCount := 0
	for _, container := range containers {
		if container.Image == imageName {
			currentCount++
		}
	}

	// Verifique se o número atual de contêineres é maior do que o número necessário
	if currentCount > numContainersNeeded {
		// Calcule a diferença entre o número atual de contêineres e o número necessário
		diff := currentCount - numContainersNeeded

		// Contador para o número de contêineres excluídos
		deletedCount := 0

		// Itere sobre os contêineres e exclua os que correspondem à imagem
		for _, container := range containers {
			if container.Image == imageName {
				err := cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{Force: true})
				if err != nil {
					return err
				}
				deletedCount++

				// Verifique se já excluímos a diferença de contêineres necessária
				if deletedCount == diff {
					break
				}
			}
		}

		fmt.Printf("Foram excluídos %d contêineres em excesso.\n", deletedCount)
	} else {
		fmt.Println("Não há contêineres em excesso.")
	}

	return nil
}

func main() {
	// Especifica a versão da API do Docker
	os.Setenv("DOCKER_API_VERSION", "1.42")

	// Conecte-se ao daemon do Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Especifique a imagem desejada
	imageName := "alpine"

	controller := PID.NewPIDController(-0.7, 0.005, 0.0)
	measured := 70.0 // Exemplo de porcentagem de CPU utilizada medida
	updateReplicas := 1.0
	stop := false
	lastInputControl := 0.0

	for {
		// Total de réplicas atualmente
		replicas, err := GetContainerCount(cli, imageName)
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
		bound := controller.Setpoint * 0.10
		diff := measured - controller.Setpoint
		if diff >= bound {
			// Calcula o número de réplicas, acrescentando uma porcentagem baseada no input control, para acelarar a medida que o input control sobe
			if inputControl > 5 {
				if inputControl > lastInputControl {
					updateReplicas += float64(replicas) * 0.1
					err = scaleOut(cli, imageName, int(updateReplicas))
					if err != nil {
						fmt.Println("Erro ao realizar o scale-out:", err)
					}
				} else {
					updateReplicas = float64(replicas)
				}
			} else {
				if inputControl < lastInputControl {
					updateReplicas -= float64(replicas) * 0.2
					err := scaleIn(imageName, int(updateReplicas))
					if err != nil {
						fmt.Println("Erro ao realizar o scale-in:", err)
					}
				} else {
					updateReplicas = float64(replicas)
				}
			}

			//percentInputControl = inputControl / 100
			//updateReplicas += float64(replicas) * percentInputControl

			fmt.Printf("Estimaded Number of Replicas: %.2f\n", updateReplicas)
			// Atualiza o número de réplicas
			//err = scaleOut(cli, imageName, int(updateReplicas))
			//if err != nil {
			//	panic(err)
			//}
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

	//verifyCPU()

}
