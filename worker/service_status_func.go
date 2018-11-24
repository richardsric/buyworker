package worker

import (
	"encoding/base64"
	"fmt"

	h "github.com/richardsric/buyworker/helper"
)

//SendTelegramIM this use to send HTML parsed  message to a telegram user.
func SendTelegramIM(chatID string, msg string) int64 {

	con, e := h.OpenConnection()
	if e != nil {
		return 0
	}
	defer con.Close()
	var mid interface{}
	q := `INSERT INTO telegram_messages(msgto, message) VALUES($1, $2) RETURNING messageid`
	e = con.Db.QueryRow(q, chatID, msg).Scan(&mid)
	if e != nil {
		return 0
	}
	if mid != nil {
		return mid.(int64)
	}

	return 0
}

//SendServiceStatusIM this use to send HTML parsed  message to a telegram user.
func SendServiceStatusIM(msg string) int64 {
	//BotKey for Error Reporting
	var adminTelegram = "430073910"

	con, e := h.OpenConnection()
	if e != nil {
		return 0
	}
	defer con.Close()
	var mid interface{}
	q := `INSERT INTO telegram_service_status_messages(msgto, message) VALUES($1, $2) RETURNING messageid`
	e = con.Db.QueryRow(q, adminTelegram, msg).Scan(&mid)
	if e != nil {
		return 0
	}
	if mid != nil {
		return mid.(int64)
	}

	return 0
}

//alertCall pushes alert to alerts service
func alertCall(chatid string, command string) bool {
	//encode command to base64
	commandEnc := base64.StdEncoding.EncodeToString([]byte(command))
	qs := fmt.Sprintf("%s/alerts?command=%s&chatid=%s", AlertsURL, commandEnc, chatid)
	//fmt.Println("Request String:", qs)
	body, err := h.GetHTTPRequest(qs)
	if err != nil {
		fmt.Println("AlertCall: could not reach alert service due to ", err)
		return false
	}
	//fmt.Println("AlertCall Response:", string(body))
	if string(body) == "ok" {
		return true
	}
	return false
}
