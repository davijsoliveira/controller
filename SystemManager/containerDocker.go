package main

import "ContainerManager/Sensor"

func main() {
	// Calula a média do uso de CPU
	//avgCPU := Sensor.CalculateAverageCPU()
	//fmt.Printf("A média de utilização de CPU dos containers da imagem %s é: %.2f%%\n", Commons.ImageName, avgCPU)
	//
	//// Instanciar controlador
	//controller := PID.NewPIDController(-0.7, 0.005, 0.0)
	//movingAvgFilter := SharedControl.NewMovingAverageFilter()
	//measured := 70.0
	//stop := false

	for {
		Sensor.CountConnections()
		//fmt.Println("                                          ")
		//fmt.Printf("Output Measured: %.2f\n", measured)
		//
		//inputFiltered := movingAvgFilter.MovingAveragesFilter(measured)
		//fmt.Println("Moving Averages Filter Value:", inputFiltered)
		//
		//inputControl := controller.Update(inputFiltered)
		//
		//fmt.Printf("Input Control: %.2f\n", inputControl)
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
	}
}
