package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	concurrency := 5     // quantos workers eu vou querer rodando
	in := make(chan int) // canal de entrada
	done := make(chan []byte)

	go func() { // incrementar canal infinitamente - simular serviço que vai receber dados
		i := 0
		for {
			in <- i
			i++
		}
	}()

	for x := 0; x < concurrency; x++ { // Cria 5 proccers workers rodando simultaneamente, exemplo 5 videos sendo processados
		go ProcessWorker(in, x) // criar workers
	}
	<-done
}

// Funcao que vai ser gerenciada por um cara maior
func ProcessWorker(in <-chan int, worker int) { // vai ser executada como uma go routine
	for x := range in {
		t := time.Duration(rand.Intn(4) * int(time.Second))
		time.Sleep(t)
		fmt.Printf("Worker %d: %d \n", worker, x)
	}
}
