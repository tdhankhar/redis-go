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
	return store[key]
}
