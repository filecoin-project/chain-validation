module github.com/filecoin-project/chain-validation

go 1.13

require (
	github.com/dave/jennifer v1.4.0
	github.com/filecoin-project/filecoin-ffi v0.30.4-0.20200716204036-cddc56607e1d
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-crypto v0.0.0-20191218222705-effae4ea9f03
	github.com/filecoin-project/go-fil-commcid v0.0.0-20200716160307-8f644712406f
	github.com/filecoin-project/specs-actors v0.8.1-0.20200720061236-f4719fdd7d90
	github.com/gorilla/rpc v1.2.0
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-cid v0.0.6
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-hamt-ipld v0.1.1-0.20200501020327-d53d20a7063e // indirect
	github.com/ipfs/go-ipfs-blockstore v1.0.0
	github.com/ipfs/go-ipld-cbor v0.0.5-0.20200428170625-a0bd04d3cbdf
	github.com/ipfs/go-ipld-format v0.2.0 // indirect
	github.com/ipfs/go-log v1.0.4
	github.com/ipsn/go-secp256k1 v0.0.0-20180726113642-9d62b9f0bc52
	github.com/kr/pretty v0.2.0 // indirect
	github.com/libp2p/go-libp2p-core v0.5.3
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/multiformats/go-multihash v0.0.14
	github.com/multiformats/go-varint v0.0.5
	github.com/stretchr/testify v1.5.1
	github.com/warpfork/go-wish v0.0.0-20200122115046-b9ea61034e4a // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20200710004633-5379fc63235d
	golang.org/x/crypto v0.0.0-20200427165652-729f1e841bcc // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/sys v0.0.0-20200427175716-29b57079015a // indirect
	golang.org/x/tools v0.0.0-20200318150045-ba25ddc85566 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	gopkg.in/yaml.v2 v2.2.8 // indirect
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
