package conditiaonalbatchexecutor

import (
	"errors"
	"sync"
	"time"
)

type Condition func(executor *Executor) bool

type Executor struct {
	queue           chan *Item
	readyTaskList   []*Item
	doFnc           func([]*Item) (map[string]interface{}, error)
	conditions      []Condition
	lastExecuteTime time.Time
	results         sync.Map
	isClose         bool
	waitCloseSign   chan struct{}
}

type Item struct {
	Key     string
	Content interface{}
}

func New(doFnc func([]*Item) (map[string]interface{}, error), conditions ...Condition) *Executor {
	worker := &Executor{doFnc: doFnc, conditions: conditions}
	worker.lastExecuteTime = time.Now()
	worker.queue = make(chan *Item, 1000)
	worker.results = sync.Map{}
	worker.waitCloseSign = make(chan struct{})
	worker.readyTaskList = make([]*Item, 0, 100)
	go worker.exec()
	return worker
}

func (w *Executor) Submit(taskID string, item interface{}) (<-chan interface{}, error) {
	if w.isClose {
		return nil, errors.New("worker is closed")
	}
	w.queue <- &Item{Content: item, Key: taskID}
	ch := make(chan interface{})
	if store, loaded := w.results.LoadOrStore(taskID, ch); !loaded {
		if v, ok := store.(chan interface{}); ok {
			return v, nil
		}
		w.results.Store(taskID, ch)
	}
	return ch, nil
}

func (w *Executor) Size() int {
	return len(w.readyTaskList)
}

func Size(i int) Condition {
	return func(worker *Executor) bool {
		return worker.Size() >= i
	}
}

func Interval(duration time.Duration) Condition {
	return func(worker *Executor) bool {
		return worker.lastExecuteTime.IsZero() || worker.lastExecuteTime.Add(duration).Before(time.Now())
	}
}

func And(first Condition, others ...Condition) Condition {
	return func(worker *Executor) bool {
		if !first(worker) {
			return false
		}
		for _, fn := range others {
			if !fn(worker) {
				return false
			}
		}
		return true
	}
}

func (w *Executor) exec() {
	for {
		select {
		case task := <-w.queue:
			w.readyTaskList = append(w.readyTaskList, task)
		default:
			break
		}

		runnable := false
		for i := range w.conditions {
			if w.conditions[i](w) {
				runnable = true
				break
			}
		}
		if !runnable || w.Size() == 0 {
			if w.Size() == 0 && w.isClose && len(w.queue) == 0 {
				w.waitCloseSign <- struct{}{}
				break
			}
			time.Sleep(time.Millisecond * 10)
			if !w.isClose {
				continue
			}
		}
		results, err := w.doFnc(w.readyTaskList)
		for itemID, result := range results {
			if ch, ok := w.results.Load(itemID); ok {
				if c, ok := ch.(chan interface{}); ok {
					if result == nil || err != nil {
						close(c)
					} else {
						c <- result
						w.results.Delete(itemID)
					}
				}
			}
		}
		w.readyTaskList = make([]*Item, 0, 100)
		w.lastExecuteTime = time.Now()
	}
}

func (w *Executor) Close() {
	w.isClose = true
	<-w.waitCloseSign
}
