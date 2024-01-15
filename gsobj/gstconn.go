package gsobj

type GSTConn struct {
	pack, unpack *GSTObj
}

func NewGSTConn(pack, unpack *GSTObj) *GSTConn {
	return &GSTConn{pack, unpack}
}
