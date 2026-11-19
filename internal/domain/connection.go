package domain

// Транспортные протоколы соединений.
const (
	ConnectionTransportProtocolTCP = "TCP"
	ConnectionTransportProtocolUDP = "UDP"
)

const (
	ConnectionApplicationProtocolSSH   = "ssh"
	ConnectionApplicationProtocolHTTP1 = "http1"
	ConnectionApplicationProtocolHTTP2 = "http2"
	ConnectionApplicationProtocolHTTP3 = "http3"
	ConnectionApplicationProtocolDNS   = "dns"
	ConnectionApplicationProtocolOther = "other"
)

const (
	ConnectionStatusOpen   = "open"
	ConnectionStatusClosed = "closed"
	ConnectionStatusReset  = "reset"
)

type Connection struct {
	ID                  int64
	ProcessID           int64
	SourceIP            string
	SourcePort          int64
	DestinationIP       string
	DestinationPort     int64
	Domain              string
	BytesReceived       int64
	BytesSent           int64
	TransportProtocol   string
	ApplicationProtocol string
	Status              string
}
