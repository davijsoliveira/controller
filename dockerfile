# Use a imagem base do Golang
FROM golang:1.16

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie os arquivos go.mod e go.sum para o diretório de trabalho
#COPY go.mod go.sum ./
COPY controller config ./

# Execute o comando go mod download para baixar as dependências
#RUN go mod download

# Copie o código fonte para o diretório de trabalho
#COPY . .

# Compile o código Go
#RUN go build -o controller ContainerManager/SystemManager


# Execute o aplicativo quando o contêiner for iniciado
CMD ["./controller"]