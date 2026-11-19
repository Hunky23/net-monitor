package domain

const (
	DnsTypeA    = "A"
	DnsTypeAAAA = "AAAA"
	DnsTypeSVCB = "SVCB"
)

type Dns struct {
	ID        int64
	ProcessID int64
	Domain    string
	Type      string
	Value     string
}
