package donut

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	CloudflareHost = "cloudflare-dns.com"
	GoogleHost     = "dns.google"
)

type Resolver struct {
	Host   string
	debug  bool
	client *http.Client
}

func New(host string, opts ...option) *Resolver {
	r := &Resolver{Host: host}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Resolver) Lookup(q Question) ([]Record, error) {
	question := encodeMessage([]Question{q})

	if r.debug {
		fmt.Printf("question: % 02x\n", question)
	}

	msg, err := r.lookup(question)
	if err != nil {
		return nil, err
	}

	if r.debug {
		fmt.Printf("answer: % 02x\n", msg)
	}

	return msg.parseMessage()
}

func (r *Resolver) LookupRaw(q []byte) ([]byte, error) {
	if r.debug {
		fmt.Printf("question: % 02x\n", q)
	}

	msg, err := r.lookup(q)
	if err != nil {
		return nil, err
	}

	if r.debug {
		fmt.Printf("answer: % 02x\n", msg)
	}

	return msg.buf, nil
}

func (r *Resolver) lookup(query []byte) (message, error) {
	url := "https://" + r.Host + "/dns-query"
	body := bytes.NewBuffer(query)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return message{}, err
	}

	req.Header.Set("Accept", "application/dns-message")
	req.Header.Set("Content-Type", "application/dns-message")

	if r.client == nil {
		r.client = http.DefaultClient
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return message{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return message{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return message{}, err
	}

	return message{buf}, nil
}
