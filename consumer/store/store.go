package store

import (
	"consumer/consumer_structs"
	"consumer/helper"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"strconv"
)

var (
	storageService StorageService
)

type StorageService struct {
	ConsumerConfig consumer_structs.ConsumerConfig
	Db             *badger.DB
}

func GetService() *StorageService {
	return &storageService
}

func InitiateStorageService() error {
	consumerConfig := helper.LoadConsumerConfiguration(os.Getenv("APP_HOME") + "/config/config.json")
	db, err := InitiateBadgerDB(consumerConfig)
	if err != nil {
		log.Printf("Consumer.InitiateStorageService. Error: [%v]", err)
		return err
	}

	storageService = StorageService{
		ConsumerConfig: consumerConfig,
		Db:             db,
	}

	return nil
}

func InitiateBadgerDB(consumerConfig consumer_structs.ConsumerConfig) (*badger.DB, error) {
	opt := badger.DefaultOptions(consumerConfig.BadgerTempDir)
	db, err := badger.Open(opt)
	if err != nil {
		log.Printf("consumer.store.InitiateBadgerDB: Error in opening badger DB connection. Err: [%v]", err)
		return nil, err
	}
	return db, nil
}

func (s *StorageService) SaveConsumedMessage(message consumer_structs.Message) error {
	key := []byte(message.Id)
	value := []byte(fmt.Sprintf("%.2f", message.Value))

	txn := s.Db.NewTransaction(true)
	defer txn.Discard()

	// Get the value for key first to check value already exists or not
	entry, er := txn.Get(key)
	if er != nil && er != badger.ErrKeyNotFound {
		log.Printf("consumer.store.SaveConsumedMessage:Error in getting value from badgerDB for key [%v]. Error: [%v]", message.Id, er)
		return er
	}
	if er == nil {
		// previous entry found, add the value to the new value
		prevValue, _ := entry.ValueCopy(nil)
		log.Printf("Found value: [%v]", string(prevValue))
		prevValueFloat, gErr := strconv.ParseFloat(string(prevValue), 64)
		if gErr != nil {
			log.Printf("consumer.store.SaveConsumedMessage: Error in converting []byte to float. Error: [%v]", gErr)
			return gErr
		}
		value = []byte(fmt.Sprintf("%.2f", message.Value+prevValueFloat))
	}

	// Set the final value
	if err := txn.Set(key, value); err != nil {
		log.Printf("consumer.store.SaveConsumedMessage: Error in setting KV in badger db. Error: [%v]", err)
		return err
	}
	if err := txn.Commit(); err != nil {
		log.Printf("consumer.store.SaveConsumedMessage: Error in committing KV in badger db. Error: [%v]", err)
		return err
	}

	return nil
}

func (s *StorageService) GetValue(id string) (consumer_structs.Message, error) {
	var message consumer_structs.Message
	txn := s.Db.NewTransaction(false)
	defer txn.Discard()

	key := []byte(id)
	entry, err := txn.Get(key)
	if err != nil {
		log.Printf("consumer.store.GetValue: Error in getting value from badgerDB for key [%v]. Error: [%v]", id, err)
		return consumer_structs.Message{}, err
	}
	if err := txn.Commit(); err != nil {
		log.Printf("consumer.store.GetValue: Error in committing KV read in badger db. Error: [%v]", err)
		return consumer_structs.Message{}, err
	}

	value, _ := entry.ValueCopy(nil)
	log.Printf("Found value: [%v]", string(value))
	valueFloat, gErr := strconv.ParseFloat(string(value), 64)
	if gErr != nil {
		log.Printf("consumer.store.GetValue: Error in converting []byte to float. Error: [%v]", gErr)
		return consumer_structs.Message{}, gErr
	}

	message.Id = id
	message.Value = valueFloat

	return message, nil
}
