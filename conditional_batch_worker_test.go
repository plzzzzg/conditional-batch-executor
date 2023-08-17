package conditiaonalbatchworker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAnd(t *testing.T) {
	size := 100
	lastRunTime := time.Now()
	w := New(func(items []*Item) (map[string]interface{}, error) {
		if len(items) < size || time.Since(lastRunTime).Seconds() < 2 {
			t.Fail()
		}
		m := map[string]interface{}{}
		for _, i := range items {
			m[i.Key] = i.Content
		}
		lastRunTime = time.Now()
		return m, nil
	}, And(Size(size), Interval(time.Second*2)))

	wg := sync.WaitGroup{}
	for i := 0; i < size*3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resultCh, err := w.Submit(fmt.Sprintf("%d", i), i)
			if err != nil {
				t.Error(err)
			}
			<-resultCh
		}(i)
	}
	wg.Wait()
}

func TestInterval(t *testing.T) {
	size := 100
	lastRunTime := time.Now()
	w := New(func(items []*Item) (map[string]interface{}, error) {
		if time.Since(lastRunTime).Seconds() < 2 {
			t.Fail()
		}
		m := map[string]interface{}{}
		for _, i := range items {
			m[i.Key] = i.Content
		}
		lastRunTime = time.Now()
		return m, nil
	}, Interval(time.Second*2))

	wg := sync.WaitGroup{}
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			resultCh, err := w.Submit(fmt.Sprintf("%d", i), i)
			if err != nil {
				t.Error(err)
			}
			<-resultCh
		}(i)
	}
	wg.Wait()
}

func TestSize(t *testing.T) {
	size := 100
	w := New(func(items []*Item) (map[string]interface{}, error) {
		if len(items) != size {
			t.Fail()
		}
		m := map[string]interface{}{}
		for _, i := range items {
			m[i.Key] = i.Content
		}
		return m, nil
	}, Size(size))

	wg := sync.WaitGroup{}
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			resultCh, err := w.Submit(fmt.Sprintf("%d", i), i)
			if err != nil {
				t.Error(err)
			}
			<-resultCh
		}(i)
	}
	wg.Wait()
}

func TestClose(t *testing.T) {
	size := 100
	w := New(func(items []*Item) (map[string]interface{}, error) {
		if len(items) != size {
			t.Fail()
		}
		m := map[string]interface{}{}
		for _, i := range items {
			m[i.Key] = i.Content
		}
		return m, nil
	}, Size(size))

	cnt := atomic.Int32{}
	wg := sync.WaitGroup{}
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			resultCh, err := w.Submit(fmt.Sprintf("%d", i), i)
			if err == nil {
				return
			}
			<-resultCh
			cnt.Add(1)
		}(i)
	}

	wg.Wait()
	if cnt.Load() == int32(size) {
		t.Fail()
	}
}
