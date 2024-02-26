package packet

type MethodRequest struct {
	Ver      byte   // Socks version
	NMethods byte   // Number of authentication methods supported
	Methods  []byte // Authentication methods 1-255
}

const (
	Ver byte = 0x05
)

type MethodResponse struct {
	Ver    byte // Socks version
	Method byte // Support method
}

func (r *MethodResponse) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, r.Ver)
	bytes = append(bytes, r.Method)
	return bytes
}

const (
	MethodNoAuth byte = 0x00 // No need auth
	MethodUPAuth byte = 0x02 // Username/Password
)

// If Method is MethodUPAuth
type NegotiationRequest struct {
	Ver    byte   // Socks version
	ULen   byte   // Username length
	UName  []byte // Username value
	PLen   byte   // Password length
	Passwd []byte // Password value
}

type NegotiationResponse struct {
	Ver    byte // Socks version
	Status byte // Auth result
}

func (r *NegotiationResponse) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, r.Ver)
	bytes = append(bytes, r.Status)
	return bytes
}

const (
	StatusSuccess byte = 0x00
	StatuFailure  byte = 0xFF // Any other value except 0x00
)

type DeliverRequest struct {
	Ver        byte    // Socks version
	Cmd        byte    // Command code
	Rsv        byte    // Spare field
	ATyp       byte    // Address type ipv4 or ipv6
	DstHostLen byte    // If ATyp is domain
	DstHost    []byte  // Desired destination address
	DstPort    [2]byte // Desired destination port
}

const (
	CmdConnect byte = 0x01
	CmdUDP     byte = 0x03
	Rsv        byte = 0x00 // Spare field
	ATypIPV4   byte = 0x01
	ATypDomain byte = 0x03
	ATypIPV6   byte = 0x04
)

type DeliverResponse struct {
	Ver        byte    // Socks version
	Rep        byte    // Result reply
	Rsv        byte    // Spare field
	ATyp       byte    // Address type ipv4 or ipv6
	BndHostLen byte    // If ATyp is domain
	BndHost    []byte  // Desired destination address
	BndPort    [2]byte // Desired destination port
}

func (r *DeliverResponse) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, r.Ver)
	bytes = append(bytes, r.Rep)
	bytes = append(bytes, r.Rsv)
	bytes = append(bytes, r.ATyp)
	if r.ATyp == ATypDomain {
		bytes = append(bytes, r.BndHostLen)
	}
	bytes = append(bytes, r.BndHost...)
	bytes = append(bytes, r.BndPort[:]...)
	return bytes
}

const (
	// Success for repling
	RepSuccess byte = 0x00
	// The server failure
	RepServerFailure byte = 0x01
	// The request not allowed
	RepNotAllowed byte = 0x02
	// The network unreachable
	RepNetworkUnreachable byte = 0x03
	// The host unreachable
	RepHostUnreachable byte = 0x04
	// The connection refused
	RepConnectionRefused byte = 0x05
	// The TTL expired
	RepTTLExpired byte = 0x06
	// The request command not supported
	RepCommandNotSupported byte = 0x07
	// The request address not supported
	RepAddressNotSupported byte = 0x08
	// Undefine
	UnknowMessage byte = 0x09
)

var ReqMessageMap map[byte]string = map[byte]string{
	RepSuccess:             "connection established successfully",
	RepServerFailure:       "socks5 server failure",
	RepNotAllowed:          "network unreachable",
	RepHostUnreachable:     "host unreachable",
	RepConnectionRefused:   "connection refused by destination host",
	RepTTLExpired:          "TTL expired",
	RepCommandNotSupported: "command not supported",
	RepAddressNotSupported: "address type not supported",
	UnknowMessage:          "unknown error",
}
