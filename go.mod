module github.com/filecoin-project/chain-validation

go 1.13

require (
	github.com/dave/jennifer v1.4.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/specs-actors v0.0.0-20200305205312-53bb01da9aeb
	github.com/gopherjs/gopherjs v0.0.0-20190812055157-5d271430af9f // indirect
	github.com/ipfs/go-cid v0.0.4
	github.com/ipfs/go-datastore v0.3.1
	github.com/ipfs/go-ipfs-blockstore v0.1.3
	github.com/ipfs/go-ipld-cbor v0.0.4
	github.com/libp2p/go-libp2p-core v0.3.0
	github.com/multiformats/go-varint v0.0.2
	github.com/stretchr/testify v1.4.0
	github.com/warpfork/go-wish v0.0.0-20190328234359-8b3e70f8e830 // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20200206220010-03c9665e2a66
	golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456 // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
