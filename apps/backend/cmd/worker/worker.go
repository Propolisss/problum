package main

import "problum/internal/worker"

func main() {
	w, _ := worker.New()
	w.Run()
}
