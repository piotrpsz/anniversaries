package events

type EventID uint

const (
	WindowOpened EventID = iota
	WindowClosed
	RefreshRequest
	PersonSelected
	PersonEdit
	AddNewPerson
	KeyEvent
)
