module github.com/julwil/bazo-client

go 1.14

require (
	github.com/boltdb/bolt v1.3.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/julwil/bazo-miner v0.0.0-20200303120255-9fe62280f40b
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/urfave/cli v1.22.3
	github.com/willf/bitset v1.1.10 // indirect
	github.com/willf/bloom v2.0.3+incompatible // indirect
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073 // indirect
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
)

replace github.com/julwil/bazo-miner => ../bazo-miner // Packages from bazo-miner are resolved locally, rather than with the specified version.
