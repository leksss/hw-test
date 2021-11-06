package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var ErrInvalidConnection = fmt.Errorf("invalid connection")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type extremelyPrimitiveTelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &extremelyPrimitiveTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (s *extremelyPrimitiveTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", s.address, s.timeout)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	s.conn = conn
	return nil
}

func (s *extremelyPrimitiveTelnetClient) Send() error {
	if s.conn == nil {
		return ErrInvalidConnection
	}
	if _, err := io.Copy(s.conn, s.in); err != nil {
		return fmt.Errorf("send %w", err)
	}
	return nil
}

func (s *extremelyPrimitiveTelnetClient) Receive() error {
	if s.conn == nil {
		return ErrInvalidConnection
	}
	if _, err := io.Copy(s.out, s.conn); err != nil {
		return fmt.Errorf("send %w", err)
	}
	return nil
}

func (s *extremelyPrimitiveTelnetClient) Close() error {
	if s.conn == nil {
		return ErrInvalidConnection
	}
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("close %w", err)
	}
	return nil
}
