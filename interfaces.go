package announce

type AnnounceNetwork interface {
	// Gossip gossips protocol messages // ethMsgCode uint64
	Gossip(message []byte) error
	// Multicast will send the eth message (with the message's payload and msgCode field set to the params
	// payload and ethMsgCode respectively) to the nodes with the signing address in the destAddresses param.
	Multicast(destAddresses []Address, message []byte, sendToSelf bool) error
	// RetrieveValidatorConnSet returns the validator connection set
	RetrieveValidatorConnSet() (map[Address]bool, error)
}

// AddressProvider provides the different addresses the announce manager needs
type AddressProvider interface {
	SelfNode() *ENodeUrl
	ValidatorAddress() Address
	IsValidating() bool
}

type NodeHelper interface {
	DecryptNode(nodeCrypted []byte) (*ENodeUrl, error)
	EncryptNode(node *ENodeUrl) ([]byte, error)
}

type PublicKeyAddressMapper interface {
	PublicKey2Address(publicKey *PublicKey) Address
	GetPublicKeyForAddress(address Address) *PublicKey
}

type PublicKeyNodeTranslator interface {
	// GetNode(address Address) *ENodeUrl
	GetEntry(address Address) *AddressEntry
	AddUpdate(address Address, version uint, node *ENodeUrl) (bool, error)
	Remove(address Address) (bool, error)
	PruneEntries(addressesToKeep []Address) error
}
