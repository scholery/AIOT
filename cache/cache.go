package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/bluele/gcache"
)

var VALUE_NOT_STRING_ERROR = errors.New("value is not string")
var VALUE_NOT_INT64_ERROR = errors.New("value is not int64")

var incrMutex = &sync.Mutex{}

// must init first
var cacheClient gcache.Cache
var cacheClient2 gcache.Cache

func Init() {
	cacheClient = gcache.New(300).ARC().Build()
	cacheClient2 = gcache.New(300).ARC().Build()
}

// client 2 start
func Set2(key, value interface{}, ttl time.Duration) error {
	return cacheClient2.SetWithExpire(key, value, ttl)
}

func Get2(key interface{}) (interface{}, error) {
	return cacheClient2.Get(key)
}

func Delete2(key string) bool {
	return cacheClient2.Remove(key)
}

// client 2 end

func Set(key, value interface{}, ttl time.Duration) error {
	return cacheClient.SetWithExpire(key, value, ttl)
}

func Get(key interface{}) (interface{}, error) {
	return cacheClient.Get(key)
}

func GetString(key interface{}) (string, error) {
	v, err := cacheClient.Get(key)
	if err != nil {
		return "", err
	}
	vs, ok := v.(string)
	if !ok {
		return "", VALUE_NOT_STRING_ERROR
	}
	return vs, err
}

func GetInt(key interface{}, defaultValue int) (int, error) {
	v, err := cacheClient.Get(key)
	if err == gcache.KeyNotFoundError {
		return defaultValue, nil
	} else if err != nil {
		return defaultValue, err
	}
	vs, ok := v.(int)
	if !ok {
		return defaultValue, VALUE_NOT_INT64_ERROR
	}
	return vs, err
}

func Delete(key string) bool {
	return cacheClient.Remove(key)
}

func GetKeys() []string {
	keys := cacheClient.Keys(false)
	var keyString = make([]string, 0)
	for _, key := range keys {
		keyString = append(keyString, key.(string))
	}
	return keyString
}

func Incr(key string, ttl time.Duration) error {
	incrMutex.Lock()
	defer incrMutex.Unlock()
	v, err := GetInt(key, 0)
	if err != nil {
		return err
	}
	return Set(key, v+1, ttl)
}

func SetWithNoExpire(key, value interface{}) error {
	return cacheClient.Set(key, value)
}
