package store

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"producer/consumer/consumer_structs"
	"producer/structs"
	"strconv"
)

type StorageService struct {
	ConsumerConfig consumer_structs.ConsumerConfig
	CommonConfig   structs.CommonConfig
}

func (s *StorageService) SaveConsumedMessage(message structs.Message) error {
	opt := badger.DefaultOptions(s.ConsumerConfig.BadgerTempDir)
	db, err := badger.Open(opt)
	if err != nil {
		log.Printf("Error in opening badger DB connection. Err: [%v]", err)
		return err
	}
	defer func(db *badger.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error in closing badger DB connection. Err: [%v]", err)
		}
	}(db)
	key := []byte(message.Id)
	value := []byte(fmt.Sprintf("%f", message.Value))

	txn := db.NewTransaction(true)
	defer txn.Discard()

	// Get the value for key first to check value already exists or not
	entry, er := txn.Get(key)
	if er != nil && er != badger.ErrKeyNotFound {
		log.Printf("Error in getting value from badgerDB for key [%v]. Error: [%v]", message.Id, er)
		return er
	}
	if er == nil {
		// previous entry found, add the value to the new value
		prevValue, _ := entry.ValueCopy(nil)
		log.Printf("Found value: [%v]", string(prevValue))
		prevValueFloat, gErr := strconv.ParseFloat(string(prevValue), 64)
		if gErr != nil {
			log.Printf("Error in converting []byte to float. Error: [%v]", gErr)
			return gErr
		}
		value = []byte(fmt.Sprintf("%f", message.Value+prevValueFloat))
	}

	// Set the final value
	if err := txn.Set(key, value); err != nil {
		log.Printf("Error in setting KV in badger db. Error: [%v]", err)
		return err
	}
	if err := txn.Commit(); err != nil {
		log.Printf("Error in committing KV in badger db. Error: [%v]", err)
		return err
	}

	return nil
}

func (s *StorageService) GetValue(id string) (structs.Message, error) {
	var message structs.Message
	opt := badger.DefaultOptions(s.ConsumerConfig.BadgerTempDir)
	db, err := badger.Open(opt)
	if err != nil {
		log.Printf("Error in opening badger DB connection. Err: [%v]", err)
		return structs.Message{}, err
	}
	defer func(db *badger.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error in closing badger DB connection. Err: [%v]", err)
		}
	}(db)
	txn := db.NewTransaction(false)
	defer txn.Discard()

	key := []byte(id)
	entry, err := txn.Get(key)
	if err != nil {
		log.Printf("Error in getting value from badgerDB for key [%v]. Error: [%v]", id, err)
		return structs.Message{}, err
	}
	value, _ := entry.ValueCopy(nil)
	log.Printf("Found value: [%v]", string(value))
	valueFloat, gErr := strconv.ParseFloat(string(value), 64)
	if gErr != nil {
		log.Printf("Error in converting []byte to float. Error: [%v]", gErr)
		return structs.Message{}, gErr
	}

	message.Id = id
	message.Value = valueFloat

	return message, nil
}
