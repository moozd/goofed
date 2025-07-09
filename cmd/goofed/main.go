package main

import (
	"log"

	"github.com/moozd/goofed/internal/graphics"
)

func main() {
	err := graphics.MainLoop()
	if err != nil {
		log.Fatalln(err)
	}
}
