package conf

import (
	"sync"
)

var (
	defaultConfiger *Configer     = nil
	singletonLock   *sync.RWMutex = new(sync.RWMutex)
)

func installConfiger(cfg *Configer) {
	singletonLock.Lock()
	defer singletonLock.Unlock()
	defaultConfiger = cfg
}

func GetValueByString(section, key string) (string, error) {
	if nil == defaultConfiger {
		return "", defaultError
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByString(section, key)
}

func GetValueByStringCarryDefault(section, key, defaultValue string) string {
	if nil == defaultConfiger {
		return defaultValue
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByStringCarryDefault(section, key, defaultValue)
}

func GetValueByBool(section, key string) (bool, error) {
	if nil == defaultConfiger {
		return false, defaultError
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByBool(section, key)
}

func GetValueByBoolCarryDefault(section, key string, defaultValue bool) bool {
	if nil == defaultConfiger {
		return defaultValue
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByBoolCarryDefault(section, key, defaultValue)
}

func GetValueByFloat64(section, key string) (float64, error) {
	if nil == defaultConfiger {
		return 0, defaultError
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByFloat64(section, key)
}

func GetValueByInt(section, key string) (int, error) {
	if nil == defaultConfiger {
		return 0, defaultError
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByInt(section, key)
}

func GetValueByIntCarryDefault(section, key string, defaultValue int) int {
	if nil == defaultConfiger {
		return defaultValue
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByIntCarryDefault(section, key, defaultValue)
}

func GetValueByInt64(section, key string) (int64, error) {
	if nil == defaultConfiger {
		return 0, defaultError
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByInt64(section, key)
}

func GetValueByStringList(section, key string) ([]string, error) {
	if nil == defaultConfiger {
		return nil, defaultError
	}

	singletonLock.RLock()
	defer singletonLock.RUnlock()
	return defaultConfiger.GetValueByStringList(section, key)
}
