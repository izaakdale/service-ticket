package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/izaakdale/lib/listener"
	"github.com/izaakdale/service-ticket/internal/mailer"
	"github.com/skip2/go-qrcode"

	"github.com/izaakdale/service-event-order/pkg/notifications"
	"github.com/izaakdale/service-event-order/pkg/proto/order"
)

func Process(ctx context.Context, m listener.Message) error {
	log.Printf("processing message from queue %s\n", m.MessageID)
	// decode message
	var payload notifications.OrderStoredPayload
	err := json.Unmarshal([]byte(m.Message), &payload)
	if err != nil {
		return err
	}
	// fetch order tickets
	o, err := client.GetOrder(ctx, &order.OrderRequest{OrderId: payload.OrderID})
	if err != nil {
		return err
	}

	for _, v := range o.Tickets {
		err := qrcode.WriteFile(v.QrPath, qrcode.Medium, 256, fmt.Sprintf("tmp/%s.jpg", v.TicketId))
		if err != nil {
			return err
		}
	}

	err = mailer.Send(o.Email, payload.OrderID, o.Tickets)
	if err != nil {
		return err
	}

	for _, v := range o.Tickets {
		err := os.Remove(fmt.Sprintf("tmp/%s.jpg", v.TicketId))
		if err != nil {
			log.Printf("%+v\n", err)
		}
	}

	return nil
}
