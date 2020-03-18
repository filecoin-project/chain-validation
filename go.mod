module github.com/filecoin-project/chain-validation

go 1.13

require (
	github.com/dave/jennifer v1.4.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-fil-commcid v0.0.0-20200208005934-2b8bd03caca5
	github.com/filecoin-project/specs-actors v0.0.0-20200306000749-99e98e61e2a0
	github.com/gopherjs/gopherjs v0.0.0-20190812055157-5d271430af9f // indirect
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-cid v0.0.5
	github.com/ipfs/go-datastore v0.4.1
	github.com/ipfs/go-ipfs-blockstore v0.1.3
	github.com/ipfs/go-ipld-cbor v0.0.4
	github.com/ipfs/go-ipld-format v0.0.2 // indirect
	github.com/ipfs/go-log v1.0.2 // indirect
	github.com/libp2p/go-libp2p-core v0.3.0
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/multiformats/go-multihash v0.0.13
	github.com/multiformats/go-varint v0.0.5
	github.com/polydawn/refmt v0.0.0-20190809202753-05966cbd336a // indirect
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/warpfork/go-wish v0.0.0-20190328234359-8b3e70f8e830 // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20200206220010-03c9665e2a66
	go.uber.org/multierr v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
