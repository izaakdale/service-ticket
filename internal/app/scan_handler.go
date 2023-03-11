package app

import (
	"log"
	"net/http"

	"github.com/izaakdale/lib/response"
	"github.com/izaakdale/service-event-order/pkg/proto/order"
	"github.com/julienschmidt/httprouter"
)

func ScanHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	ticketID := params.ByName("id")
	log.Printf("handling scan of ticket %s\n", ticketID)

	resp, err := client.ScanTicket(r.Context(), &order.ScanRequest{
		TicketId: ticketID,
	})
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.NewError(err.Error()))
		return
	}

	response.WriteJson(w, http.StatusOK, resp)
	return
}
