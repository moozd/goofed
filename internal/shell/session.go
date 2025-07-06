package shell

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
)

type Session struct {
	cmd    *exec.Cmd
	pty    *os.File
	stdin  io.Writer
	stdout io.Reader
	done   chan struct{}
}

func NewSession(shell string, args ...string) (*Session, error) {
	cmd := exec.Command(shell, args...)
	fd, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	session := &Session{
		cmd:    cmd,
		pty:    fd,
		stdin:  fd,
		stdout: fd,
		done:   make(chan struct{}),
	}

	go session.handleWindowResize()
	go session.forwardSignals()

	return session, nil
}

func (s *Session) Write(data []byte) (int, error) {
	return s.stdin.Write(data)
}

func (s *Session) Read(data []byte) (int, error) {
	return s.stdout.Read(data)
}

func (s *Session) Close() error {
	close(s.done)
	if err := s.pty.Close(); err != nil {
		return err
	}

	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}

	return nil
}

func (s *Session) handleWindowResize() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	for {
		select {
		case <-s.done:
			return
		case <-sig:
			_ = pty.InheritSize(os.Stdin, s.pty)
		}
	}
}

func (s *Session) forwardSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-s.done:
			return
		case sig := <-signals:
			if s.cmd != nil && s.cmd.Process != nil {
				_ = s.cmd.Process.Signal(sig)
			}
		}
	}
}
