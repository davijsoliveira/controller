package Sensor

import (
	"ContainerManager/ContainersFunc"
	"context"
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"time"
)

type Sensor struct {
	Measured int
}

type StatsResponse struct {
	TotalRequests     int `json:"totalRequests"`
	RequestsPerSecond int `json:"requestsPerSecond"`
}

func NewSensor() *Sensor {
	return &Sensor{}
}

func listPods(clientset *kubernetes.Clientset, namespace string) ([]v1.Pod, error) {
	// Criando um contexto
	ctx := context.TODO()

	// Obtendo a lista de Pods no namespace especificado
	podList, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}

func (s *Sensor) CountConnections(toController chan int) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	//ticker := time.NewTicker(1 * time.Second)
	//defer ticker.Stop()

	var lastTotalRequests int
	var lastTime time.Time

	//for range ticker.C {
	for {
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
		s.Measured = int(float64(totalRequests-lastTotalRequests) / elapsedTime)

		lastTotalRequests = totalRequests
		lastTime = currentTime

		//fmt.Printf("Total Requests: %d\n", totalRequests)
		fmt.Printf("Requests per Second: %d\n", s.Measured)

		// Tempo para coleta
		time.Sleep(time.Second * 1)
		toController <- s.Measured
	}
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
