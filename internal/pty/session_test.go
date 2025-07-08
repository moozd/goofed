package pty

import (
	"bufio"
	"context"
	"strings"
	"testing"
	"time"
)

func TestSession_Echo(t *testing.T) {

	session, err := NewSession(context.Background(), "--norc", "--noprofile")
	if err != nil {
		t.Fatalf("failed to start terminal session: %v", err)
	}
	defer session.Close()

	_, err = session.Write([]byte("echo hello-world\n"))
	if err != nil {
		t.Fatalf("failed to write to session: %v", err)
	}

	output := make(chan string)
	go func() {
		scanner := bufio.NewScanner(session)
		var result strings.Builder
		for scanner.Scan() {
			line := scanner.Text()
			result.WriteString(line + "\n")
			if strings.Contains(line, "hello-world") {
				break
			}
		}
		output <- result.String()
	}()

	select {
	case out := <-output:
		if !strings.Contains(out, "hello-world") {
			t.Errorf("output did not contain expected text: %q", out)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for terminal output")
	}

}
