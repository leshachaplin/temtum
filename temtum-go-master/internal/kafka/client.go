package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Client struct {
	transactions *kafka.Conn
	blockReader  *kafka.Conn
}

// Close kafka connection
func (k *Client) Close() error {
	err := k.transactions.Close()
	err = k.blockReader.Close()
	return err
}

// New kafka client
func New(topicTransaction string, topicBlocks string, url string) (*Client, error) {

	senderMsg, err := kafka.DialLeader(context.Background(), "tcp", url, topicBlocks, 0)
	if err != nil {
		return nil, err
	}

	readTransactions, err := kafka.DialLeader(context.Background(), "tcp", url, topicTransaction, 0)
	if err != nil {
		return nil, err
	}

	return &Client{
		blockReader:  senderMsg,
		transactions: readTransactions,
	}, nil
}

// read message from kafka
func (k *Client) ReadWithConn(maxBytes int) (kafka.Message, error) {
	m, err := k.blockReader.ReadMessage(maxBytes)
	return m, err
}

//set offset in kafka
func (k *Client) SetOffsetWithConn(offset int64, whence int) error {
	_, err := k.transactions.Seek(offset, whence)
	return err
}
