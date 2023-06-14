package Sensor

import (
	"ContainerManager/ContainersFunc"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type StatsResponse struct {
	TotalRequests     int64 `json:"totalRequests"`
	RequestsPerSecond int64 `json:"requestsPerSecond"`
}

func CountConnections() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastTotalRequests int64
	var lastTime time.Time

	for range ticker.C {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/stats", nil)
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to make request: %v", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Request failed with status: %s", resp.Status)
			continue
		}

		var stats StatsResponse
		err = json.NewDecoder(resp.Body).Decode(&stats)
		if err != nil {
			log.Printf("Failed to decode response: %v", err)
			continue
		}

		currentTime := time.Now()
		elapsedTime := currentTime.Sub(lastTime).Seconds()

		totalRequests := stats.TotalRequests
		requestsPerSecond := int64(float64(totalRequests-lastTotalRequests) / elapsedTime)

		lastTotalRequests = totalRequests
		lastTime = currentTime

		fmt.Printf("Total Requests: %d\n", totalRequests)
		fmt.Printf("Requests per Second: %d\n", requestsPerSecond)
	}
	//client := &http.Client{
	//	Timeout: 10 * time.Second,
	//}
	//
	//ticker := time.NewTicker(1 * time.Second)
	//defer ticker.Stop()
	//
	//for range ticker.C {
	//	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/stats", nil)
	//	if err != nil {
	//		log.Printf("Failed to create request: %v", err)
	//		continue
	//	}
	//
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		log.Printf("Failed to make request: %v", err)
	//		continue
	//	}
	//	defer resp.Body.Close()
	//
	//	if resp.StatusCode != http.StatusOK {
	//		log.Printf("Request failed with status: %s", resp.Status)
	//		continue
	//	}
	//
	//	var stats StatsResponse
	//	err = json.NewDecoder(resp.Body).Decode(&stats)
	//	if err != nil {
	//		log.Printf("Failed to decode response: %v", err)
	//		continue
	//	}
	//
	//	fmt.Printf("Total Requests: %d\n", stats.TotalRequests)
	//	fmt.Printf("Requests per Second: %d\n", stats.RequestsPerSecond)
	//}
}

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
