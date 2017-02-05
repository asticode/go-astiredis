package astiredis

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/rs/xlog"
	"gopkg.in/redis.v5"
)

// Client represents a client
type Client struct {
	client *redis.Client
	Logger xlog.Logger
	prefix string
}

// New returns a client based on a configuration
func New(c Configuration) *Client {
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr: c.Addr,
		}),
		Logger: xlog.NopLogger,
		prefix: c.Prefix,
	}
}

// key builds a key with the prefix
func (c Client) key(k string) string {
	if len(c.prefix) == 0 {
		return k
	}
	return c.prefix + "." + k
}

// Del deletes a key
func (c Client) Del(k string) error {
	return c.client.Del(c.key(k)).Err()
}

// Get gets a value
func (c Client) Get(k string, v interface{}) error {
	b, err := c.client.Get(c.key(k)).Bytes()
	if err != nil {
		return err
	}

	// Decode
	return gob.NewDecoder(bytes.NewReader(b)).Decode(v)
}

// Set sets a value
func (c Client) Set(k string, v interface{}, ttl time.Duration) error {
	// Encode
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(v)
	if err != nil {
		return err
	}

	// Set
	return c.client.Set(c.key(k), buf.Bytes(), ttl).Err()
}