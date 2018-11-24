package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	s "strings"
	"sync"
	"time"

	h "github.com/richardsric/buyworker/helper"
)

//GatewayURL HOLDS GATEWAY URL
var GatewayURL string

//TelegramURL HOLDS Telegram URL
var TelegramURL string

//ServicePort holds service port for reporting
var ServicePort int64

//MaxBuyProcs holds maximum number of concurrent processes to start for the worker
var MaxBuyProcs int64

//AlertsURL stores alerts URL
var AlertsURL string

//RequestTimeout holds the timeinterval for the service to timeout when it places an request to gateway
var RequestTimeout time.Duration
var mutex = &sync.Mutex{}
var job dbJob
var jobControl = make(chan int, MaxBuyProcs)

//BuyUpdateWorkService Starts and controls the service
func BuyUpdateWorkService() {
	for {

		BuyUpdateWorker()
		time.Sleep(time.Millisecond * 200)
	}
}

//GetServericeSettings populates the default settings needed by the server
func GetServericeSettings() {
	fmt.Println("Settings Service Defaults From DB")
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("GetServericeSettings: DB conenction could not be established")
		return
	}
	var gwurl, sp, reqtimeout, maxbuyprocs, tgurl, alurl interface{}
	defer con.Close()
	qr := "SELECT gateway_url,service_port,request_timeout,max_buy_procs,telegram_url,alerts_url FROM buy_worker_settings LIMIT 1"
	err = con.Db.QueryRow(qr).Scan(&gwurl, &sp, &reqtimeout, &maxbuyprocs, &alurl)
	if err != nil {
		fmt.Println("GetServericeSettings: could not get and assign server operation settings")
		return
	}
	if gwurl != nil {
		mutex.Lock()
		GatewayURL = gwurl.(string)
		mutex.Unlock()
	}
	if sp != nil {
		mutex.Lock()
		ServicePort = sp.(int64)
		mutex.Unlock()
	}
	if reqtimeout != nil {
		mutex.Lock()
		RequestTimeout = reqtimeout.(time.Duration)
		mutex.Unlock()
	}
	if maxbuyprocs != nil {
		mutex.Lock()
		MaxBuyProcs = maxbuyprocs.(int64)
		mutex.Unlock()
	}
	if tgurl != nil {
		mutex.Lock()
		TelegramURL = tgurl.(string)
		mutex.Unlock()
	}
	if alurl != nil {
		mutex.Lock()
		AlertsURL = alurl.(string)
		mutex.Unlock()
	}

}

