package PID

const DeltaTime = 1

type PIDController struct {
	Kp, Ki, Kd         float64
	Setpoint, Integral float64
	LastError          float64
	SumPrevErrors      float64
	Output             float64
	Min, Max           float64
}

func NewPIDController(kp, ki, kd float64) *PIDController {
	return &PIDController{
		Kp:            kp,
		Ki:            ki,
		Kd:            kd,
		Setpoint:      70.0,
		Min:           1.0,
		Max:           5.0,
		Integral:      0.0,
		LastError:     0.0,
		SumPrevErrors: 0.0,
		Output:        0.0,
	}
}

func (pid *PIDController) Update(measured float64) float64 {
	// errors
	err := pid.Setpoint - measured

	// Proportional
	proportional := pid.Kp * err

	// Integrator
	pid.Integral += DeltaTime * err
	integrator := pid.Integral * pid.Ki

	// Differentiator
	differentiator := pid.Kd * (err - pid.LastError) / DeltaTime

	// control law
	pid.Output = proportional + integrator + differentiator

	//if pid.Output > pid.Max {
	//	pid.Output = pid.Max
	//} else if pid.Output < pid.Min {
	//	pid.Output = pid.Min
	//}

	pid.LastError = err
	pid.SumPrevErrors = pid.LastError + err

	return pid.Output
}
