package sharding

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nem0z/wiki-pathfinder/storage"
)

type Shard struct {
	*storage.DB
}

type ShardMap struct {
	shards   map[uint64]*Shard
	nbShards uint64
}

func Init(name string, nbShards uint64) (*ShardMap, error) {
	shards := make(map[uint64]*Shard)

	for i := uint64(0); i < nbShards; i++ {
		path := fmt.Sprintf("./shards/%v/shard%v.db", name, i)
		db, err := storage.Init(path)
		if err != nil {
			return nil, err
		}

		shards[i] = &Shard{db}
	}

	return &ShardMap{shards, nbShards}, nil
}

func bytesToUint64(data []byte) uint64 {
	bytes := make([]byte, 8)

	if len(data) >= 8 {
		copy(bytes[:], data[:8])
	} else {
		copy(bytes[8-len(data):], data)
	}

	return binary.BigEndian.Uint64(bytes[:])
}

func (shardMap *ShardMap) GetShard(item interface{}) (*Shard, error) {
	bytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	hash := sha1.Sum(bytes)
	value := bytesToUint64(hash[:])
	key := value % shardMap.nbShards

	shard, ok := shardMap.shards[key]
	if !ok {
		return nil, errors.New("no shard can be associated to this item")
	}

	return shard, nil
}
