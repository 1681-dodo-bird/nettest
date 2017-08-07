package msg

/*
Message is the message object
*/
type Message struct {
	Index       uint64
	StartAtSec  int64
	StartAtNSec int64
	Payload     [512]byte
}
