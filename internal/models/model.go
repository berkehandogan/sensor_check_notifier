package models

import "time"

type OpcSystemLogs struct {
	OpcServerId   int       `db:"opc_server_id"`
	Message       string    `db:"message"`
	Date          time.Time `db:"date"`
	TotalSensor   int       `db:"toplam_sensor"`
	SuccessSensor int       `db:"basarili_sensor"`
}

type TrendAnalysisServer struct {
	Id             int    `db:"id"`
	ServerEndPoint string `db:"server_endpoint"`
	ServerName     string `db:"server_name"`
}

type OpcToMail struct {
	OpcServerId int    `db:"opc_server_id"`
	ToMail      string `db:"to_mail"`
}
