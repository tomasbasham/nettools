/**
 * This is an implementation of the DNS wire format as defined in
 * https://datatracker.ietf.org/doc/html/rfc1035
 */
package donut

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type RData interface{}

type Question struct {
	FQDN  string
	Type  RecordType
	Class RecordClass
}

type Record struct {
	Name  string      `json:"name"`
	Type  RecordType  `json:"type"`
	Class RecordClass `json:"class"`
	TTL   uint32      `json:"ttl"`
	Data  RData       `json:"data"`
}

func encodeMessage(q []Question) []byte {
	message := bytes.NewBuffer(nil)

	// All communications inside of the domain protocol are carried in a single
	// format called a message.  The top level format of message is divided into 5
	// sections (some of which are empty in certain cases) shown below:
	//
	//  +---------------------+
	//  |        Header       |
	//  +---------------------+
	//  |       Question      | the question for the name server
	//  +---------------------+
	//  |        Answer       | RRs answering the question
	//  +---------------------+
	//  |      Authority      | RRs pointing toward an authority
	//  +---------------------+
	//  |      Additional     | RRs holding additional information
	//  +---------------------+
	//
	// The header section is always present.  The header includes fields that
	// specify which of the remaining sections are present, and also specify
	// whether the message is a query or a response, a standard query or some
	// other opcode, etc.

	// Header Section

	// The header contains the following fields:
	//
	//                                  1  1  1  1  1  1
	//    0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                      ID                       |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                    QDCOUNT                    |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                    ANCOUNT                    |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                    NSCOUNT                    |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                    ARCOUNT                    |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

	// ID is a 16 bit identifier assigned by the program. We'll set this to 0
	// but it could be any value.
	message.Write([]byte{0, 0})

	// QR is a 1 bit field that specifies whether this message is a query (0) or
	// a response (1). We're sending a query so this is 0.
	message.WriteByte((0 << 7) | (0 << 3) | (0 << 1) | 1)

	// Opcode is a 4 bit field that specifies kind of query in this message. This
	// is 0 for a standard query.
	message.WriteByte(1 << 4)

	// QDCOUNT is a 16 bit field that specifies the number of entries in the
	// question section. We're sending a single question so this is 1.
	message.Write(binaryBigEndianUint16(uint16(len(q))))

	// ANCOUNT is a 16 bit field that specifies the number of resource records in
	// the answer section. We're not sending any answers so this is 0.
	message.Write([]byte{0, 0})

	// NSCOUNT is a 16 bit field that specifies the number of name server resource
	// records in the authority records section. We're not sending any authority
	// records so this is 0.
	message.Write([]byte{0, 0})

	// ARCOUNT is a 16 bit field that specifies the number of resource records in
	// the additional records section. We're not sending any additional records so
	// this is 0.
	message.Write([]byte{0, 0})

	// Question Section
	for _, question := range q {
		message.Write(encodeQuestion(question))
	}

	return message.Bytes()
}

func encodeQuestion(q Question) []byte {
	question := bytes.NewBuffer(nil)

	// The question section is used to carry the "question" in most queries, i.e.,
	// the parameters that define what is being asked. The section contains the
	// following fields:
	//
	//                                  1  1  1  1  1  1
	//    0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                                               |
	//  /                     QNAME                     /
	//  /                                               /
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                     QTYPE                     |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                     QCLASS                    |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

	// QNAME is a domain name represented as a sequence of labels, where each
	// label consists of a length octet followed by that number of octets. The
	// domain name terminates with the zero length octet for the null label of the
	// root. We need to convert the FQDN into this format.

	labels := strings.Split(q.FQDN, ".")
	for _, label := range labels {
		question.WriteByte(byte(len(label)))
		question.WriteString(label)
	}
	question.WriteByte(0)

	// The QTYPE field specifies the type of the query.
	question.Write(binaryBigEndianUint16(uint16(q.Type)))

	// The QCLASS field specifies the class of the query.
	question.Write(binaryBigEndianUint16(uint16(q.Class)))

	return question.Bytes()
}

func binaryBigEndianUint16(v uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return b
}

type message struct {
	buf []byte
}