// BuyUpdateWorker is func to update buy order
func BuyUpdateWorker() {
	fmt.Println("Starting buyUpdate Worker")
	con, err := h.OpenConnection()
	if err != nil {
		//return err
		fmt.Println("BuyUpdateWorker: DB conenction failed:", err)
		return
	}
	defer con.Close()
	queryString := `SELECT market,order_number,account_id,exchange_id,akey,work_started_on,
	buy_order_timeout,partialbuy_timeout,partialbuy_timeout_pl,partialbuy_detected_on,ask_bid,
	jobid,update_hash,order_type,actual_rate,cost,quantity,actual_quantity,buy_capital,trade_mode,
	instance_id,trade_profile 
	 FROM buy_orders WHERE work_status = 0`
	//re := helper.DBSelectRow(queryString)
	orderRows, err := con.Db.Query(queryString)
	if err != nil {
		fmt.Println("buy_update_worker.go: BuyUpdateWorker func - error in selecting buy_orders due to ", err)
		return
	}

	for orderRows.Next() {
		var instanceid, tradeprofile, trademode, buycapital, ordertype, actualrate, cost, quantity, actualquantity, orderid, accountid, exchangeid, apikey, markt, starttime, partialbuydetectedtime, buyordertimeout, partialbuytimeout, partialbuytimeoutpl, askbid, jobid, updatehash interface{}

		err = orderRows.Scan(&markt, &orderid, &accountid, &exchangeid, &apikey, &starttime, &buyordertimeout,
			&partialbuytimeout, &partialbuytimeoutpl, &partialbuydetectedtime, &askbid, &jobid,
			&updatehash, &ordertype, &actualrate, &cost, &quantity, &actualquantity, &buycapital, &trademode, &instanceid, &tradeprofile)
		if err != nil {
			if s.Contains(fmt.Sprintf("%v", err), "no rows") != true {
				fmt.Println("buy_update_worker.go: BuyUpdateWorker func - Row Scan Failed Due To: ", err)
			}
			return
		}

		//check for nill cases
		if markt != nil {
			job.market = markt.(string)
		}
		if orderid != nil {
			job.orderID = orderid.(string)
		}
		if accountid != nil {
			job.accountID = accountid.(int64)
		}
		if exchangeid != nil {
			job.exchangeID = exchangeid.(int64)
		}
		if apikey != nil {
			job.apiKey = apikey.(string)
		}

		if starttime != nil {
			job.startTime = starttime.(time.Time)
		}
		if buyordertimeout != nil {
			job.buyOrderTimeout = buyordertimeout.(int64)
		}
		if partialbuytimeout != nil {
			job.partailBuyTimeout = partialbuytimeout.(int64)
		}
		if partialbuytimeoutpl != nil {
			job.partailBuyTimeoutPl = partialbuytimeoutpl.(float64)
		}
		if partialbuydetectedtime != nil {
			job.partialBuyDetectedTime = partialbuydetectedtime.(time.Time)
		}
		if askbid != nil {
			job.askBid = askbid.(float64)
		}
		if jobid != nil {
			job.jobID = jobid.(int64)
		}
		if updatehash != nil {
			job.updateHash = updatehash.(string)
		}
		if ordertype != nil {
			job.orderType = ordertype.(string)
		}
		if actualrate != nil {
			job.actualRate = actualrate.(float64)
		}
		if cost != nil {
			job.cost = cost.(float64)
		}
		if quantity != nil {
			job.quantity = quantity.(float64)
		}
		if actualquantity != nil {
			job.actualQuantity = actualquantity.(float64)
		}
		if buycapital != nil {
			job.buyCapital = buycapital.(float64)
		}
		if trademode != nil {
			job.tradeMode = trademode.(string)
		}
		if instanceid != nil {
			job.instanceID = instanceid.(int64)
		}
		if tradeprofile != nil {
			job.tradeProfile = tradeprofile.(string)
		}
		//end check for nil cases
		//call routine to handle job
		var i int64
		for i = 1; i <= MaxBuyProcs; i++ {
			jobControl <- 1
			go processJob(job, jobControl)
		}
	} //end of data loop
}

func isBuyOrderTimedOut(timeStarted time.Time, buyTimeOut int64) bool {
	dur := time.Since(timeStarted).Seconds()
	fmt.Println("duration in seconds is ", dur)
	if dur > float64(buyTimeOut) {
		return true
	}
	return false
}

func isPartialOrderTimedOut(timeStarted time.Time, partialOrderTimeOut int64) bool {
	dur := time.Since(timeStarted).Seconds()
	fmt.Println("duration in seconds is ", dur)
	if dur > float64(partialOrderTimeOut) {
		return true
	}
	return false
}

func isPartailOrderPlMet(timeOutPl float64, market string, exchangeID int64, askbid float64) bool {
	mInfo := getAskBid(market, exchangeID)
	//get profit and round it up to 4db
	profitLoss := h.Round((((mInfo.Bid - askbid) / askbid) * 100), 4, "u")

	if profitLoss > timeOutPl {
		return true
	}
	return false
}

//CancelOrder fulfils cancel request
func CancelOrder(orderID string, apiKey string, accountID int64, exchangeID int64) (bool, error) {
	var oInfo OrderResponse
	body, err := h.GetHTTPRequest(fmt.Sprintf("%s/cancelOrder?apiKey=%s&uuid=%s&eid=%v&aid=%v", GatewayURL, apiKey, orderID, exchangeID, accountID))

	fmt.Println(string(body))
	if err != nil {
		fmt.Println("buy_update_worker.go: Error On GetHTTPRequest for cancel order in Partial fill status due to ", err)
		return false, err
	}
	err = json.Unmarshal(body, &oInfo)
	if err != nil {
		fmt.Println("buy_update_worker.go: BuyUpdateWorker func - orderResponse unmarshalling error due to ", err)
		return false, err
	}
	if oInfo.Result == "error" {
		fmt.Println("got error in trying to cancel order due to ", oInfo.Message)
		return false, errors.New(oInfo.Message)
	}
	return true, nil
}

