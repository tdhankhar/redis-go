package core

import (
	"log"
	"time"
)

const sampleLimit int = 20
const fractionThreshold float32 = 0.25

func deleteSample() float32 {
	deletedCount, i := 0, 0
	for key, obj := range store {
		if i >= sampleLimit {
			break
		}
		if obj.ExpiresAt != -1 {
			if obj.ExpiresAt <= time.Now().UnixMilli() {
				deletedCount += del(key)
			}
			i++
		}
	}
	return float32(deletedCount) / float32(sampleLimit)
}

func DeleteExpiredKeys() {
	for {
		deletedFraction := deleteSample()
		if deletedFraction < fractionThreshold {
			break;
		}
	}
	log.Println("expired keys deleted. total keys:", len(store))
}