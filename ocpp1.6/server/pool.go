package server

import (
	"reflect"
	"sync"
)

type object interface {
	init(reflect.Type)
	get(t reflect.Type) interface{}
	put(t reflect.Type, x interface{})
}

// reset defines Reset method for pooled object.
type Reset interface {
	Reset()
}

var reflectInstance = func(t reflect.Type) interface{} {
	var argv reflect.Value
	if t.Kind() == reflect.Ptr {
		argv = reflect.New(t.Elem())
	} else {
		argv = reflect.New(t)
	}
	return argv.Interface()
}

type ocppTypePools struct {
	mu    sync.RWMutex
	pools map[reflect.Type]*sync.Pool
	New   func(t reflect.Type) interface{}
}

func (p *ocppTypePools) init(t reflect.Type) {
	tp := &sync.Pool{}
	tp.New = func() interface{} {
		return p.New(t)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pools[t] = tp
}

func (p *ocppTypePools) put(t reflect.Type, x interface{}) {
	if o, ok := x.(Reset); ok {
		o.Reset()
	}
	p.mu.RLock()
	pool := p.pools[t]
	p.mu.RUnlock()
	pool.Put(x)
}

func (p *ocppTypePools) get(t reflect.Type) interface{} {
	p.mu.RLock()
	pool := p.pools[t]
	p.mu.RUnlock()
	return pool.Get()
}

type ocppType struct{}

func (p *ocppType) init(t reflect.Type) {
}

func (p *ocppType) get(t reflect.Type) interface{} {
	return reflectInstance(t)
}
func (p *ocppType) put(t reflect.Type, x interface{}) {
	if o, ok := x.(Reset); ok {
		o.Reset()
	}
	return
}

func get(t reflect.Type) interface{} {
	return options.object.get(t)
}
func put(t reflect.Type, x interface{}) {
	options.object.put(t, x)
}

func SupportObjectPool(support bool) opt {
	return func(o *option) {
		switch support {
		case true:
			o.object = &ocppTypePools{
				pools: make(map[reflect.Type]*sync.Pool),
				New:   reflectInstance,
			}
		default:
			o.object = &ocppType{}
		}
	}
}

/*pool used for epoll trigger task*/
type TaskFunc func(interface{}) error
type Task struct {
	RunFunc TaskFunc
	RunArg  interface{}
}

var taskPool = sync.Pool{New: func() interface{} { return new(Task) }}

func GetTask() *Task {
	return taskPool.Get().(*Task)
}
func PutTask(task *Task) {
	task.RunFunc, task.RunArg = nil, nil
	taskPool.Put(task)
}
