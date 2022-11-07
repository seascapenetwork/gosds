/*Spaghetti transaction without method name and without clear input parameters*/
package spaghetti

import (
	"encoding/json"
	"fmt"
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

func (b *Log) ToString() string {
	interfaces := b.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return ""
	}

	return string(byt)
}

func ParseLog(log map[string]interface{}) Log {
	rawTopics := log["topics"].([]interface{})
	topics := make([]string, len(rawTopics))
	for i, t := range rawTopics {
		topics[i] = t.(string)
	}

	return Log{
		networkId: log["network_id"].(string),
		txId:      log["txid"].(string),
		logIndex:  uint(log["log_index"].(float64)),
		data:      log["data"].(string),
		address:   log["address"].(string),
		topics:    topics,
	}
}

func ParseLogs(rawLogs []interface{}) []Log {
	logs := make([]Log, len(rawLogs))
	for i, rawLog := range rawLogs {
		if rawLog == nil {
			continue
		}
		fmt.Println("Log to parse: ", rawLog)
		l := ParseLog(rawLog.(map[string]interface{}))
		logs[i] = l
	}
	return logs
}
