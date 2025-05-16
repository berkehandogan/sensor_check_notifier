package queries

import (
	"context"
	"fmt"
	"log"
	"sensor_check_notifier/internal/models"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

// checkSensorPerformance fonksiyonu, sensörlerin başarı oranlarını kontrol eder.
func CheckSensorPerformance(dbPool *pgxpool.Pool) {
	fmt.Println("Sensör performans kontrolü başlatılıyor...")

	// 1. Şuan sadece "Sistem veri akışı tamamlandı." mesajına sahip logları çeker. Burayı şimdilik böyle bırakıyorum diğer mesajlar içinde ayarla.
	targetMessages := []string{
		"Sistem veri akışı tamamlandı.",
		"System data flow completed",
	}

	date := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC) // istediğin tarih
	dateNextDay := date.Add(24 * time.Hour)

	query := `
		SELECT opc_server_id, message, toplam_sensor, basarili_sensor, date
		FROM techupdb.techup.opc_system_logs
		WHERE message = ANY($1)
		AND date >= $2
		AND date < $3
	`
	rows, err := dbPool.Query(context.Background(), query, targetMessages, date, dateNextDay)
	if err != nil {
		log.Printf("opc_system_logs sorgulanırken hata: %v\n", err)
		return
	}
	defer rows.Close()

	var problematicServers []models.TrendAnalysisServer
	processedServerIDs := make(map[int]bool) // Aynı server ID için birden fazla uyarıyı engellemek için

	fmt.Println("İlgili log kayıtları işleniyor...")
	// 2. Her bir log kaydını işle
	for rows.Next() {
		var logEntry models.OpcSystemLogs
		err := rows.Scan(
			&logEntry.OpcServerId,
			&logEntry.Message,
			&logEntry.TotalSensor,
			&logEntry.SuccessSensor,
			&logEntry.Date,
		)
		if err != nil {
			log.Printf("opc_system_logs satırı okunurken hata: %v\n", err)
			continue // Bu satırı atla, sonrakine geç
		}

		successRatio := (float64(logEntry.SuccessSensor) / float64(logEntry.TotalSensor)) * 100.0

		//log.Printf("DEBUG: OpcServerId: %d, Başarılı: %d, Toplam: %d, Oran: %.2f%% (Tarih: %s)\n",
		//logEntry.OpcServerId, logEntry.SuccessSensor, logEntry.TotalSensor, successRatio, logEntry.Date)

		// 3. Başarı oranı %70'in altındaysa ilgili sunucu bilgilerini çek
		if successRatio < 70.0 {
			log.Printf("UYARI: OpcServerId %d için başarı oranı (%.2f%%) %%70'in altında (Tarih: %s).\n", logEntry.OpcServerId, successRatio, logEntry.Date)

			var serverInfo models.TrendAnalysisServer
			serverQuery := `
				SELECT id, server_endpoint, server_name
				FROM techupdb.techup.trendanalysis_servers
				WHERE id = $1
			`
			err := dbPool.QueryRow(context.Background(), serverQuery, logEntry.OpcServerId).Scan(
				&serverInfo.Id, // Bu zaten logEntry.OpcServerId ile aynı olacak
				&serverInfo.ServerEndPoint,
				&serverInfo.ServerName,
			)

			if err != nil {
				if err == pgx.ErrNoRows { // Veritabanında ilgili ID'ye sahip sunucu bulunamazsa
					log.Printf("UYARI: OpcServerId %d için trendanalysis_servers tablosunda eşleşen sunucu bulunamadı.\n", logEntry.OpcServerId)
				} else {
					log.Printf("trendanalysis_servers sorgulanırken hata (OpcServerId: %d): %v\n", logEntry.OpcServerId, err)
				}
				continue // Bu sunucu bilgisini atla
			}

			// Eğer bu server ID'yi daha önce işlemediysek listeye ekle
			if _, ok := processedServerIDs[serverInfo.Id]; !ok {
				problematicServers = append(problematicServers, serverInfo)
				processedServerIDs[serverInfo.Id] = true
			}
		}
	}

	if rows.Err() != nil { // Döngü sırasında bir hata oluştuysa
		log.Printf("opc_system_logs satırları işlenirken döngü sonrası hata: %v\n", rows.Err())
	}

	// 4. Problemli sunucuları listele
	if len(problematicServers) > 0 {
		fmt.Println("\n--- DİKKAT EDİLMESİ GEREKEN SUNUCULAR (Başarı Oranı < %70) ---")
		for _, server := range problematicServers {
			fmt.Printf("  Server ID: %d, Adı: '%s', Endpoint: '%s'\n", server.Id, server.ServerName, server.ServerEndPoint)
		}
		fmt.Println("-----------------------------------------------------------")
	} else {
		fmt.Println("\n--- Kontrol edilen 'Sistem veri akışı tamamlandı' loglarında başarı oranı %70'in altında olan sunucu bulunamadı. ---")
	}
}
