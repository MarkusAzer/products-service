package handler

import (
	"fmt"
	"time"

	"github.com/markus-azer/products-service/pkg/entity"

	"github.com/markus-azer/products-service/pkg/brand"
)

//MakeBrandHandlers make msg handlers
func MakeBrandHandlers(msgRepo brand.MessagesRepository, service brand.UseCase) {
	c := msgRepo.GetMessages()

	go func() {
		for msg := range c {
			switch msg.Type {
			case "BRAND_CREATED":
				fmt.Println("BRAND_CREATED", msg)
				b := &entity.Brand{
					ID:      entity.ID(msg.ID),
					Version: msg.Version,
					//Name:        msg.Payload["Name"],
					// Description: msg.Payload["Description"],
					// Slug:        msg.Payload["Slug"],
					// CreatedAt:   msg.Payload["CreatedAt"],
				}

				for key, item := range msg.Payload {
					value, _ := item.(string)
					fmt.Println("value =>>>>>>>>>", value)

					switch key {
					case "CreatedAt":
						layout := "2006-01-02T15:04:05.000Z"
						t, _ := time.Parse(layout, value)
						b.CreatedAt = time.Time(t)
						fmt.Println("t =>>>>>>>>>", t)
					case "Description":
						b.Description = value
					case "Name":
						b.Name = value
					case "Slug":
						b.Slug = value
					}
				}
				fmt.Println("=>>>>>>>>>", b)
				// brand := *entity.Brand{ID: msg.ID, Version: msg.Version, Name:}
				service.Create(b)
			case "BRAND_UPDATED":
				fmt.Println("BRAND_UPDATED")
			case "BRAND_DELETED":
				fmt.Println("BRAND_DELETED")
			default:
				fmt.Println("No handler")
			}
		}
	}()

}
