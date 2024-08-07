package telemetry

import "sync"

type ErrorDetail struct {
	Id    int    `json:"id"`
	Error string `json:"error"`
}

type MetricErrors struct {
	mutex  sync.Mutex
	errors []ErrorDetail
}

func (me *MetricErrors) Append(id int, err string) {
	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.errors = append(me.errors, ErrorDetail{Id: id, Error: err})
}

func (me *MetricErrors) Get() []ErrorDetail {
	me.mutex.Lock()
	defer me.mutex.Unlock()
	return me.errors
}

func (me *MetricErrors) Len() int {
	me.mutex.Lock()
	defer me.mutex.Unlock()
	return len(me.errors)
}
