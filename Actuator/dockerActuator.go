package Actuator

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"log"
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

func Alpine(containerName string) {
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

// Realiza o scale-out dos containeres
func ScaleOut(cli *client.Client, imageName string, newValue int) error {
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

// Realiza o scale-in dos containeres
func ScaleIn(cli *client.Client, imageName string, numContainersNeeded int) error {

	// Obtenha a lista de todos os contêineres
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	// Obtenha o número atual de contêineres da imagem
	//currentCount := 0
	currentCount, err := GetContainerCount(cli, imageName)
	if err != nil {
		return err
	}
	//for _, container := range containers {
	//	if container.Image == imageName {
	//		currentCount++
	//	}
	//}

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
	}
	return nil
}
