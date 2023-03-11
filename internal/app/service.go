package app

import (
	"net/http"
)

type Service struct {
	Name   string
	server *http.Server
}
