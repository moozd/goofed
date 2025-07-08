package pty

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
)

type Session struct {
	cmd      *exec.Cmd
	fd       *os.File
	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewSession(ctx context.Context, shell string, args ...string) (*Session, error) {
	ctx, cancel := context.WithCancel(ctx)

	cmd := exec.Command(shell, args...)
	fd, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	session := &Session{
		cmd:      cmd,
		fd:       fd,
		ctx:      ctx,
		cancelFn: cancel,
	}

	go session.startWindowResizer()
	go session.startSignalForwarder()

	return session, nil
}

func (s *Session) Write(data []byte) (int, error) {
	return s.fd.Write(data)
}

func (s *Session) Read(data []byte) (int, error) {
	return s.fd.Read(data)
}

func (s *Session) Close() error {
	s.cancelFn()
	if err := s.fd.Close(); err != nil {
		return err
	}

	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}

	return nil
}

func (s *Session) startWindowResizer() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-sig:
			_ = pty.InheritSize(os.Stdin, s.fd)
		}
	}
}

func (s *Session) startSignalForwarder() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("Session: signal forwarder [stopped]")
			return
		case sig := <-signals:
			if s.cmd != nil && s.cmd.Process != nil {
				_ = s.cmd.Process.Signal(sig)
			}
		}
	}
}
