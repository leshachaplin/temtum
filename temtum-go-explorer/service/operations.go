package service

import (
	"bitbucket.org/dragoninfo/temtum-go-data/data"
)

type Operation interface {
}

type GetLastBlockOperation struct {
	Msg data.Message
}

type SyncOperation struct {
}

type ReadBlocksOperation struct {
	Offset int64
}

type CancelReadBlocks struct {
}
