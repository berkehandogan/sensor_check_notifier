package main

import (
	"context"
	"fmt"
	"log"
	"sensor_check_notifier/internal/database"
)

func main() {

	fmt.Println("Hello, World!")
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

	var version string
	err = dbPool.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Sorgu çalıştırılırken hata: %v", err)
	}

	fmt.Printf("PostgreSQL Versiyonu: %s\n", version)

	fmt.Println("Test sorgusu başarıyla çalıştı!")
}
