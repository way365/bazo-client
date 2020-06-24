package cstorage

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/julwil/bazo-client/util"
	"log"
	"time"
)

var (
	db     *bolt.DB
	logger *log.Logger
)

const (
	ERROR_MSG                = "Initiate storage aborted: "
	LAST_BLOCK_HEADER_BUCKET = "lastblockheader"
	BLOCK_HEADER_BUCKET      = "blockheaders"
	ACCOUNT_TX_BUCKET        = "account_transactions"
	FUND_TX_BUCKET           = "fund_transactions"
	CONFIG_TX_BUCKET         = "config_transactions"
	STAKING_TX_BUCKET        = "staking_transactions"
	UPDATE_TX_BUCKET         = "update_transactions"
)

//Entry function for the storage package
func Init(dbname string) {
	logger = util.InitLogger()

	var err error
	db, err = bolt.Open(dbname, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		logger.Fatal(ERROR_MSG, err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte(BLOCK_HEADER_BUCKET))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}

		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte(LAST_BLOCK_HEADER_BUCKET))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}

		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(ACCOUNT_TX_BUCKET))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(FUND_TX_BUCKET))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(CONFIG_TX_BUCKET))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(STAKING_TX_BUCKET))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(UPDATE_TX_BUCKET))
		if err != nil {
			return fmt.Errorf(ERROR_MSG+"Create bucket: %s", err)
		}
		return nil
	})
}

func TearDown() {
	db.Close()
}
