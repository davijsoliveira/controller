package PID

import "fmt"

const DeltaTime = 1

type PIDController struct {
	Kp, Ki, Kd                  float64
	Setpoint                    int
	Integral                    float64
	LastError, LastInputControl float64
	SumPrevErrors               float64
	Output                      float64
	Min, Max                    float64
}

func NewPIDController(kp, ki, kd float64) *PIDController {
	return &PIDController{
		Kp:               kp,
		Ki:               ki,
		Kd:               kd,
		Setpoint:         10.0,
		Min:              1.0,
		Max:              5.0,
		Integral:         0.0,
		LastError:        0.0,
		LastInputControl: 0.0,
		SumPrevErrors:    0.0,
		Output:           0.0,
	}
}

func (pid *PIDController) Update(fromSensor chan int, toHysteresis chan []float64) {
	//func (pid *PIDController) Update(fromSensor chan int) {
	for {
		measured := <-fromSensor

		fmt.Println("valor de measured: ", measured)

		// errors
		err := pid.Setpoint - measured

		// Proportional
		proportional := pid.Kp * float64(err)

		// Integrator
		pid.Integral += DeltaTime * float64(err)
		integrator := pid.Integral * pid.Ki

		// Differentiator
		differentiator := pid.Kd * (float64(err) - pid.LastError) / DeltaTime

		// control law
		pid.Output = proportional + integrator + differentiator

		//if pid.Output > pid.Max {
		//	pid.Output = pid.Max
		//} else if pid.Output < pid.Min {
		//	pid.Output = pid.Min
		//}

		pid.LastError = float64(err)
		pid.SumPrevErrors = pid.LastError + float64(err)

		//toActuator <- pid.Output
		fmt.Printf("Input Control: %.2f\n", pid.Output)

		//Preenchendo o slice com os valores para enviar para o filtro de hysteresis
		filterHysteresis := make([]float64, 4)
		filterHysteresis[0] = pid.LastInputControl
		filterHysteresis[1] = pid.Output
		filterHysteresis[2] = float64(pid.Setpoint)
		filterHysteresis[3] = float64(measured)

		// Envia as informações para o filtro de hysteresis
		toHysteresis <- filterHysteresis
	}
}
