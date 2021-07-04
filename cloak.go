package snake

import (
	"sync"
	"time"
)

// Cloak is the interface that wraps a ticker.
//
// Start starts the ticker.
//
// Tick wraps the ticker receive channel.
//
// Stop stops the ticker.
type Cloak interface {
	Start(d time.Duration)
	Tick() <-chan time.Time
	Stop()
}

// Ticker wrapper implementation
type DefaultCloak struct {
	ticker *time.Ticker
	wg     *sync.WaitGroup
}

// NewCloak returns a pointer to DefaultCloak. Start must be called
// in order to receive ticks.
//
// NewCloak initializes an internal sync.WaitGroup which allows to start
// the internal ticker on a later time.
func NewCloak() *DefaultCloak {
	var wg sync.WaitGroup
	wg.Add(1)
	return &DefaultCloak{nil, &wg}
}

// Start initializes the internal time.Ticker releasing the internal sync.WaitGroup.
func (c *DefaultCloak) Start(d time.Duration) {
	defer c.wg.Done()
	c.ticker = time.NewTicker(d)
}

// Tick returns the internal time.Ticker.C receive channel waiting for the cloak start.
func (c *DefaultCloak) Tick() <-chan time.Time {
	c.wg.Wait()
	return c.ticker.C
}

// Stop stops the internal ticker after waiting for the cloak to start.
func (c *DefaultCloak) Stop() {
	c.wg.Wait()
	c.ticker.Stop()
}
