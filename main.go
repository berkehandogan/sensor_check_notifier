package main

import (
	"fmt"
	"log"
	"sensor_check_notifier/internal/database"
	"sensor_check_notifier/internal/queries"
)

func main() {

	err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Veritabanı başlatılamadı: %v", err)
	}

	dbPool := database.GetDB()
	if dbPool == nil {
		log.Fatal("DB bağlantı havuzu alınamadı.")
		return
	}
	fmt.Println("Veritabanına başarıyla bağlanıldı ve bağlantı havuzu alındı.")
	queries.CheckSensorPerformance(dbPool)

}
