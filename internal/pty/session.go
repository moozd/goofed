package pty

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type Session struct {
	cmd    *exec.Cmd
	fd     *os.File
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSession(c context.Context, shell string, args ...string) (*Session, error) {
	ctx, cancel := context.WithCancel(c)

	cmd := exec.Command(shell, args...)
	fd, err := pty.Start(cmd)

	session := &Session{
		cmd:    cmd,
		fd:     fd,
		ctx:    ctx,
		cancel: cancel,
	}

	if err != nil {
		return nil, err
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
	s.cancel()
	if err := s.fd.Close(); err != nil {
		return err
	}

	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}

	return nil
}

func (s *Session) startWindowResizer() {
	signals := getResizeSignals()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-signals:
			_ = pty.InheritSize(os.Stdin, s.fd)
		}
	}
}

func (s *Session) startSignalForwarder() {
	signals := getStopSignals()
	for {
		select {
		case <-s.ctx.Done():
			return
		case signal := <-signals:
			if s.cmd != nil && s.cmd.Process != nil {
				_ = s.cmd.Process.Signal(signal)
			}
		}
	}
}
