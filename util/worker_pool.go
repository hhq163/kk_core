package util

import (
	"context"
	"sync"

	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
)

type Work func(c interface{})

type WorkerPool struct {
	work chan Work
	wg   sync.WaitGroup
}

func NewWorkerPool(maxGoroutines int, c interface{}) *WorkerPool {
	p := WorkerPool{
		work: make(chan Work),
	}
	if c == nil {
		p.wg.Add(maxGoroutines)
		for i := 0; i < maxGoroutines; i++ {
			go func() {
				for w := range p.work {
					w(nil)
				}
				p.wg.Done()
			}()
		}
	}else if client, ok := c.(*redis.ClusterClient); ok { //redis集群连接池
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
	} else if client, ok := c.(*redis.Client); ok { //redis连接池
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
	}

	return &p
}

func (p *WorkerPool) Run(w Work) {
	p.work <- w
}

func (p *WorkerPool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
