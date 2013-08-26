package java2go

import (
    "fmt"
)

type Message struct {
    path string
}

func (msg *Message) PrintFuncs() {
    fmt.Println("Found", msg.path)
}

func NewMessage(path string) (*Message, error) {
    msg := new(Message)
    msg.path = path

    return msg, nil
}
