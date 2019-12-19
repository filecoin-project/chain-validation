module github.com/filecoin-project/chain-validation

go 1.13

require (
	github.com/filecoin-project/go-address v0.0.0-20191219011437-af739c490b4f
	github.com/filecoin-project/go-amt-ipld v0.0.0-20191205011053-79efc22d6cdc
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190812055157-5d271430af9f // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/ipfs/go-cid v0.0.4-0.20191112011718-79e75dffeb10
	github.com/ipfs/go-datastore v0.1.0
	github.com/ipfs/go-hamt-ipld v0.0.12-0.20190910032255-ee6e898f0456
	github.com/ipfs/go-ipfs-blockstore v0.1.0
	github.com/ipfs/go-ipld-cbor v0.0.3
	github.com/ipfs/go-ipld-format v0.0.2 // indirect
	github.com/ipfs/go-log v1.0.0 // indirect
	github.com/libp2p/go-libp2p-core v0.2.4
	github.com/multiformats/go-multihash v0.0.9
	github.com/polydawn/refmt v0.0.0-20190809202753-05966cbd336a
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/warpfork/go-wish v0.0.0-20190328234359-8b3e70f8e830 // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20191216205031-b047b6acb3c0
	golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
