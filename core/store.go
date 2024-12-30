package core

import "time"

type storeObj struct {
	Value interface{}
	ExpiresAt int64
}

var store map[string]*storeObj

func init() {
	store = make(map[string]*storeObj)
}

func createStoreObj(value interface{}, expiryMs int64) *storeObj {
	var expiresAt int64 = -1
	if expiryMs > 0 {
		expiresAt = time.Now().UnixMilli() + expiryMs
	}
	return &storeObj{
		Value: value,
		ExpiresAt: expiresAt,
	}
}

func put(key string, obj *storeObj) {
	store[key] = obj
}

func get(key string) *storeObj {
	obj := store[key]
	if obj != nil {
		if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
			delete(store, key)
			return nil
		}
	}
	return obj
}

func del(key string) int {
	if _, ok := store[key]; ok {
		delete(store, key)
		return 1
	}
	return 0
}
