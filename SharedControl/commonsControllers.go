package SharedControl

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type MovingAverageFilter struct {
	NumberRequests []int
}

func NewMovingAverageFilter() *MovingAverageFilter {
	return &MovingAverageFilter{
		NumberRequests: []int{},
	}
}

func (maf *MovingAverageFilter) MovingAveragesFilter(fromSensor chan int, toController chan int) {
	// Criar um arquivo CSV para escrita
	file, err := os.Create("dados-requisicoes.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Criar um escritor CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Definir o separador como ponto e vírgula
	writer.Comma = ';'

	for {
		measured := <-fromSensor
		maf.NumberRequests = append(maf.NumberRequests, measured)
		if len(maf.NumberRequests) > 5 {
			for i := range maf.NumberRequests {
				if i == 0 {
					maf.NumberRequests[i] = measured
				} else {
					maf.NumberRequests[i-1] = maf.NumberRequests[i]
				}
			}
			maf.NumberRequests = maf.NumberRequests[:5]
		}
		fmt.Println("Valores da janela", maf.NumberRequests)
		sum := 0
		for _, value := range maf.NumberRequests {
			sum += value
		}
		filtered := sum / len(maf.NumberRequests)
		fmt.Println("Valor da média", filtered)

		// GERAÇÃO DE ARQUIVO COM OS VALORES DA MÉDIA
		outputMeasuredStr := strconv.Itoa(filtered)
		err = writer.Write([]string{outputMeasuredStr})
		if err != nil {
			log.Fatal(err)
		}

		// Esvaziar o buffer do escritor CSV
		writer.Flush()

		toController <- filtered
	}
}
