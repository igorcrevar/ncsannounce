package announce

import (
	"fmt"
	"sync"
)

type MemoryAddressNodeTranslator struct {
	lock            sync.RWMutex
	dict            map[Address]*AddressEntry
	pkAddressMapper PublicKeyAddressMapper
}

func NewMemoryAddressNodeTranslator(pkAddressMapper PublicKeyAddressMapper) *MemoryAddressNodeTranslator {
	r := &MemoryAddressNodeTranslator{
		dict:            make(map[Address]*AddressEntry),
		pkAddressMapper: pkAddressMapper,
	}
	return r
}

func (m *MemoryAddressNodeTranslator) GetEntry(address Address) *AddressEntry {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var (
		entry *AddressEntry
		ok    bool
	)
	if entry, ok = m.dict[address]; !ok {
		return nil
	}
	return entry
}

func (m *MemoryAddressNodeTranslator) AddUpdate(address Address, version uint, node *ENodeUrl) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var (
		entry *AddressEntry
		ok    bool
	)
	if entry, ok = m.dict[address]; ok {
		if entry.Version < version {
			entry.PublicKey = m.pkAddressMapper.GetPublicKeyForAddress(address)
			entry.Node = node
			entry.Version = version
			return true, nil
		} else {
			return true, fmt.Errorf("entry for address %v version is greater than %d", address, version)
		}
	} else {
		m.dict[address] = &AddressEntry{
			PublicKey: m.pkAddressMapper.GetPublicKeyForAddress(address),
			Node:      node,
			Version:   version,
		}
		return false, nil
	}
}

func (m *MemoryAddressNodeTranslator) Remove(address Address) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.dict[address]
	delete(m.dict, address)
	return ok, nil
}

func (m *MemoryAddressNodeTranslator) PruneEntries(addressesToKeep []Address) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make(map[Address]*AddressEntry, len(addressesToKeep))
	for _, address := range addressesToKeep {
		if e, ok := m.dict[address]; ok {
			result[address] = e
		}
	}
	m.dict = result
	return nil
}
