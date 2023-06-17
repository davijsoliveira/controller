package main

import (
	"ContainerManager/PID"
	"ContainerManager/Sensor"
	"sync"
)

func main() {
	// Instanciar controlador
	controller := PID.NewPIDController(-0.5, 0.005, 0.0)

	// Instancia o sensor
	sensor := Sensor.NewSensor()
	// Instancia o Filtro
	//movingAvgFilter := SharedControl.NewMovingAverageFilter()
	//measured := 0.0
	//stop := false

	// Cria os canais
	sensorToController := make(chan int)
	//controllerToActuator := make(chan float64)

	// Cria o wait group para controlar as go routines
	var wg sync.WaitGroup
	wg.Add(3)

	//for {
	go sensor.CountConnections(sensorToController)
	//fmt.Println("                                          ")
	//fmt.Printf("Output Measured: %.2f\n", measured)
	//
	//inputFiltered := movingAvgFilter.MovingAveragesFilter(measured)
	//fmt.Println("Moving Averages Filter Value:", inputFiltered)
	//
	//go controller.Update(sensorToController, controllerToActuator)
	go controller.Update(sensorToController)

	//measured := <-controllerToActuator

	//fmt.Printf("Input Control: %.2f\n", measured)
	wg.Wait()
	//
	//// Filtro de hysteresis
	//SharedControl.Hysteresis(controller.LastInputControl, inputControl, controller.Setpoint, measured)
	//
	//// Simular mudanças nos valores medidos e de controle
	//if measured < 90 && stop == false {
	//	measured += 1.0
	//} else {
	//	measured -= 1.0
	//	stop = true
	//}
	//
	//time.Sleep(time.Second)
	//controller.LastInputControl = inputControl
	//}
}
