package gsmap

import "sync"

type Gsmap struct {
	data      map[string]string
	lock_data sync.Mutex
}

func NewGSMap() *Gsmap {
	return &Gsmap{data: make(map[string]string)}
}

func (m *Gsmap) Add(key, value string) {
	m.lock_data.Lock()
	defer m.lock_data.Unlock()
	_, ok := m.data[key]
	if ok {
		panic("The key of map already exists")
	}
	m.data[key] = value
}
