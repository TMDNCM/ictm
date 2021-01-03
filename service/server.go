package service

import (
	"fmt"
	"net/http"
)

type Dispatcher struct {
	Server http.Server
}

func NewDispatcher(address string, port uint) *Dispatcher {
	dispatcher := new(Dispatcher)
	dispatcher.Server.Addr = fmt.Sprintf("%s:%d", address, port)
}
