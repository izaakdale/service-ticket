package app

import (
	"net/http"

	"github.com/izaakdale/lib/router"
)

func Router() http.Handler {
	r := router.New(
		router.WithRoute(http.MethodPost, "/qr/:id", ScanHandler),
	)
	return r
}
