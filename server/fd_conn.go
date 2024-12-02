package server

import "syscall"

type fdConn struct {
	fd int
}

func (f fdConn) Read(bytes []byte) (int, error) {
	return syscall.Read(f.fd, bytes)
}

func (f fdConn) Write(bytes []byte) (int, error) {
	return syscall.Write(f.fd, bytes)
}
