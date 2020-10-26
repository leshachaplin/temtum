package operations

type Operation interface {
}

type ChangeActiveLieder struct {
}

type ExitOperation struct {
	Msg string
}

type SendRoleOperation struct {
	Msg string
}

