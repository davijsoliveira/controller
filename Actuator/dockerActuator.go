package Actuator

import (
	"ContainerManager/Commons"
	"ContainerManager/ContainersFunc"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func ScaleDeployment(replicas int32) error {
	clientset, err := ContainersFunc.GetKubernetsCLient()

	ctx := context.TODO()

	// Obtendo as informações do deployment atual
	deployment, err := clientset.AppsV1().Deployments(Commons.Namespace).Get(ctx, Commons.DeployName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Atualizando o número de réplicas
	deployment.Spec.Replicas = &replicas

	// Atualizando o deployment
	_, err = clientset.AppsV1().Deployments(Commons.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Realiza o scale-out dos containeres
func ScaleOut(newValue int) error {
	cli := ContainersFunc.GetConnection()
	// Obtem o total de containers atual
	currentValue, err := ContainersFunc.GetContainerCount(cli)
	if err != nil {
		return err
	}

	// Verifica se o novo valor é maior que o valor atual
	if newValue > currentValue {
		// Crie um contexto para a chamada de API
		ctx := context.Background()

		// Calcula a diferença entre o novo valor e o valor atual
		diff := newValue - currentValue

		containerConfig := &container.Config{
			Image: "alpine",
			Cmd:   []string{"sh", "-c", "while true; do yes > /dev/null; done"},
		}

		// Inicia novos containers para atingir o valor estimado
		for i := 0; i < diff; i++ {
			currentTime := time.Now().Format("20060102150405")
			name := Commons.ImageName + currentTime
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
func ScaleIn(numContainersNeeded int) error {
	cli := ContainersFunc.GetConnection()

	// Obtenha a lista de todos os containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	// Obtenha o número atual de containers da imagem
	currentCount, err := ContainersFunc.GetContainerCount(cli)
	if err != nil {
		return err
	}

	// Verifique se o número atual de containers é maior do que o número necessário
	if currentCount > numContainersNeeded {
		// Calcule a diferença entre o número atual de containers e o número necessário
		diff := currentCount - numContainersNeeded

		// Contador para o número de containers excluídos
		deletedCount := 0

		// Itera no slice de containers e exclue os que correspondem à imagem
		for _, container := range containers {
			if container.Image == Commons.ImageName {
				err := cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{Force: true})
				if err != nil {
					return err
				}
				deletedCount++

				// Verifica se todos os containers foram excluídos
				if deletedCount == diff {
					break
				}
			}
		}
	}
	return nil
}
