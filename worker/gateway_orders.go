package worker

import (
	"encoding/json"
	"fmt"

	h "github.com/richardsric/buyworker/helper"
)

//IBuyOrder performs the gateway call to execute the instruction for Instant Buy
func IBuyOrder(eID, aID, rate, Qty, key, market, timeOut string) OrderResult {
	/*
		This timeout is timeout for BUY. If it doesn't fill up or partial fill, and dat timeout passes, it will cancel
		Gateway places order, inserts in database with timeout and returns order ID.
		Actually, d IBUY should hv another parameter which is BuyOrderTabelID
		Actually, on second thought, d IBUY should not hv dat rowID. It should hv just timeout. Bcos it is d worker dat uses d rowID.
		So the only param to be added is d timeout on d IBUY.

		So, IBUY takes all d params dat d BUYORDER takes now in addition to timetout
		DBUY takes same parameters.
		Difference is dat DBUY does not leave gatway

		IBUY: works like d normal buy on gatway. But takes timeout param.
		Places order on exchange snd inserts in buy orders table and returns order number and message

	*/

	res := OrderResult{}

	//url := "" + baseURL + "/buyOrder?market=" + market + "&quantity=" + Qty + "&rate=" + rate + "&eid=" + eID + "&apiKey=" + key + "&aid=" + aID + ""
	url := "" + GatewayURL + "IBUY?market=" + market + "&quantity=" + Qty + "&rate=" + rate + "&eid=" + eID + "&apiKey=" + key + "&aid=" + aID + "&timeout=" + timeOut + " "

	body, err := h.GetHTTPRequest(url)
	if err != nil {
		fmt.Println("IBuyOrder Failed due to ", err)
		result := OrderResult{
			Message:     "Sorry I could not get any response for you at the moment. \n Try again later",
			OrderNumber: "",
		}
		return result
	}

	if len(body) == 0 {
		fmt.Println("IBuyOrder, Nil Response Gotten From The Request For Sell Order")
		fmt.Println("Kindly Check Your Internet Connection")
		result := OrderResult{
			Message:     "Sorry I could not get any response for you at the moment. \n Try again later",
			OrderNumber: "",
		}
		return result
	}

	// unmarshal the json response.
	var m map[string]interface{}

	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println("IBuyOrder:", err)
	}

	if m == nil {
		fmt.Println("IBuyOrder: Invalid Response Received From The Request", string(body))
		result := OrderResult{
			Message:     "Sorry I could not get any response for you at the moment. \n Try again later",
			OrderNumber: "",
		}
		return result
	}

	result := m["result"]
	message := m["message"]
	orderNo := m["order_number"]

	if result == "error" {
		fmt.Println("IBuyOrder:Order For buy Encountred The Following error: ", message)
		//return false, message.(string)
		result := OrderResult{
			Message:     message.(string),
			OrderNumber: "",
		}
		return result
	}
	if result == "success" {
		//fmt.Println("Order For Buy Placed Order ID: ", orderNo)
		//return true, orderNo.(string)
		result := OrderResult{
			Message:     "",
			OrderNumber: orderNo.(string),
		}
		return result

	}

	//return false, "I could Not process Your buy request.\nTry again later"
	res = OrderResult{
		Message:     "I could Not process Your buy request.\nTry again later",
		OrderNumber: "",
	}
	return res
}

//DBuyOrder executes instruction for a deferred order on gateway
func DBuyOrder(eID, aID, rate, Qty, key, market, timeOut string) OrderRowID {
	/*
	   it should be placed by d keyword DBUY
	   It will take normal params as d other ones except that it is not placed on d exchange. It is placed on d gatway, yes.
	   But it doesn't leave d gateway.
	   When we place a DBUY call to d gatway, gatway accepts and inserts into d database and returns d rowID of d insertion.
	   The DBUY simply inserts into d database and returns rowID.
	   DBUY takes same parameters.
	   	Difference is dat DBUY does not leave gatway

	   	DBUY: has same params with IBUY.
	   On gateway, it inserts order params in DB and returns d insert ID.

	   [3:06 PM, 10/17/2017] +234 810 069 9833: The Dbuy goes to the gateway  then insert into buy order table then waits when the price is reached?
	   [5:31 PM, 10/17/2017] Boss: It inserts into d order table and dat finishes.
	   [5:31 PM, 10/17/2017] Boss: The waiting and checking is done by d buy update worker.
	*/

	res := OrderRowID{}

	//url := "" + baseURL + "/buyOrder?market=" + market + "&quantity=" + Qty + "&rate=" + rate + "&eid=" + eID + "&apiKey=" + key + "&aid=" + aID + ""
	url := "" + GatewayURL + "DBUY?market=" + market + "&quantity=" + Qty + "&rate=" + rate + "&eid=" + eID + "&apiKey=" + key + "&aid=" + aID + "&timeout=" + timeOut + " "

	body, err := h.GetHTTPRequest(url)
	if err != nil {
		fmt.Println("DBuyOrder Failed due to ", err)
		result := OrderRowID{
			Message: "Sorry I could not get any response for you at the moment. \n Try again later",
			RowID:   0,
		}
		return result
	}

	if len(body) == 0 {
		fmt.Println("DBuyOrder: Nil Response Gotten From The Request For Sell Order")
		fmt.Println("Kindly Check Your Internet Connection")
		result := OrderRowID{
			Message: "Sorry I could not get any response for you at the moment. \n Try again later",
			RowID:   0,
		}
		return result
	}

	// unmarshal the json response.
	var m map[string]interface{}

	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println("DBuyOrder: unmarshal", err)
	}

	if m == nil {
		fmt.Println("DBuyOrder. Invalid Response Received From The Request", string(body))
		result := OrderRowID{
			Message: "Sorry I could not get any response for you at the moment. \n Try again later",
			RowID:   0,
		}
		return result
	}

	result := m["result"]
	message := m["message"]
	rowID := m["rowID"]

	if result == "error" {
		fmt.Println("DBuyOrder Encountred The Following error: ", message)
		//return false, message.(string)
		result := OrderRowID{
			Message: message.(string),
			RowID:   0,
		}
		return result
	}
	if result == "success" {
		//fmt.Println("Order For Buy Placed Order ID: ", rowID)
		//return true, orderNo.(string)
		result := OrderRowID{
			Message: "",
			RowID:   rowID.(int),
		}
		return result

	}

	//return false, "I could Not process Your buy request.\nTry again later"
	res = OrderRowID{
		Message: "I could Not process Your buy request.\nTry again later",
		RowID:   0,
	}
	return res

}
