package zaia

import (
	"encoding/json"
	"sync"
)

type Task[T any] struct {
	taskLocker     sync.RWMutex
	TaskMapList    map[string]*T
	ProjectName    string
	StatusCallback func(taskId string, t *T)
}

func NewTask[T any]() *Task[T] {
	return &Task[T]{
		taskLocker:  sync.RWMutex{},
		TaskMapList: make(map[string]*T),
		ProjectName: "",
	}
}

func (tc *Task[T]) Set(taskId string, t *T) {
	tc.taskLocker.Lock()
	defer tc.taskLocker.Unlock()
	tc.TaskMapList[taskId] = t
	tc.Dump()
	if tc.StatusCallback != nil {
		go func() {
			tc.StatusCallback(taskId, t)
		}()
	}
}

func (tc *Task[T]) Get(taskId string) *T {
	tc.taskLocker.Lock()
	defer tc.taskLocker.Unlock()
	v, _ := tc.TaskMapList[taskId]
	if v != nil {
		return &*v
	}
	return nil
}

func (tc *Task[T]) Find(callback func(t *T) bool) (string, *T) {
	tc.taskLocker.Lock()
	defer tc.taskLocker.Unlock()
	for k, v := range tc.TaskMapList {
		if callback(v) {
			return k, &*v
		}
	}
	return "", nil
}

func (tc *Task[T]) Remove(taskId string) {
	tc.taskLocker.Lock()
	defer tc.taskLocker.Unlock()
	v, _ := tc.TaskMapList[taskId]
	if v != nil {
		delete(tc.TaskMapList, taskId)
	}
	tc.Dump()
}

func (tc *Task[T]) Dump() {
	if len(tc.ProjectName) > 0 {
		DumpInterface(tc.ProjectName, tc.TaskMapList)
	}
}

func (tc *Task[T]) Load() {
	if len(tc.ProjectName) > 0 {
		v := make(map[string]*T)
		b := ReadFileToByte(tc.ProjectName)
		if b != nil {
			json.Unmarshal(b, &v)
			tc.TaskMapList = v
		}
	}
}
