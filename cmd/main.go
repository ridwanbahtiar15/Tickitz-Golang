package main

import (
	"gilangrizaltin/Backend_Golang/internal/routers"
	"gilangrizaltin/Backend_Golang/pkg"
	"log"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}
func main() {
	database, err := pkg.PostgreSQLDB()
	if err != nil {
		log.Fatal(err)
	}
	routers := routers.New(database)
	server := pkg.Server(routers)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
