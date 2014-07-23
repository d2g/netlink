package netlink

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"
	"syscall"
)

type Connection struct {
	socketID       int
	connectionType ConnectionType

	socket         int
	sequenceNumber uint32
	socketAddress  syscall.SockaddrNetlink

	handleFunc (func([]byte) error)
}

type ioreader struct {
	Connection int
}

func (s ioreader) Read(b []byte) (n int, err error) {
	nr, _, e := syscall.Recvfrom(s.Connection, b, 0)

	//2014-07-21 Note: This stops the reader creating a panic in bufio/bufio.go line 99
	if nr < 0 {
		return 0, os.NewSyscallError("recvfrom", e)
	} else {
		return nr, os.NewSyscallError("recvfrom", e)
	}
}

func GetNetlinkSocket(socketid int, connectiontype ConnectionType) Connection {
	netlinkConnection := Connection{}
	netlinkConnection.socketID = socketid
	netlinkConnection.connectionType = connectiontype
	return netlinkConnection
}

func (t *Connection) SocketID() int {
	return t.socketID
}

func (t *Connection) SetHandleFunc(fn func([]byte) error) {
	t.handleFunc = fn
}

func (t *Connection) HandleFunc() func([]byte) error {
	return t.handleFunc
}

func (t *Connection) Connect() error {
	var err error
	t.socket, err = syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, t.socketID)

	if err != nil {
		log.Printf("Cannot create netlink socket: %s\n", os.NewSyscallError("socket", err))
		return err
	}

	t.socketAddress = syscall.SockaddrNetlink{}
	t.socketAddress.Family = syscall.AF_NETLINK
	t.socketAddress.Pid = 0
	t.socketAddress.Groups = uint32(t.connectionType)

	err = syscall.Bind(t.socket, &t.socketAddress)
	if err != nil {
		log.Printf("Cannot bind netlink socket: %s", os.NewSyscallError("bind", err))
		syscall.Close(t.socket)
		return err
	}
	return nil
}

func (t *Connection) Read() ([]byte, error) {
	var readBuffer *bufio.Reader

	message := syscall.NetlinkMessage{}

	//IOreader
	reader := ioreader{}
	reader.Connection = t.socket //Socket!
	readBuffer = bufio.NewReader(reader)

	binary.Read(readBuffer, binary.LittleEndian, &message.Header)

	if message.Header.Len > syscall.NLMSG_HDRLEN {
		message.Data = make([]byte, message.Header.Len-syscall.NLMSG_HDRLEN)

		_, err := readBuffer.Read(message.Data)
		if err != nil {
			syscall.Close(t.socket)
			return []byte{}, err
		}

		return message.Data, nil
	} else {
		return []byte{}, nil
	}
}

func (t *Connection) Write(message []byte) error {

	netlinkMessage := syscall.NetlinkMessage{}

	t.sequenceNumber++

	netlinkMessage.Data = message
	netlinkMessage.Header.Seq = t.sequenceNumber
	netlinkMessage.Header.Pid = uint32(os.Getpid())

	messageBuffer := bytes.NewBuffer(nil)
	binary.Write(messageBuffer, binary.LittleEndian, netlinkMessage.Header)

	_, err := messageBuffer.Write(netlinkMessage.Data)
	if err != nil {
		log.Println("Error Writing to Message Buffer")
		return err
	}

	err = syscall.Sendto(t.socket, messageBuffer.Bytes(), 0, &t.socketAddress)
	if err != nil {
		return err
	}
	return nil
}

func (t *Connection) Close() error {
	return syscall.Close(t.socket)
}

func (t *Connection) ListenAndServe() error {
	if t.handleFunc == nil {
		return errors.New("HandleFunc Not Set")
	}

	if t.socket == 0 {
		err := t.Connect()
		defer t.Close()
		if err != nil {
			return err
		}
	}

	for {
		message, err := t.Read()
		if err != nil {
			return err
		} else {
			err = t.HandleFunc()(message)
			if err != nil {
				return err
			}
		}
	}

}
