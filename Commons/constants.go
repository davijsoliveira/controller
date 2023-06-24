package Commons

const ImageName = "alpine"
const App = "processor"
const DeployName = "processor-deploy"

const URLManagedSystem = "http://processor-svc/stats"

// const URLManagedSystem = "http://processor.trt6.jus.br/stats"
//const URLManagedSystem = "http://localhost:8080/stats"

const Kubeconfig = "/app/config"

// const Kubeconfig = "config-TRT"
const Namespace = "default"

// Configure o setpoint para a aplicação em número de requisições por segundo, e.g., 10 requisições por segundo
const SetPoint = 10

// Configure o range desejada para a hysteresis em porcentagem, e.g., 20 é igual a 20%
const HysteresisRange = 20

const LowerBound = float64(SetPoint) - (float64(SetPoint) * (HysteresisRange * 0.01))

const UpperBound = float64(SetPoint) + (float64(SetPoint) * (HysteresisRange * 0.01))
