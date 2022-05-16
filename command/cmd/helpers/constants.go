package helpers

const (

	// DefaultNamespace is the default namespace for all operations
	DefaultNamespace = "default"

	DefaultStorageclass = ""

	DefaultCAImage   = "hyperledger/fabric-ca"
	DefaultCAVersion = "1.4.9"

	DefaultPeerImage   = "hyperledger/fabric-peer"
	DefaultPeerVersion = "2.4.1"
	//DefaultPeerVersion = "1.4.2"

	DefaultFSServerImage   = "quay.io/kfsoftware/fs-peer"
	DefaultFSServerVersion = "amd64-2.2.0-0.0.1"

	DefaultCouchDBImage   = "couchdb"
	DefaultCouchDBVersion = "3.1.1"

	DefaultOrdererImage   = "hyperledger/fabric-orderer"
	DefaultOrdererVersion = "2.4.1"
	//DefaultOrdererVersion = "1.4.2"
)
