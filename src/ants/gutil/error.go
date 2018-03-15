package gutil

import "fmt"

//属于error
type IError interface {
	Error() string
	MessageID() int
}

type MsgError struct {
	IError
	code    int
	message string
}

func NewError(message string) IError {
	return NewMsgError(0, message)
}

func NewMsgError(code int, message string) IError {
	this := new(MsgError)
	this.SetError(code, message)
	return this
}

func (this *MsgError) SetError(code int, message string) {
	this.code = code
	this.message = message
}

func (this *MsgError) Error() string {
	return fmt.Sprintf("Msg=%d,Error=", this.code, this.message)
}

func (this *MsgError) MessageID() int {
	return this.code
}
