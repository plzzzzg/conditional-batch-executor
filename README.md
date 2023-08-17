# Conditional Batch Worker

[![Go Reference](https://pkg.go.dev/badge/github.com/plzzzzg/conditional-batch-worker#section-readme.svg)](https://pkg.go.dev/github.com/plzzzzg/conditional-batch-worker#section-readme)
[![Go Report Card](https://goreportcard.com/badge/github.com/plzzzzg/conditional-batch-worker)](https://goreportcard.com/report/github.com/plzzzzg/conditional-batch-worker)
[![Go](https://github.com/plzzzzg/conditional-batch-worker/actions/workflows/go.yml/badge.svg)](https://github.com/plzzzzg/conditional-batch-worker/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/plzzzzg/conditional-batch-worker/graph/badge.svg?token=KONGT3A0I8)](https://codecov.io/gh/plzzzzg/conditional-batch-worker)

A batch worker that collects tasks and executes them when conditions are met.

Caller can get the results asynchronously.

## Install

```shell
go install github.com/plzzzzg/condition-batch-worker
```

## Examples

```go
// init
worker := conditiaonalbatchworker.New(func(tasks []interface{}) map[string]interface{} {
		m := map[string]interface{}{}
		// do something ...
		return m // key is the id from Submit
	}, conditiaonalbatchworker.Size(3), conditiaonalbatchworker.Interval(time.Second*2)) // execute if size of tasks >= 3 OR after 2 seconds since last execution  

// submit
resultReciever, err := worker.Submit(idStr, i)

// receive
result := <-resultReciever
```

## Supported Conditions

### Interval

```go
worker := conditiaonalbatchworker.New(func(tasks []interface{}) map[string]interface{} {
		m := map[string]interface{}{}
		return m 
	}, conditiaonalbatchworker.Interval(time.Second*2)) // execute every 2 seconds 
```

### Size

```go
worker := conditiaonalbatchworker.New(func(tasks []interface{}) map[string]interface{} {
		m := map[string]interface{}{}
		return m 
	}, conditiaonalbatchworker.Size(2)) // execute then buffer size reaches 2
```

### And


```go
worker := conditiaonalbatchworker.New(func(tasks []interface{}) map[string]interface{} {
		m := map[string]interface{}{}
		return m 
	}, conditiaonalbatchworker.And(Size(2), conditiaonalbatchworker.Interval(time.Second*2))) // execute then buffer size reaches 2 AND last execution happened more than 2 min ago  
```