package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/moozd/goofed/internal/pty"
	"github.com/moozd/goofed/internal/vte"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP)
	defer stop()

	session, err := pty.NewSession(ctx, "zsh")
	defer session.Close()

	if err != nil {
		fmt.Println(err)
	}

	parser := vte.NewParser(ctx, session)
	defer parser.Close()

	session.Write([]byte("echo hi\r\n"))

	for {
		select {
		case token := <-parser.Queue:
			fmt.Println(token.String())
		case <-ctx.Done():
			os.Exit(0)
		}

	}

}
