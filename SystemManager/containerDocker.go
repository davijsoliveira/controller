package main

import (
	"ContainerManager/PID"
	"ContainerManager/Sensor"
	"ContainerManager/SharedControl"
	"sync"
)

func main() {
	// Instanciar controlador
	//controller := PID.NewPIDController(-0.5, 0.005, 0.0)
	controller := PID.NewPIDController(-0.5, 0.001, 0.0)

	// Instancia os componentes
	sensor := Sensor.NewSensor()
	movingAvgFilter := SharedControl.NewMovingAverageFilter()
	hysteresisFilter := SharedControl.NewHysteresisFilter()

	// Cria os canais
	sensorToFilter := make(chan int)
	filterToController := make(chan int)
	controllertoHysteris := make(chan []float64)
	//controllerToActuator := make(chan float64)

	// Cria o wait group para controlar as go routines
	var wg sync.WaitGroup
	wg.Add(4)

	go sensor.CountConnections(sensorToFilter)
	go movingAvgFilter.MovingAveragesFilter(sensorToFilter, filterToController)
	go controller.Update(filterToController, controllertoHysteris)
	go hysteresisFilter.Hysteresis(controllertoHysteris)

	wg.Wait()
}
