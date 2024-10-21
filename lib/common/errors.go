package common

import (
	"sync"

	"go.uber.org/zap"
)

type ErrorDetail struct {
	Id    int    `json:"id"`
	Error string `json:"error"`
}

type MetricErrors struct {
	logger *zap.Logger
	mutex  sync.Mutex
	errors []ErrorDetail
}

func NewMetricErrors(logger *zap.Logger) *MetricErrors {
	return &MetricErrors{
		logger: logger,
	}
}

func (me *MetricErrors) Append(id int, err string, skipLogging ...bool) {
	if me.logger != nil && (len(skipLogging) == 0 || !skipLogging[0]) {
		me.logger.Error(err)
	}
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
