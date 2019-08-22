package rocksdb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
)

func openLevelDB(filepath string) (*leveldb.DB, error) {
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
		Strict: opt.StrictAll,
	}
	db, err := leveldb.OpenFile(filepath, o)
	if err == nil {
		return db, nil
	}
	if _, ok := err.(*errors.ErrCorrupted); ok {
		log.Printf("recovering leveldb: %v", err)
		db, err = leveldb.RecoverFile(filepath, o)
		if err != nil {
			log.Printf("failed to recover leveldb: %v", err)
			return nil, err
		}
		return db, nil
	}
	log.Printf("failed to open leveldb: %v", err)
	return nil, err
}