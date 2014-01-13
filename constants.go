package netlink

type ConnectionType int

const (
	Unicast   ConnectionType = 0
	Broadcast ConnectionType = 1
)
