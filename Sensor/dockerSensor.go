package Sensor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"
)

func calculateCPUPercentUnix(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	cpuPercent := (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	return cpuPercent
}

func VerifyCPU(cli *client.Client, imageName string) {
	// Especifique o ID ou nome do container
	containerID := "myalpine0"

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
