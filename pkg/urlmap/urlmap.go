package urlmap

import "errors"

type URLMap interface {
	Add(key, url string)
	GetUrl(key string) (string, error)
	Delete(key string) bool
}

type MemoryURLMap map[string]string

func New() MemoryURLMap {
	return make(MemoryURLMap)
}

func (m MemoryURLMap) Add(key, url string) {
	m[key] = url
}

func (m MemoryURLMap) GetUrl(key string) (string, error) {
	url, ok := m[key]
	if !ok {
		return "", errors.New("not found")
	} 
	return url, nil
}

func (m MemoryURLMap) Delete(key string) bool{
	_, ok := m[key]
	delete(m, key)
	return ok
}