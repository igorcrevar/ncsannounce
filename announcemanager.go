package announce

import (
	"log"
	"os"
	"time"
)

type AnnounceManagerConfig struct {
	logger *log.Logger
}

func (m *AnnounceManagerConfig) WithLogger(log *log.Logger) {
	m.logger = log
}

func DefaultAnnounceManagerConfig() AnnounceManagerConfig {
	return AnnounceManagerConfig{
		logger: log.New(os.Stdout, "announce_manager", log.LstdFlags),
	}
}

type AnnounceManager struct {
	config       AnnounceManagerConfig
	addrProvider AddressProvider
	network      AnnounceNetwork

	gossipCache GossipCache

	// The enode certificate message map contains the most recently generated
	// enode certificates for each external node ID (e.g. will have one entry per proxy
	// for a proxied validator, or just one entry if it's a standalone validator).
	// Each proxy will just have one entry for their own external node ID.
	// Used for proving itself as a validator in the handshake for externally exposed nodes,
	// or by saving latest generated certificate messages by proxied validators to send
	// to their proxies.
	// enodeCertificateMsgMap     map[enode.ID]*istanbul.EnodeCertMsg
	// enodeCertificateMsgVersion uint
	// enodeCertificateMsgMapMu   sync.RWMutex // This protects both enodeCertificateMsgMap and enodeCertificateMsgVersion

	lastVersionCertificatesGossiped map[Address]time.Time
	// lastVersionCertificatesGossipedMu sync.RWMutex

	lastQueryEnodeGossiped map[Address]time.Time
	// lastQueryEnodeGossipedMu sync.RWMutex
}

// NewAnnounceManager creates a new AnnounceManager using the valEnodeTable given. It is
// the responsibility of the caller to close the valEnodeTable, the AnnounceManager will
// not do it.
func NewAnnounceManager(network AnnounceNetwork,
	addrProvider AddressProvider,
	gossipCache GossipCache,
	config AnnounceManagerConfig) *AnnounceManager {
	am := &AnnounceManager{
		config:                          config,
		network:                         network,
		addrProvider:                    addrProvider,
		gossipCache:                     gossipCache,
		lastQueryEnodeGossiped:          make(map[Address]time.Time),
		lastVersionCertificatesGossiped: make(map[Address]time.Time),
	}
	// versionCertificateTable, err := enodes.OpenVersionCertificateDB(config.VcDbPath)
	// if err != nil {
	// 	am.logger.Crit("Can't open VersionCertificateDB", "err", err, "dbpath", config.VcDbPath)
	// }
	// am.versionCertificateTable = versionCertificateTable
	return am
}

func (m *AnnounceManager) Close() error {
	// No need to close valEnodeTable since it's a reference,
	// the creator of this announce manager is the responsible for
	// closing it.
	// return m.versionCertificateTable.Close()
	return nil
}
