package Sensor

import (
	"ContainerManager/ContainersFunc"
	"fmt"
)

// getCPUUsageByImage obtém o uso de CPU de todos os contêineres que usam uma determinada imagem
//func GetCPUUsageByImage(cli *client.Client, imageName string) ([]float64, error) {
//	// Obtenha a lista de todos os containeres
//	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	// Slice para armazenar os valores de uso de CPU
//	cpuUsage := make([]float64, 0)
//
//	// Itera sobre os containeres e obtem o uso de CPU dos que correspondem à imagem
//	for _, container := range containers {
//		if container.Image == imageName {
//			stats, err := cli.ContainerStats(context.Background(), container.ID, false)
//			if err != nil {
//				return nil, err
//			}
//
//			// Decodifique as estatísticas em JSON
//			var statsJSON types.StatsJSON
//			if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
//				panic(err)
//			}
//
//			// Obtem o uso de CPU do container
//			cpuPercent := calculateCPUPercentUnix(&statsJSON)
//			cpuUsage = append(cpuUsage, cpuPercent)
//		}
//	}
//
//	return cpuUsage, nil
//}
//
//// calculateCPUPercentage calcula a porcentagem de uso de CPU com base nos dados de estatísticas
//func calculateCPUPercentUnix(stats *types.StatsJSON) float64 {
//	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
//	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
//	cpuPercent := (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
//	return cpuPercent
//}

// calculateAverageCPU calcula a média de utilização de CPU a partir dos valores de uso de CPU fornecidos
func CalculateAverageCPU() float64 {
	cli := ContainersFunc.GetConnection()
	avgCPU := 0.0
	cpuUsage, err := ContainersFunc.GetCPUUsageByImage(cli)
	if err != nil {
		fmt.Println("Erro ao obter uso de CPU:", err)
	} else {
		if len(cpuUsage) == 0 {
			return 0.0
		}

		totalCPU := 0.0
		for _, cpu := range cpuUsage {
			totalCPU += cpu

		}
		fmt.Println("A quantidade total de containers é:", len(cpuUsage))
		avgCPU = totalCPU / float64(len(cpuUsage))
	}
	return avgCPU
}

//func VerifyCPU(cli *client.Client) {
//	// Especifique o ID ou nome do container
//	containerID := "myalpine0"
//
//	// Crie um contexto para a chamada de API
//	ctx := context.Background()
//	for {
//		time.Sleep(5 * time.Second)
//		// Obtenha as estatísticas do container
//		stats, err := cli.ContainerStats(ctx, containerID, false)
//		if err != nil {
//			panic(err)
//		}
//		defer stats.Body.Close()
//
//		// Decodifique as estatísticas em JSON
//		var statJSON types.StatsJSON
//		if err := json.NewDecoder(stats.Body).Decode(&statJSON); err != nil {
//			panic(err)
//		}
//
//		// Obtenha o uso de CPU do container
//		cpuPercent := calculateCPUPercentUnix(&statJSON)
//
//		fmt.Printf("Uso de CPU do container %s: %.2f%%\n", containerID, cpuPercent)
//	}
//}
