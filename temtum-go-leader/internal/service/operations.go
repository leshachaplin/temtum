package service

type Operation interface {
}

type ReadRoleOperation struct {
	Role string
}

type ReadTransactionOperation struct {
}

type CancelReadBlocks struct {
}

type CancelWriteBlocks struct {
}

type ReadBlockOperation struct {
}

type SendLastBlockOperation struct {
	Msg []byte
}

type SendLastBlockHeaderOperation struct {
	Msg []byte
}
