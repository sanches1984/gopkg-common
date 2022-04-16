package closer

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

// GlobalCloser automatically calls when app is terminating
var globalCloser = New()

// Add adds `func() error` callback to the globalCloser
func Add(name string, f func() error) {
	globalCloser.Add(name, f)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs map[string]func() error
}

// New returns new Closer, if []os.Signal is specified Closer will automatically call CloseAll when one of signals is received from OS
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{}), funcs: make(map[string]func() error)}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

func (c *Closer) Add(name string, f func() error) {
	defer c.Unlock()
	c.Lock()
	if _, ok := c.funcs[name]; ok {
		panic("Closer " + name + " already used")
	}
	c.funcs[name] = f
}

func (c *Closer) Wait() {
	select {
	case <-c.done:
	}
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.Unlock()

		var wg sync.WaitGroup
		wg.Add(len(funcs))
		for name, fn := range funcs {
			go func(name string, fn func() error) {
				defer wg.Done()
				if err := fn(); err != nil {
					fmt.Printf("Error on close %s: %+v", name, err)
				}
			}(name, fn)
		}

		wg.Wait()
	})
}
