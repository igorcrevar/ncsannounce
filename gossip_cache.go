package announce

import (
	"log"
	"os"

	lru "github.com/hashicorp/golang-lru"
)

type GossipCache interface {
	MarkMessageProcessedByPeer(peerNodeAddr Address, payloadHash Hash)
	CheckIfMessageProcessedByPeer(peerNodeAddr Address, payloadHash Hash) bool

	MarkMessageProcessedBySelf(payloadHash Hash)
	CheckIfMessageProcessedBySelf(payloadHash Hash) bool
}

type LRUGossipCacheConfig struct {
	InmemoryPeers    int
	InmemoryMessages int
	Logger           *log.Logger
}

func (m *LRUGossipCacheConfig) WithLogger(log *log.Logger) {
	m.Logger = log
}

func GetDefaultLRUGossipCacheConfig() LRUGossipCacheConfig {
	return LRUGossipCacheConfig{
		InmemoryPeers:    40,
		InmemoryMessages: 1024,
		Logger:           log.New(os.Stdout, "gossip_cache", log.LstdFlags),
	}
}

type LRUGossipCache struct {
	peerRecentMessages *lru.ARCCache        // the cache of peer's recent messages
	selfRecentMessages *lru.ARCCache        // the cache of self recent messages
	config             LRUGossipCacheConfig // config
}

func NewLRUGossipCache(config LRUGossipCacheConfig) *LRUGossipCache {
	peerRecentMessages, err := lru.NewARC(config.InmemoryPeers)
	if err != nil {
		config.Logger.Fatal("Failed to create recent messages cache", "err", err)
	}
	selfRecentMessages, err := lru.NewARC(config.InmemoryMessages)
	if err != nil {
		config.Logger.Fatal("Failed to create known messages cache", "err", err)
	}
	return &LRUGossipCache{
		peerRecentMessages: peerRecentMessages,
		selfRecentMessages: selfRecentMessages,
		config:             config,
	}
}

func (gc *LRUGossipCache) MarkMessageProcessedByPeer(peerNodeAddr Address, payloadHash Hash) {
	ms, ok := gc.peerRecentMessages.Get(peerNodeAddr)
	var m *lru.ARCCache
	if ok {
		m, _ = ms.(*lru.ARCCache)
	} else {
		m, _ = lru.NewARC(gc.config.InmemoryMessages)
		gc.peerRecentMessages.Add(peerNodeAddr, m)
	}
	m.Add(payloadHash, true)
}

func (gc *LRUGossipCache) CheckIfMessageProcessedByPeer(peerNodeAddr Address, payloadHash Hash) bool {
	ms, ok := gc.peerRecentMessages.Get(peerNodeAddr)
	var m *lru.ARCCache
	if ok {
		m, _ = ms.(*lru.ARCCache)
		_, ok := m.Get(payloadHash)
		return ok
	}

	return false
}

func (gc *LRUGossipCache) MarkMessageProcessedBySelf(payloadHash Hash) {
	gc.selfRecentMessages.Add(payloadHash, true)
}

func (gc *LRUGossipCache) CheckIfMessageProcessedBySelf(payloadHash Hash) bool {
	_, ok := gc.selfRecentMessages.Get(payloadHash)
	return ok
}
