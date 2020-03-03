package cstorage

import (
	"github.com/boltdb/bolt"
	"github.com/julwil/bazo-miner/protocol"
)

func ReadBlockHeader(hash [32]byte) (header *protocol.Block) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockheaders"))
		encodedHeader := b.Get(hash[:])
		header = header.Decode(encodedHeader)

		return nil
	})

	if header == nil {
		return nil
	}

	return header
}

func ReadLastBlockHeader() (header *protocol.Block) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lastblockheader"))
		cb := b.Cursor()
		_, encodedHeader := cb.First()
		header = header.Decode(encodedHeader)

		return nil
	})

	if header == nil {
		return nil
	}

	return header
}
