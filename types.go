package announce

import (
	"crypto/ecdsa"
	"encoding/hex"
	"time"
)

const (
	AddressLength = 20 // HashLength is the expected length of the hash
	HashLength    = 32 // AddressLength is the expected length of the address
)

// Address of etherium account
type Address [AddressLength]byte

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

type PublicKey ecdsa.PublicKey

type ENodeUrl string // libp2p url?

// type AddressPubKeyPair struct {
// 	address Address
// 	parsed  bool
// 	PubKey  *PublicKey
// }

// AddressEntry is an entry for the valEnodeTable.
type AddressEntry struct {
	// Address                      Address
	PublicKey                    *PublicKey
	Node                         *ENodeUrl
	Version                      uint
	HighestKnownVersion          uint
	NumQueryAttemptsForHKVersion uint
	LastQueryTimestamp           *time.Time
}

func (a *Address) Hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

// String implements fmt.Stringer.
func (a *Address) String() string {
	return string(a.Hex())
}
