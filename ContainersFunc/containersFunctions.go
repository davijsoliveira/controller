package ContainersFunc

import (
	"ContainerManager/Commons"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

type ContainersStats struct {
	ScaleNumberReplicas   float64
	CurrentNumberReplicas int
}

var ContainersStatsRepository = NewContainerStats()

func NewContainerStats() *ContainersStats {
	return &ContainersStats{
		ScaleNumberReplicas:   1.0,
		CurrentNumberReplicas: 1,
	}
}

func GetConnection() *client.Client {
	// Especifica a versão da API do Docker
	os.Setenv("DOCKER_API_VERSION", "1.42")

	// Conexão com o Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return cli
}

func GetKubernetsCLient() (*kubernetes.Clientset, error) {
	// Obtendo o caminho para o arquivo kubeconfig
	//home := homedir.HomeDir()
	//kubeconfig := filepath.Join(home, ".kube", "config")

	// Criando a configuração do cliente do Kubernetes usando o kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", Commons.Kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Criando o cliente do Kubernetes
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset, nil
}

func (stats *ContainersStats) GetReplicaCount() {

	clientset, err := GetKubernetsCLient()

	// Criando um contexto
	ctx := context.TODO()

	// Criando os seletores para buscar os pods com base no nome da aplicação
	labelSelector := fmt.Sprintf("app=%s", Commons.App)
	podList, err := clientset.CoreV1().Pods(Commons.Namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		fmt.Println(err)
	}

	// Verificando se encontrou algum pod
	if len(podList.Items) == 0 {
		fmt.Println("Nenhum pod encontrado para a aplicação %s", Commons.App)
	}

	// Selecionando o primeiro pod da lista
	podApp := podList.Items[0]

	// Obtendo o nome do conjunto de réplicas (ReplicaSet) associado ao pod
	ownerReferences := podApp.ObjectMeta.OwnerReferences
	replicaSetName := ownerReferences[0].Name

	// Obtendo o objeto ReplicaSet usando o nome do namespace e o nome do ReplicaSet
	replicaSet, err := clientset.AppsV1().ReplicaSets(Commons.Namespace).Get(ctx, replicaSetName, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}

	// Obtendo o número de réplicas definido no objeto ReplicaSet
	numReplica := replicaSet.Spec.Replicas
	fmt.Println("Número de réplicas atualmente:", *numReplica)

	stats.CurrentNumberReplicas = int(*numReplica)
}

func (stats *ContainersStats) CurrentNumberContainers() {
	cli := GetConnection()
	currentReplicas, err := GetContainerCount(cli)
	stats.CurrentNumberReplicas = currentReplicas
	if err != nil {
		panic(err)
	}
}

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

// getCPUUsageByImage obtém o uso de CPU de todos os contêineres que usam uma determinada imagem
func GetCPUUsageByImage(cli *client.Client) ([]float64, error) {
	// Obtenha a lista de todos os containeres
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	// Slice para armazenar os valores de uso de CPU
	cpuUsage := make([]float64, 0)

	// Itera sobre os containeres e obtem o uso de CPU dos que correspondem à imagem
	for _, container := range containers {
		if container.Image == Commons.ImageName {
			stats, err := cli.ContainerStats(context.Background(), container.ID, false)
			if err != nil {
				return nil, err
			}

			// Decodifique as estatísticas em JSON
			var statsJSON types.StatsJSON
			if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
				panic(err)
			}

			// Obtem o uso de CPU do container
			cpuPercent := calculateCPUPercentUnix(&statsJSON)
			cpuUsage = append(cpuUsage, cpuPercent)
		}
	}

	return cpuUsage, nil
}

// calculateCPUPercentage calcula a porcentagem de uso de CPU com base nos dados de estatísticas
func calculateCPUPercentUnix(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	cpuPercent := (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	return cpuPercent
}

func GetContainerCount(cli *client.Client) (int, error) {
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
		if container.Image == Commons.ImageName {
			count++
		}
	}

	return count, nil
}
