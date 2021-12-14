package apiserver

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/westcoastcode-se/gitgo/api"
	"html"
	"io/ioutil"
	"net/http"
	"time"
)

// Client is a type which is used when calling the API server
type Client struct {
	Address    string
	httpClient *http.Client
}

// FindUserUsingPublicKey fetches a user that has the supplied public key registered
func (c *Client) FindUserUsingPublicKey(fingerprint string) (*api.User, error) {
	// Get user that has the supplied user fingerprint
	resp, err := c.httpClient.Get(c.Address + "/api/v1/users?fingerprint=" + html.EscapeString(fingerprint))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	user := &api.User{}
	err = json.Unmarshal(data, user)
	return user, err
}

// NewClient creates a new https TLS client used when communicating with the API server
func NewClient(address string, certPath string, keyPath string, caPath string,
	insecureSkipVerify bool) (*Client, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Prepare client-side certificate
	var tlsConfig *tls.Config
	if insecureSkipVerify {
		tlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: insecureSkipVerify,
		}
	} else {
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
	}

	// Configure the actual https client
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    20,
		IdleConnTimeout: 5 * time.Minute,
	}
	client := &Client{
		Address: address,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   5000 * time.Millisecond,
		},
	}
	return client, nil
}
