package main

import (
	"github.com/izaakdale/service-ticket/internal/app"
)

func main() {
	// mailer.Send([]order.Ticket{
	// 	{
	// 		FirstName:  "izaak",
	// 		Surname:    "dale",
	// 		QrPath:     "testpath",
	// 		TicketType: "adult",
	// 	},
	// })

	app.Run()
}
