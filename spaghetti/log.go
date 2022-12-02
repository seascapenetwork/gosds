/*Spaghetti transaction without method name and without clear input parameters*/
package spaghetti

import (
	"encoding/json"
	"errors"

	"github.com/blocklords/gosds/message"
)

type Log struct {
	networkId string
	txId      string // txId column
	logIndex  uint
	data      string // text data type
	topics    []string
	address   string
}

func (b *Log) NetworkID() string {
	return b.networkId
}

func (b *Log) LogIndex() uint {
	return b.logIndex
}

func (b *Log) Data() string {
	return b.data
}

func (b *Log) Topics() []string {
	return b.topics
}

func (b *Log) TxId() string {
	return b.txId
}

// JSON representation of the spaghetti.Log
func (b *Log) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"network_id": b.networkId,
		"txid":       b.txId,
		"log_index":  b.logIndex,
		"data":       b.data,
		"topics":     b.topics,
		"address":    b.address,
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

	return &Log{
		networkId: network_id,
		address:   address,
		txId:      txid,
		logIndex:  uint(log_index),
		data:      data,
		topics:    topics,
	}, nil
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
