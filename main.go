package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP)
	defer stop()

	session, err := NewSession(ctx, "zsh")
	defer session.Close()
	if err != nil {
		fmt.Println(err)
	}

	parser := NewParser(ctx, session)
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
