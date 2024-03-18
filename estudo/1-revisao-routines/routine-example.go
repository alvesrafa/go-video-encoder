package main

import (
	"fmt"
	"math/rand"
	"time"
)

func hello(msg string) {
	fmt.Println(msg + " - GoRoutine")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Second) // Tempo vai variar
}

func main() {
	go hello("Hello 1")
	go hello("Hello 2")
	time.Sleep(1 * time.Second)
	fmt.Println("Chamada normal")
}

// As chamadas do hello são executadas simultaneamente, mas a chamada normal é executada primeiro.
// Ou seja, a chamada normal é executada antes de as chamadas do hello.

// Todas as vezes que vc chama uma goRoutine, as chamadas sao executadas em "pararelo" (quase) basicamente iniciar simultaneamente
