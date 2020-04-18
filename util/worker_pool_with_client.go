package util

import (
	"context"
	"sync"

	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
)

type Work func(c interface{})

type WorkerPoolWithClient struct {
	work chan Work
	wg   sync.WaitGroup
}

func NewWorkerPoolWithClient(maxGoroutines int, c interface{}) *WorkerPoolWithClient {
	p := WorkerPoolWithClient{
		work: make(chan Work),
	}
	if client, ok := c.(*redis.ClusterClient); ok {
		p.wg.Add(maxGoroutines)
		for i := 0; i < maxGoroutines; i++ {
			go func() {
				cClient := client.WithContext(context.Background())
				for w := range p.work {
					w(cClient)
				}
				p.wg.Done()
			}()
		}
	} else if client, ok := c.(*redis.Client); ok {
		p.wg.Add(maxGoroutines)
		for i := 0; i < maxGoroutines; i++ {
			go func() {
				cClient := client.WithContext(context.Background())
				for w := range p.work {
					w(cClient)
				}
				p.wg.Done()
			}()
		}
	} else if client, ok := c.(*mgo.Session); ok {
		p.wg.Add(maxGoroutines)
		for i := 0; i < maxGoroutines; i++ {
			go func() {
				cClient := client.Clone()
				for w := range p.work {
					w(cClient)
				}
				p.wg.Done()
			}()
		}
	} else {
		return nil
	}

	return &p
}

func (p *WorkerPoolWithClient) Run(w Work) {
	p.work <- w
}

func (p *WorkerPoolWithClient) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
