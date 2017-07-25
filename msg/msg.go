package msg

/*
Message is the message object
*/
type Message struct {
	I           int64
	StartAtSec  int64
	StartAtNSec int64
	Payload     [512]byte
}
