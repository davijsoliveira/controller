package main

import (
	"ContainerManager/Actuator"
	"ContainerManager/PID"
	"ContainerManager/Sensor"
	"ContainerManager/SharedControl"
	"sync"
)

func main() {
	// Instanciar controlador
	//controller := PID.NewPIDController(0.5, 0.005, 0.0) >> Manual
	//controller := PID.NewPIDController(0.5, 0.001, 0.0) >> Manual
	// >> CHR
	//controller := PID.NewPIDController(0.016, 2.22, 0.0)
	// >> Amigo
	//controller := PID.NewPIDController(0.27, 0.06, 0.0)
	// >> Ziegler-Nichols
	//controller := PID.NewPIDController(0.48, 0.29, 0.0)
	// >> Cohen-Coon PID
	//controller := PID.NewPIDController(0.48, 4.13, 0.5)
	// >> Amigo PID
	//controller := PID.NewPIDController(0.27, 0.06, 1.2)
	// >> Ziegler-Nichols PID
	//controller := PID.NewPIDController(0.48, 0.29, 0.06)

	// Instancia os componentes
	sensor := Sensor.NewSensor()
	movingAvgFilter := SharedControl.NewMovingAverageFilter()
	// >> Cohen-Coon
	controller := PID.NewPIDController(0.48, 4.13, 0.0)
	actuator := Actuator.NewActuator()

	// Cria os canais
	sensorToFilter := make(chan int)
	filterToController := make(chan int)
	controllerToActuator := make(chan []float64)

	// Cria o wait group para controlar as go routines
	var wg sync.WaitGroup
	wg.Add(4)

	go sensor.CountConnections(sensorToFilter)
	go movingAvgFilter.MovingAveragesFilter(sensorToFilter, filterToController)
	go controller.Update(filterToController, controllerToActuator)
	//go controller.Update(sensorToFilter, controllerToActuator)
	go actuator.Scale(controllerToActuator)

	wg.Wait()
}
