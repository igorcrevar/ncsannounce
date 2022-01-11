package announce

type MessageHandler interface {
	Handle([]byte) error
}

type AnnounceMsg struct {
	Code    uint64
	Content []byte // serialized QueryEncodeMsg, NodeCertificateMsg or VersionCertificatesMsg
}

// queryEnodeMsg (0x12)
// This message is sent from a NearlyElectedValidator and regossipped through the network, with the intention of reach out for other NearlyElectedValidator nodes and discover their eNodeURL values, while at the same time sending its own. Note that the sender sends many queries in one message. A query is a pair <destination address, encrypted eNodeURL> so that only the destination address can decrypt the encrypted eNodeURL in it.
// [[dest_address: B_20, encrypted_enode_url: B], version: P, timestamp: P]
// dest_address: the destination validator address this message is querying to.
// encrypted_enode_url: the origin's eNodeURL encrypted with the destination public key.
// version: the current announce version for the origin's eNodeURL.
// timestamp: a message generation timestamp, used to bypass hash caches for messages.
type DestAddressENodeURLPair struct {
	DestAddress   Address
	EncryptedNode []byte
}

type QueryEncodeMsg struct {
	Addresses []DestAddressENodeURLPair
	Version   uint
	Timestamp uint
}

// nodeCertificateMsg (0x17)
// This message holds only ONE eNodeURL and it's meant to be send directly from NearlyElectedValidator to NearlyElectedValidator in a direct connection as a response to a queryEnodeMsg.
// [enode_url: B, version: P]
type NodeCertificateMsg struct {
	Node    ENodeUrl
	Version uint
}

// versionCertificatesMsg (0x16)
// This messages holds MANY version certificates. It is used mostly in two ways:
// To share the WHOLE version table a FullNode has.
// To share updated parts of the table of a FullNode.
// [[version: P, signature: B]]
// version: the current highest known announce version for a validator (deduced from the signature), according to the peer that gossipped this message.
// signature: the signature for the version payload string ('versionCertificate'|version) by the emmitter of the certificate. Note that the address and
// public key can be deduced from this signature.
// TODO:!
type VersionCertificate struct {
	Version   uint
	Node      ENodeUrl
	PublicKey *PublicKey
	Signature []byte // TODO:!
}

type VersionCertificateMsg struct {
	Certificates []VersionCertificate
}
