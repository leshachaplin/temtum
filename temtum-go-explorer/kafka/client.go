package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	isEmptyBuff = 27
)

type Client struct {
	reader *kafka.Conn
}

// Close kafka connection
func (k *Client) Close() error {
	err := k.reader.Close()
	return err
}

// New kafka client
func New(topic string, url string, batchGroup string) (*Client, error) {

	dialer := kafka.Dialer{
		ClientID: batchGroup,
	}

	readerMsg, err := dialer.DialLeader(context.Background(), "tcp", url, topic, 0)
	if err != nil {
		return nil, err
	}

	return &Client{
		reader: readerMsg,
	}, nil
}

//Batch read message from kafka
func (k *Client) BatchRead(maxBlocksAmount int, timeout int, batchSize int,
	maxBytesPerMessage int, do func(msg []byte) error)  error {
	//log
	err := k.reader.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
	if err != nil {
		return err
	}
	batch := k.reader.ReadBatch(1, batchSize)

	buff := make([]byte, 0, batchSize)
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
	if err == nil || len(buff) < isEmptyBuff {
		err	= batch.Close()
	}

	return err
}

func (k *Client) SetOffset(offset int64, whence int) error {
	_, err := k.reader.Seek(offset, whence)
	return err
}

func (k *Client) GetOffset() (int64, int) {
	return k.reader.Offset()
}
