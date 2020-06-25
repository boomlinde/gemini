package client

import (
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
)

// UntrustedError is returned by Request when verification against a pinned certificate signature has failed
var UntrustedError error = errors.New("untrusted cert signature")

// TofuStore is an interface to the TOFU storage backend
type TofuStore interface {
	Get(string) (string, error)
	Pin(string, string) error
}

// Client represents a Gemini client
type Client struct {
	Tofu TofuStore
	// If Dialer is nil, tls.Dial will be used to connect
	Dialer *net.Dialer
	// Contains client certificates
	Certs []tls.Certificate
}

// NewClient returns a new Client with the default TofuStore pointing at the given path
func NewClient(tofuPath string) *Client {
	return &Client{
		Tofu: &tofustore{tofuPath},
	}
}

// Request will request the given URI, return its header and a ReadCloser.
// The caller is responsible for closing the ReadCloser.
func (c *Client) Request(uri string) (io.ReadCloser, error) {
	cfg := &tls.Config{}
	if c.Certs != nil {
		cfg.Certificates = c.Certs
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}
	signature, err := c.Tofu.Get(u.Hostname())
	// can't get a signature; request with full verification
	if err != nil {
		return c.request(u, cfg, nil)
	}

	// Otherwise, skip verify and only verify that the certs match the signature
	verify := func(certs []*x509.Certificate) error {
		serverSig := digest(certs)
		if serverSig != signature {
			return UntrustedError
		}
		return nil
	}

	cfg.InsecureSkipVerify = true
	return c.request(u, cfg, verify)
}

// Pin will pin a signature of the certificates of the given URI
func (c *Client) Pin(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("failed to parse URI: %w", err)
	}

	certs, err := c.getcerts(u)
	if err != nil {
		return fmt.Errorf("failed to get certs: %v", err)
	}

	sig := digest(certs)
	return c.Tofu.Pin(u.Hostname(), sig)
}

func (c *Client) dial(u *url.URL, config *tls.Config) (*tls.Conn, error) {
	port := u.Port()

	if port == "" {
		port = "1965"
	}

	addr := u.Hostname() + ":" + port
	if c.Dialer == nil {
		return tls.Dial("tcp", addr, config)
	}

	return tls.DialWithDialer(c.Dialer, "tcp", addr, config)
}

func (c *Client) request(u *url.URL, config *tls.Config, verify func([]*x509.Certificate) error) (io.ReadCloser, error) {
	conn, err := c.dial(u, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	connstate := conn.ConnectionState()
	if verify != nil {
		if err := verify(connstate.PeerCertificates); err != nil {
			return nil, fmt.Errorf("failed to verify pinned cert: %w", err)
		}
	}

	_, err = conn.Write([]byte(u.String() + "\r\n"))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("write request failed: %w", err)
	}

	return conn, nil
}

func (c *Client) getcerts(u *url.URL) ([]*x509.Certificate, error) {
	config := &tls.Config{InsecureSkipVerify: true}
	conn, err := c.dial(u, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	connstate := conn.ConnectionState()
	return connstate.PeerCertificates, nil
}

// Untrusted checks whether an error is due to either unknown root CA or failed pin verification
func Untrusted(err error) bool {
	var autherr x509.UnknownAuthorityError
	return errors.As(err, &autherr) || err == UntrustedError
}

// Invalid checks whether the error is due to the certificate being invalid
func Invalid(err error) bool {
	var invaliderr x509.CertificateInvalidError
	return errors.As(err, &invaliderr)
}

// digest calculates an opaque digest string of a certificate chain
func digest(certs []*x509.Certificate) string {
	total := sha1.New()
	for _, cert := range certs {
		total.Write(cert.Raw)
	}
	sum := total.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}
