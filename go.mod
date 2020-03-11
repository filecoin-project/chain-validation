module github.com/filecoin-project/chain-validation

go 1.13

require (
	github.com/dave/jennifer v1.4.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-fil-commcid v0.0.0-20200208005934-2b8bd03caca5
	github.com/filecoin-project/go-sectorbuilder v0.0.2-0.20200309211213-75e9124a1904
	github.com/filecoin-project/specs-actors v0.0.0-20200306000749-99e98e61e2a0
	github.com/ipfs/go-cid v0.0.5
	github.com/ipfs/go-datastore v0.4.1
	github.com/ipfs/go-ipfs-blockstore v0.1.3
	github.com/ipfs/go-ipld-cbor v0.0.4
	github.com/ipfs/go-log v1.0.2 // indirect
	github.com/libp2p/go-libp2p-core v0.3.0
	github.com/multiformats/go-multihash v0.0.13
	github.com/multiformats/go-varint v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20200206220010-03c9665e2a66
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
