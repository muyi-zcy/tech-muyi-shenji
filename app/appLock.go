package app

import (
	"sync"
)

// AppLock 按 appCode 分片互斥；使用 sync.Map 降低全局锁竞争。条目随应用数量增长，通常与应用规模同阶。
type AppLock struct {
	locks sync.Map // string -> *sync.RWMutex
}

var (
	appLockInstance *AppLock
	once            sync.Once
)

func GetAppLock() *AppLock {
	once.Do(func() {
		appLockInstance = &AppLock{}
	})
	return appLockInstance
}

func (al *AppLock) getMutex(appcode string) *sync.RWMutex {
	v, _ := al.locks.LoadOrStore(appcode, &sync.RWMutex{})
	return v.(*sync.RWMutex)
}

func (al *AppLock) RLock(appcode string) {
	al.getMutex(appcode).RLock()
}

func (al *AppLock) RUnlock(appcode string) {
	if v, ok := al.locks.Load(appcode); ok {
		v.(*sync.RWMutex).RUnlock()
	}
}

func (al *AppLock) Lock(appcode string) {
	al.getMutex(appcode).Lock()
}

func (al *AppLock) Unlock(appcode string) {
	if v, ok := al.locks.Load(appcode); ok {
		v.(*sync.RWMutex).Unlock()
	}
}
