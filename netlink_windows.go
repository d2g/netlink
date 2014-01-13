package netlink

import (
	"errors"
)

type Connection struct {
}

func GetNetlinkSocket(socketid int, connectiontype ConnectionType) Connection {
	return Connection{}
}

func (t *Connection) SetHandleFunc(fn func([]byte) error) {

}

func (t *Connection) HandleFunc(fn func([]byte) error) {
	return
}

func (t *Connection) ListenAndServe() error {
	return errors.New("Netlink Is not Supported On Windows")
}

func (t *Connection) Connect() error {
	return errors.New("Netlink Is not Supported On Windows")
}

func (t *Connection) Read() ([]byte, error) {
	return make([]byte, 0), errors.New("Netlink Is not Supported On Windows")
}

func (t *Connection) Write(message []byte) error {
	return errors.New("Netlink Is not Supported On Windows")
}

func (t *Connection) Close() error {
	return errors.New("Netlink Is not Supported On Windows")
}