func getAskBid(market string, exchangeID int64) *AskBid {
	var marketInfo AskBid
	url := fmt.Sprintf("%s/pair/price?pair=%s&eid=%v", GatewayURL, market, exchangeID)
	body, err := h.GetHTTPRequest(url)
	//fmt.Println(string(body))
	if err != nil {
		fmt.Println("buy_update_worker.go: getAskBid Error On GetHTTPRequest for cancel order in Partial fill status due to ", err)
		return nil
	}
	err = json.Unmarshal(body, &marketInfo)
	if err != nil {
		fmt.Println("buy_update_worker.go: getAskBid func - orderResponse unmarshalling error due to ", err)
		return nil
	}
	return &marketInfo
}

func processDBUYJob(job dbJob) {
	//deferred order. Check the price if it is conducieve and exit.
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("processDBUYJob: DB conenction could not be established")
		return
	}
	defer con.Close()
	//check timeout
	if isBuyOrderTimedOut(job.startTime, job.buyOrderTimeout) {
		//order timed out. mark it as done and revert funds
		qs := "UPDATE buy_orders SET work_status = $1, order_status = $2, order_date = $3 WHERE jobid = $4 AND update_hash = $5"
		res, err := con.Db.Exec(qs, 1, "CANCELED", time.Now(), job.jobID, job.updateHash)
		if err != nil {
			fmt.Println("BuyUpdateWorker: DBUY: failed to update buy_orders table for CANCELED due to ", err)
			return
		}
		aRow, _ := res.RowsAffected()
		if aRow > 0 {

			//DBUY Cancel was successful. now reverse funds.
			coinName := h.PsCoin(job.market).P
			debitResult := UpdateInternalWalletBalToGateway(job.apiKey, job.exchangeID, job.accountID, coinName, job.tradeMode, job.buyCapital, "used", "debit")
			if debitResult.Result == "error" {
				//debit failed
				return
			} else if debitResult.Result == "success" {
				//Fund has bn succesfully reversed
				//Send Message to Order Owner

				return
			}
		}
		//order update did not update the row in the timeout
		return
	}
	//order did not time out...check if it has reached BUY condition
	//get ask price
	marketPrice := getAskBid(job.market, job.exchangeID)

	if marketPrice.Ask <= job.askBid {
		//condition has been met. place IBUY order
		if placeIBUYOrder(job) {

			//if order is successful, then disable the existing jobid
		}
		return
	}

}

