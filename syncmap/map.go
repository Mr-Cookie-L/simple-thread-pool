package syncmap

type SyncMap struct {
	m            map[interface{}]interface{}
	readSig      chan *readSig
	writeSig     chan *writeSig
	delSig       chan *delSig
	terminateSig chan struct{}
}

type readSig struct {
	key   interface{}
	value interface{}
	ok    chan bool
}

type writeSig struct {
	key   interface{}
	value interface{}
	ok    chan bool
}

type delSig struct {
	key interface{}
	ok  chan bool
}

func New(kind interface{}) *SyncMap {
	m := &SyncMap{
		m:            make(map[interface{}]interface{}),
		readSig:      make(chan *readSig),
		writeSig:     make(chan *writeSig),
		delSig:       make(chan *delSig),
		terminateSig: make(chan struct{}),
	}
	m.run()
	return m
}

func (m *SyncMap) run() {
	for {
		select {
		case readSig := <-m.readSig:
			if value, ok := m.m[readSig.key]; ok {
				readSig.value = value
				readSig.ok <- true
			} else {
				readSig.ok <- false
			}
		case writeSig := <-m.writeSig:
			m.m[writeSig.key] = writeSig.value
			writeSig.ok <- true
		case delSig := <-m.delSig:
			delete(m.m, delSig.key)
			delSig.ok <- true
		case <-m.terminateSig:
			break
		}
	}
}

func (m *SyncMap) Get(key interface{}) (interface{}, bool) {
	readSig := &readSig{
		key: key,
	}
	m.readSig <- readSig
	ok := <-readSig.ok
	return readSig.value, ok
}

func (m *SyncMap) Set(key, value interface{}) bool {
	writeSig := &writeSig{
		key:   key,
		value: value,
	}
	m.writeSig <- writeSig
	ok := <-writeSig.ok
	return ok
}

func (m *SyncMap) Delete(key interface{}) bool {
	delSig := &delSig{
		key: key,
	}
	m.delSig <- delSig
	ok := <-delSig.ok
	return ok
}