func (m *message) parseMessage() ([]Record, error) {
	qdcount := binary.BigEndian.Uint16(m.buf[4:6])
	ancount := binary.BigEndian.Uint16(m.buf[6:8])

	// Skip the header section and move to the question section.
	buf := m.buf[12:]

	// Loop through the number of questions, we're only interested in the answers.
	var i uint16
	for i = 0; i < qdcount; i++ {
		_, offset := m.decodeQuestion(buf)
		buf = buf[offset:]
	}

	answers := make([]Record, ancount)
	for i = 0; i < ancount; i++ {
		answer, offset := m.decodeAnswer(buf)
		answers[i] = answer
		buf = buf[offset:]
	}

	return answers, nil
}

func (m *message) decodeQuestion(q []byte) (Question, int) {
	// The question section is used to carry the "question" in most queries, i.e.,
	// the parameters that define what is being asked. The section contains the
	// following fields:
	//
	//                                  1  1  1  1  1  1
	//    0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                                               |
	//  /                     QNAME                     /
	//  /                                               /
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                     QTYPE                     |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                     QCLASS                    |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

	// QNAME is a domain name represented as a sequence of labels, where each
	// label consists of a length octet followed by that number of octets. The
	// domain name terminates with the zero length octet for the null label of the
	// root. We need to convert the FQDN into this format.
	name, offset := m.parseName(q)

	// The QTYPE field specifies the type of the query.
	qtype := binary.BigEndian.Uint16(q[offset : offset+2])

	// The QCLASS field specifies the class of the query.
	qclass := binary.BigEndian.Uint16(q[offset+2 : offset+4])

	return Question{
		FQDN:  name,
		Type:  RecordType(qtype),
		Class: RecordClass(qclass),
	}, offset + 4
}

func (m *message) decodeAnswer(a []byte) (Record, int) {
	// The answer section is used to carry the "answer" in response to a
	// query. The answer section contains the following fields:
	//
	//                                  1  1  1  1  1  1
	//    0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                                               |
	//  /                                               /
	//  /                      NAME                     /
	//  |                                               |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                      TYPE                     |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                     CLASS                     |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                      TTL                      |
	//  |                                               |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	//  |                   RDLENGTH                    |
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
	//  /                     RDATA                     /
	//  /                                               /
	//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

	// NAME is a domain name to which this resource record pertains.
	name, offset := m.parseName(a)

	// TYPE is two octets containing one of the RR type codes. This field
	// specifies the meaning of the data in the RDATA field.
	qtype := binary.BigEndian.Uint16(a[offset : offset+2])

	// CLASS is two octets which specify the class of the data in the RDATA field.
	qclass := binary.BigEndian.Uint16(a[offset+2 : offset+4])

	// TTL is a 32 bit unsigned integer that specifies the time interval that the
	// resource record may be cached before it should be discarded. Zero values
	// are interpreted to mean that the RR can only be used for the transaction
	// in progress, and should not be cached.
	ttl := binary.BigEndian.Uint32(a[offset+4 : offset+8])

	// RDLENGTH is an unsigned 16 bit integer that specifies the length in octets
	// of the RDATA field.
	rdlength := binary.BigEndian.Uint16(a[offset+8 : offset+10])

	// RDATA is a variable length string of octets that describes the resource.
	// The format of this information varies according to the TYPE and CLASS of
	// the resource record.
	rdata := a[offset+10 : offset+10+int(rdlength)]

	return Record{
		Name:  name,
		Type:  RecordType(qtype),
		Class: RecordClass(qclass),
		TTL:   ttl,
		Data:  rdata,
	}, offset + 10 + int(rdlength)
}

func (m *message) parseName(b []byte) (string, int) {
	var name string
	var offset int

	for {
		length := int(b[offset])
		if length == 0 {
			offset++
			break
		}

		// If the two most significant bits of the first octet are set then
		// compression is being used for the domain name. The next 14 bits
		// represent the offset from the start of the message where the domain
		// name is stored.
		if length&0xC0 == 0xC0 {
			ptr := binary.BigEndian.Uint16([]byte{b[offset], b[offset+1]}) & 0x3FFF
			name, _ = m.parseName(m.buf[ptr:])
			return name, 2
		} else {
			offset++
			name += string(b[offset:offset+length]) + "."
			offset += length
		}
	}

	return name, offset
}
