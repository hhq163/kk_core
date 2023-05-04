package util

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hhq163/logger"
	"github.com/smallnest/chanx"
)

//ordered goroutine pool
type OrderWorkers struct {
	works  []*OrderWorker
	Length int32
	wg     sync.WaitGroup
	log    *logger.Logger
}

/**
* 	goNum: goroutine num
* 	chanLen: the init length of goroutine chan
 */
func NewOrderWorkers(goNum int32, chanLen int32, clog *logger.Logger) *OrderWorkers {
	if goNum <= 0 {
		return nil
	}

	w := &OrderWorkers{
		works:  make([]*OrderWorker, goNum),
		Length: goNum,
		log:    clog,
	}
	for k := range w.works {
		w.wg.Add(1)
		w.works[k] = NewOrderWorker(chanLen)
		SafeGo(w.works[k].Run, w.log, &w.wg)
	}
	return w
}

func (o *OrderWorkers) Push(id int32, f func()) error {
	i := id % o.Length
	return o.works[i].Push(f)
}

//Close the worker pool
func (o *OrderWorkers) Close() {
	o.wg.Wait()
	if len(o.works) > 0 {
		for k := range o.works {
			o.works[k].Close()
		}

	}
}

type OrderWorker struct {
	taskChan *chanx.UnboundedChan
	closed   bool
}

func NewOrderWorker(chanLen int32) *OrderWorker {
	w := &OrderWorker{
		taskChan: chanx.NewUnboundedChan((int)(chanLen)),
		closed:   false,
	}
	return w
}

func (w *OrderWorker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		f, ok := <-w.taskChan.Out
		if !ok {
			break
		}
		f.(func())()
	}
}

func (w *OrderWorker) Push(task interface{}) error {
	var err error
	if w.closed {
		return errors.New("taskChan closed")
	}
	if f, ok := task.(func()); ok {
		w.taskChan.In <- f
	} else {
		return fmt.Errorf("task=%v is not a fun", task)
	}

	if i, ok := task.(int); ok {
		if i == -1 {
			w.closed = true
			close(w.taskChan.In)
		}
	}

	return err
}

func (w *OrderWorker) Close() {
	w.Push(-1)
}
