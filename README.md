# Conditional Batch Executor

[![Go Reference](https://pkg.go.dev/badge/github.com/plzzzzg/conditional-batch-executor#section-readme.svg)](https://pkg.go.dev/github.com/plzzzzg/conditional-batch-executor#section-readme)
[![Go Report Card](https://goreportcard.com/badge/github.com/plzzzzg/conditional-batch-executor)](https://goreportcard.com/report/github.com/plzzzzg/conditional-batch-executor)
[![Go](https://github.com/plzzzzg/conditional-batch-executor/actions/workflows/go.yml/badge.svg)](https://github.com/plzzzzg/conditional-batch-executor/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/plzzzzg/conditional-batch-executor/graph/badge.svg?token=KONGT3A0I8)](https://codecov.io/gh/plzzzzg/conditional-batch-executor)

A batch worker that collects tasks and executes them when conditions are met.

Caller can get the results asynchronously.

## Install

```shell
go get github.com/plzzzzg/conditional-batch-executor
```

## Examples

```go
// init
worker := conditiaonalbatchexecutor.New(func(tasks []interface{}) map[string]interface{} {
		m := map[string]interface{}{}
		// do something ...
		return m // key is the id from Submit
	}, conditiaonalbatchexecutor.Size(3), conditiaonalbatchexecutor.Interval(time.Second*2)) // execute if size of tasks >= 3 OR after 2 seconds since last execution  

// submit
resultReciever, err := worker.Submit(idStr, i)

// receive
result := <-resultReciever
```

## Supported Conditions

### Interval

```go
worker := conditiaonalbatchexecutor.New(func(tasks []interface{}) map[string]interface{} {
		m := map[string]interface{}{}
		return m 
	}, conditiaonalbatchexecutor.Interval(time.Second*2)) // execute every 2 seconds 
```

### Size

```go
worker := conditiaonalbatchexecutor.New(func(tasks []interface{}) map[string]interface{} {
		m := map[string]interface{}{}
		return m 
	}, conditiaonalbatchexecutor.Size(2)) // execute then buffer size reaches 2
```

### And


```go
worker := conditiaonalbatchexecutor.New(func(tasks []interface{}) map[string]interface{} {
		m := map[string]interface{}{}
		return m 
	}, conditiaonalbatchexecutor.And(Size(2), conditiaonalbatchexecutor.Interval(time.Second*2))) // execute then buffer size reaches 2 AND last execution happened more than 2 min ago  
```