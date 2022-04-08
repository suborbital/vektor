package vtest

import (
	"net"
	"net/http"
	"time"
)

type fakeConn struct {
	http.ResponseWriter
}

// Read reads data from the connection.
func (f *fakeConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

// Write writes data to the connection.
func (f *fakeConn) Write(b []byte) (n int, err error) {
	return f.ResponseWriter.Write(b)
}

// Close closes the connection.
func (f *fakeConn) Close() error {
	return nil
}

// LocalAddr returns the local network address, if known.
func (f *fakeConn) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr returns the remote network address, if known.
func (f *fakeConn) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
func (f *fakeConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
func (f *fakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
func (f *fakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}
