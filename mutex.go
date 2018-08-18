package astiredis

import (
	"time"

	"context"

	"github.com/asticode/go-astitools/time"
	"github.com/pkg/errors"
)

// Mutex represents a distributed mutex
type Mutex struct {
	c       *Client
	ctx     context.Context
	key     string
	sleep   time.Duration
	timeout time.Duration
}

// MutexOptions represents mutex options
type MutexOptions struct {
	Key     string
	Sleep   time.Duration
	Timeout time.Duration
}

// NewMutex creates a new mutex
func (c *Client) NewMutex(ctx context.Context) *Mutex {
	return c.NewMutexWithOptions(ctx, MutexOptions{
		Key:   "distributed.mutex",
		Sleep: 500 * time.Millisecond,
	})
}

// NewMutexWithOptions creates a new mutex with options
func (c *Client) NewMutexWithOptions(ctx context.Context, o MutexOptions) *Mutex {
	return &Mutex{
		c:       c,
		ctx:     ctx,
		key:     o.Key,
		sleep:   o.Sleep,
		timeout: o.Timeout,
	}
}

// Lock lock the mutex
func (l *Mutex) Lock() (err error) {
	for {
		// Check context error
		if l.ctx.Err() != nil {
			err = errors.Wrap(err, "astiredis: context error")
			return
		}

		// Try to set the key
		var ok bool
		if ok, err = l.c.SetNX(l.key, true, l.timeout); err != nil {
			err = errors.Wrapf(err, "astiredis: setting mutex key %s failed", l.key)
			return
		}

		// Key already exists
		if !ok {
			astitime.Sleep(l.ctx, l.sleep)
			continue
		}
		return
	}
}

// Unlock unlocks the mutex
func (l *Mutex) Unlock() (err error) {
	if err = l.c.Del(l.key); err != nil {
		err = errors.Wrapf(err, "astiredis: deleting mutex key %s failed", l.key)
		return
	}
	return
}