func processIBUYJob(job dbJob) {
	con, err := h.OpenConnection()
	if err != nil {
		fmt.Println("processIBUYJob: DB conenction could not be established")
		return
	}
	defer con.Close()
	body, err := h.GetHTTPRequest(fmt.Sprintf("%s/getOrderInfo?apiKey=%s&uuid=%s&eid=%v&aid=%v", GatewayURL, job.apiKey, job.orderID, job.exchangeID, job.accountID))
	//fmt.Println(string(body))
	if err != nil {
		fmt.Println("buy_update_worker.go: Error On Bittrex GetTicker Func due to ", err)
		return
	}
	// unmarshal the json response.
	var m OrderInfo

	err = json.Unmarshal(body, &m)
	if err != nil {
		//panic(err)
		fmt.Println("buy_update_worker.go: BuyUpdateWorker func - getOrderInfo unmarshalling error due to ", err)
	}

	//check for order_status
	if m.OrderStatus == "COMPLETED" {
		//Status_COMPLETED(con,orderID,accountID,exchangeID,m.Market,m.OrderStatus,m.ActualRate,m.OrderDate)
		//update other parameters of the order and check if order needs refactoring.
		if job.cost == 0 {
			//first time process. Update the orderInfo parameters.Order completed before any pass of the buy update worker

		}
		queryString := "UPDATE buy_orders SET work_status = $1, order_status = $2, order_date = $3 WHERE jobid = $4 AND update_hash = $5"
		res, err := con.Db.Exec(queryString, 1, "COMPLETED", m.OrderDate, job.jobID, job.updateHash)
		if err != nil {
			fmt.Println("failed to update buy_orders table for COMPLETED due to ", err)
		}
		aRow, _ := res.RowsAffected()
		if aRow != 1 {
			fmt.Println("affected row for COMPLETED status not equal to 1")
		}
	} else if m.OrderStatus == "OPEN" {
		if check := isBuyOrderTimedOut(job.startTime, job.buyOrderTimeout); check {
			//cancel the order
			checkCancel, err := CancelOrder(job.orderID, job.apiKey, job.accountID, job.exchangeID)
			if err != nil || checkCancel == false {
				fmt.Println("IsPartialOrderTimedOut cancel failed due to", err)
			}
			//fmt.Println("order successfully canceled")
		} else {
			queryString := "UPDATE buy_orders SET order_status = $2, order_date = $3 WHERE jobid = $4 AND update_hash = $5"
			res, err := con.Db.Exec(queryString, "OPEN", m.OrderDate, job.jobID, job.updateHash)
			if err != nil {
				fmt.Println("failed to update buy_orders table for OPEN due to ", err)
			}
			aRow, _ := res.RowsAffected()
			if aRow != 1 {
				fmt.Println("affected row for OPEN status not equal to 1")
			}
		}
	} else if m.OrderStatus == "CANCELED" {
		queryString := "UPDATE buy_orders SET work_status = $1, order_status = $2, order_date = $3 WHERE jobid = $4 AND update_hash = $5"
		res, err := con.Db.Exec(queryString, 1, "CANCELED", m.OrderDate, job.jobID, job.updateHash)
		if err != nil {
			fmt.Println("failed to update buy_orders table for CANCELED due to ", err)
		}
		aRow, _ := res.RowsAffected()
		if aRow != 1 {
			fmt.Println("affected row for CANCELED status not equal to 1")
		}
	} else if m.OrderStatus == "PARTIAL_FILL" {
		//first update the partialbuy_detected_on column to currenc time
		queryString := "UPDATE buy_orders SET partialbuy_detected_on = $1 WHERE jobid = $2 AND update_hash = $3"
		res, err := con.Db.Exec(queryString, time.Now(), job.jobID, job.updateHash)
		if err != nil {
			fmt.Println("failed to update buy_orders table due to ", err)
		}
		aRow, _ := res.RowsAffected()
		if aRow != 1 {
			fmt.Println("updating partialbuy_detected_on column affected row not equal to 1")
		}
		if check := isPartialOrderTimedOut(job.partialBuyDetectedTime, job.partailBuyTimeout); check {
			//cancel order
			checkCancel, err := CancelOrder(job.orderID, job.apiKey, job.accountID, job.exchangeID)
			if err != nil || checkCancel == false {
				fmt.Println("IsPartialOrderTimedOut cancel failed due to", err)
			}
		}
		//check if partialFill_pl is met
		if checkPartialPl := isPartailOrderPlMet(job.partailBuyTimeoutPl, job.market, job.exchangeID, job.askBid); checkPartialPl {
			//TODO calculate the unused capital and place a sell order with the purcahsed secondary coin
			if m.QuantityRemaining == 0 {
				fmt.Println("quantity remaining is equal to zero")
			} else {
				purchasedSecCoin := m.ActualQuantity - m.QuantityRemaining
				amountUsed := m.PricePerUnit * (m.ActualQuantity - m.QuantityRemaining)
				amountUnused := m.PricePerUnit * m.QuantityRemaining

				//TODO: call sell limit with the purchasedSecCoin as the quantity

				fmt.Printf("purchasedSecCoin - %v\t amountUsed - %v\t amountUnused - %v \n", purchasedSecCoin, amountUsed, amountUnused)
			}
		}
	} else if m.OrderStatus == "PARTIAL_CANCELED" {
		queryString := "UPDATE buy_orders SET work_status = $1, order_status = $2, order_date = $3 WHERE order_number = $4 AND account_id = $5 AND exchange_id = $6"
		res, err := con.Db.Exec(queryString, 1, "PARTIAL_CANCELED", m.OrderDate, job.orderID, job.accountID, job.exchangeID)
		if err != nil {
			fmt.Println("failed to update buy_orders table for PARTIAL_CANCELED due to ", err)
		}
		aRow, _ := res.RowsAffected()
		if aRow != 1 {
			fmt.Println("affected row for PARTIAL_CANCELED status not equal to 1")
		}
	} else {
		fmt.Println("unknown order status gotten as ", m.OrderStatus)
	}
} //end of IBUY loop

func processJob(job dbJob, jobControl chan int) {

	if job.orderType == "DBUY" {
		processDBUYJob(job)
	}
	if job.orderType == "IBUY" {
		processIBUYJob(job)
	}
	//free job slot
	<-jobControl
}

//placeIBUYOrder places an IBUY order on gateway
func placeIBUYOrder(job dbJob) bool {
	//call wallet funcs from here

	return true
}
