/*Spaghetti transaction without method name and without clear input parameters*/
package spaghetti

import (
	"encoding/json"
	"errors"

	"github.com/blocklords/gosds/message"
)

type Log struct {
	NetworkId      string
	Txid           string // txId column
	BlockNumber    uint64
	BlockTimestamp uint64
	LogIndex       uint
	Data           string // text data type
	Topics         []string
	Address        string
}

// JSON representation of the spaghetti.Log
func (b *Log) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"network_id":      b.NetworkId,
		"txid":            b.Txid,
		"block_timestamp": b.BlockTimestamp,
		"block_number":    b.BlockNumber,
		"log_index":       b.LogIndex,
		"data":            b.Data,
		"topics":          b.Topics,
		"address":         b.Address,
	}
}

// JSON string representation of the spaghetti.Log
func (b *Log) ToString() string {
	interfaces := b.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return ""
	}

	return string(byt)
}

// Convert the JSON into spaghetti.Log
func ParseLog(parameters map[string]interface{}) (*Log, error) {
	topics, err := message.GetStringList(parameters, "topics")
	if err != nil {
		return nil, err
	}
	network_id, err := message.GetString(parameters, "network_id")
	if err != nil {
		return nil, err
	}
	txid, err := message.GetString(parameters, "txid")
	if err != nil {
		return nil, err
	}
	log_index, err := message.GetUint64(parameters, "log_index")
	if err != nil {
		return nil, err
	}
	data, err := message.GetString(parameters, "data")
	if err != nil {
		return nil, err
	}
	address, err := message.GetString(parameters, "address")
	if err != nil {
		return nil, err
	}

	block_timestamp, err := message.GetUint64(parameters, "block_timestamp")
	if err != nil {
		return nil, err
	}
	block_number, err := message.GetUint64(parameters, "block_number")
	if err != nil {
		return nil, err
	}

	return &Log{
		NetworkId:      network_id,
		Address:        address,
		Txid:           txid,
		BlockNumber:    block_number,
		BlockTimestamp: block_timestamp,
		LogIndex:       uint(log_index),
		Data:           data,
		Topics:         topics,
	}, nil
}

// Serielizes the Log.Topics into the byte array
func (b *Log) TopicRaw() []byte {
	byt, err := json.Marshal(b.Topics)
	if err != nil {
		return []byte{}
	}

	return byt
}

// Converts the byte series into the topic list
func (b *Log) ParseTopics(raw []byte) error {
	var topics []string
	err := json.Unmarshal(raw, &topics)
	if err != nil {
		return err
	}
	b.Topics = topics

	return nil
}

// Parse list of Logs into array of spaghetti.Log
func ParseLogs(raw_logs []interface{}) ([]*Log, error) {
	logs := make([]*Log, len(raw_logs))
	for i, raw := range raw_logs {
		if raw == nil {
			continue
		}
		log_map, ok := raw.(map[string]interface{})
		if !ok {
			return nil, errors.New("the log is not a map")
		}
		l, err := ParseLog(log_map)
		if err != nil {
			return nil, err
		}
		logs[i] = l
	}
	return logs, nil
}
