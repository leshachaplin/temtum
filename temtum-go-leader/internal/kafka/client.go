package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"io"
	"time"
)

type Client struct {
	transactionsReader *kafka.Conn
	blockReader        *kafka.Conn
	blockWriter        *kafka.Conn
	batchBlocksReader  *kafka.Conn
}

// Close kafka connection
func (k *Client) Close() error {
	err := k.transactionsReader.Close()
	err = k.blockReader.Close()
	err = k.batchBlocksReader.Close()
	return err
}

// New kafka client
func New(topicTransaction string, topicBlocks string, url string,
	batchBlockGroup string, blockGroup string, network string) (*Client, error) {

	readTransactions, err := kafka.DialLeader(context.Background(),
		"tcp", url, topicTransaction, 0)
	if err != nil {
		return nil, err
	}

	blockDialer := kafka.Dialer{
		ClientID: blockGroup,
	}

	readBlock, err := blockDialer.DialLeader(context.Background(), network, url, topicBlocks, 0)
	if err != nil {
		return nil, err
	}

	batchBlocksDialer := kafka.Dialer{
		ClientID: batchBlockGroup,
	}

	batchReadBlocks, err := batchBlocksDialer.DialLeader(context.Background(), network, url, topicBlocks, 0)
	if err != nil {
		return nil, err
	}

	blockWriter, err := kafka.DialLeader(context.Background(), network, url, topicBlocks, 0)

	return &Client{
		transactionsReader: readTransactions,
		blockReader:        readBlock,
		blockWriter:        blockWriter,
		batchBlocksReader:  batchReadBlocks,
	}, nil
}

// write message to kafka
func (k *Client) WriteMessage(msg []byte) error {
	//log
	_, err := k.blockWriter.WriteMessages(kafka.Message{
		Value: msg,
	})
	return err
}

//Batch read message from kafka
func (k *Client) ReadMessage(batchMaxSize, maxBytePerMessage, maxBlockAmount, timeout int) ([]byte, int64, error) {
	//log
	t := time.Now().Add(time.Second * time.Duration(timeout))
	err := k.transactionsReader.SetReadDeadline(t)
	if err != nil {
		return nil, 0, err
	}
	batch := k.transactionsReader.ReadBatch(1, batchMaxSize)

	buff := make([]byte, 0, batchMaxSize)
	buff = append(buff, byte('{'))
	buff = append(buff, byte('"'))
	buff = append(buff, []byte("transactions")...)
	buff = append(buff, byte('"'))
	buff = append(buff, byte(':'))
	buff = append(buff, byte('['))
	i := 0

	for {
		b := make([]byte, maxBytePerMessage)
		n, err := batch.Read(b)
		if err != nil && err == io.EOF && !time.Now().After(t) {
			continue
		} else if err != nil && time.Now().After(t) {
			break
		}
		trimB := b[:n]
		buff = append(buff, trimB...)
		buff = append(buff, byte(','))
		i++
		if i >= maxBlockAmount {
			break
		}
	}
	offset := batch.Offset()
	err = batch.Close()

	if !time.Now().After(t) {
		time.Sleep(t.Sub(time.Now()))
	}

	buff[len(buff)-1] = byte(']')
	buff = append(buff, byte('}'))
	return buff, offset, err
}

//Batch read block from kafka
func (k *Client) BatchRead(maxBlocksAmount int, maxBytesPerMessage int,
	BatchSize int, timeout int, do func(msg []byte) error) error {
	t := time.Now().Add(time.Second * time.Duration(timeout))
	err := k.batchBlocksReader.SetReadDeadline(t)
	if err != nil {
		return err
	}
	batch := k.batchBlocksReader.ReadBatch(1, BatchSize)

	buff := make([]byte, 0, BatchSize)
	buff = append(buff, byte('{'))
	buff = append(buff, byte('"'))
	buff = append(buff, []byte("blocks")...)
	buff = append(buff, byte('"'))
	buff = append(buff, byte(':'))
	buff = append(buff, byte('['))

	i := 0

	for {
		b := make([]byte, maxBytesPerMessage)
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		trimB := b[:n]
		buff = append(buff, trimB...)
		buff = append(buff, byte(','))
		i++
		if i >= maxBlocksAmount {
			break
		}
	}


	buff[len(buff)-1] = byte(']')
	buff = append(buff, byte('}'))

	err = do(buff)
	if err == nil || len(buff) < 40 {
		return batch.Close()
	}

	return err
}

func (k *Client) SetOffsetInBatchReader(offset int64, whence int) error {
	_, err := k.batchBlocksReader.Seek(offset, whence)
	return err
}

func (k *Client) GetOffsetInBatchReader() (int64, int) {
	return k.batchBlocksReader.Offset()
}
