package cstorage

import (
	"github.com/boltdb/bolt"
	"github.com/way365/bazo-miner/protocol"
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

func ReadTransaction(txHash [32]byte) (transaction protocol.Transaction) {
	var encodedTx []byte

	var accTx *protocol.AccTx
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ACCOUNT_TX_BUCKET))
		encodedTx = b.Get(txHash[:])
		return nil
	})
	if encodedTx != nil {
		return accTx.Decode(encodedTx)
	}

	var fundsTx *protocol.FundsTx
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FUND_TX_BUCKET))
		encodedTx = b.Get(txHash[:])
		return nil
	})
	if encodedTx != nil {
		return fundsTx.Decode(encodedTx)
	}

	var configTx *protocol.ConfigTx
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CONFIG_TX_BUCKET))
		encodedTx = b.Get(txHash[:])
		return nil
	})
	if encodedTx != nil {
		return configTx.Decode(encodedTx)
	}

	var stakingTx *protocol.StakeTx
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(STAKING_TX_BUCKET))
		encodedTx = b.Get(txHash[:])
		return nil
	})
	if encodedTx != nil {
		return stakingTx.Decode(encodedTx)
	}

	var updateTx *protocol.UpdateTx
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UPDATE_TX_BUCKET))
		encodedTx = b.Get(txHash[:])
		return nil
	})
	if encodedTx != nil {
		return updateTx.Decode(encodedTx)
	}

	return nil
}
