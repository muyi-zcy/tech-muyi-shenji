package app

import (
	"sync"
)

// AppLock 是一个全局可用的读写锁工具类
type AppLock struct {
	mu      sync.Mutex
	rwLocks map[string]*sync.RWMutex
}

var (
	appLockInstance *AppLock
	once            sync.Once
)

// GetAppLock 返回唯一的 AppLock 实例
func GetAppLock() *AppLock {
	once.Do(func() {
		appLockInstance = &AppLock{
			rwLocks: make(map[string]*sync.RWMutex),
		}
	})
	return appLockInstance
}

// RLock 获取读锁
func (al *AppLock) RLock(appcode string) {
	al.mu.Lock()
	if _, exists := al.rwLocks[appcode]; !exists {
		al.rwLocks[appcode] = &sync.RWMutex{}
	}
	al.mu.Unlock()

	al.rwLocks[appcode].RLock()
}

// RUnlock 释放读锁
func (al *AppLock) RUnlock(appcode string) {
	if rwLock, exists := al.rwLocks[appcode]; exists {
		rwLock.RUnlock()
	}
}

// Lock 获取写锁
func (al *AppLock) Lock(appcode string) {
	al.mu.Lock()
	if _, exists := al.rwLocks[appcode]; !exists {
		al.rwLocks[appcode] = &sync.RWMutex{}
	}
	al.mu.Unlock()

	al.rwLocks[appcode].Lock()
}

// Unlock 释放写锁
func (al *AppLock) Unlock(appcode string) {
	if rwLock, exists := al.rwLocks[appcode]; exists {
		rwLock.Unlock()
	}
}
