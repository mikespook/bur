package bur

type State struct {
	Tx, Rx float64
}

type User struct {
	Name, Password string
	State          UserState
}

type UserState struct {
	Login int
	State
}
