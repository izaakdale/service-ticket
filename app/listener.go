package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/izaakdale/lib/listener"

	"github.com/izaakdale/service-event-order/pkg/notifications"
	"github.com/izaakdale/service-event-order/pkg/proto/order"
)

func Process(m listener.Message) error {
	log.Printf("processing message from queue %s\n", m.MessageID)
	// decode message
	var payload notifications.OrderStoredPayload
	err := json.Unmarshal([]byte(m.Message), &payload)
	if err != nil {
		return err
	}
	// fetch order tickets
	o, err := client.GetOrder(context.Background(), &order.OrderRequest{OrderId: payload.OrderID})
	if err != nil {
		return err
	}

	log.Printf("%+v\n", o)
	// compose email
	// send
	return nil
}
