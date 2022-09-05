package main

import (
	"fmt"

	"customer-partner/internal/db"
	"customer-partner/internal/domain"
	"customer-partner/internal/web"
)

func main() {
	fmt.Println("Starting Server")
	repo := db.NewPartnerInMemoryRepository()
	service := domain.NewPartnerService(repo)
	api := web.NewPartnerAPI(service)
	api.ListenAndServe()
}
