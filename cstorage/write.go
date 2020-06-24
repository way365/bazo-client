package cstorage

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/julwil/bazo-miner/protocol"
)

func WriteBlockHeader(header *protocol.Block) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockheaders"))
		err := b.Put(header.Hash[:], header.EncodeHeader())

		return err
	})

	return err
}

//Before saving the last block header, delete all existing entries.
func WriteLastBlockHeader(header *protocol.Block) (err error) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lastblockheader"))
		b.ForEach(func(k, v []byte) error {
			b.Delete(k)

			return nil
		})

		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("lastblockheader"))
		err := b.Put(header.Hash[:], header.EncodeHeader())

		return err
	})

	return err
}

func WriteTransaction(txHash [32]byte, tx protocol.Transaction) (err error) {
	var bucket string
	switch tx.(type) {
	case *protocol.AccTx:
		bucket = ACCOUNT_TX_BUCKET
	case *protocol.FundsTx:
		bucket = FUND_TX_BUCKET
	case *protocol.ConfigTx:
		bucket = CONFIG_TX_BUCKET
	case *protocol.StakeTx:
		bucket = STAKING_TX_BUCKET
	case *protocol.UpdateTx:
		bucket = UPDATE_TX_BUCKET
	default:
		return errors.New("invalid tx type")
	}

	err = db.Update(func(boltTx *bolt.Tx) error {
		b := boltTx.Bucket([]byte(bucket))
		err := b.Put(txHash[:], tx.Encode())

		return err
	})

	return err
}
