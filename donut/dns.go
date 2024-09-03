package donut

type RecordType uint16

const (
	A     RecordType = 1
	NS    RecordType = 2
	CNAME RecordType = 5
	SOA   RecordType = 6
	PTR   RecordType = 12
	MX    RecordType = 15
	TXT   RecordType = 16
	AAAA  RecordType = 28
	SRV   RecordType = 33
)

type RecordClass uint16

const (
	IN RecordClass = 1
	CS RecordClass = 2
	CH RecordClass = 3
	HS RecordClass = 4
)
