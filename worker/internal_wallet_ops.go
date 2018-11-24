package worker

import (
	"encoding/json"
	"fmt"

	h "github.com/richardsric/buyworker/helper"
)

//GetInternalWalletBalFrmGateway calls gateway and gets the wallet balance.
func GetInternalWalletBalFrmGateway(apiKey string, exchangeID int64, accountID int64, coinName string, walletType string) InternalBalanceResponse {
	var res InternalBalanceResponse
	//"/getInternalBalance"
	gwendPoint := fmt.Sprintf("%s/getInternalBalance?apiKey=%s&eid=%v&aid=%v&coinName=%s&walletType=%s", GatewayURL, apiKey, exchangeID, accountID, coinName, walletType)
	body, err := h.GetHTTPRequest(gwendPoint)
	if err != nil {
		fmt.Println("GetBalanceFailed due to: ", err)
		res.Result = "error"
		res.Message = "GetBalanceFailed"
		return res
	}

	if len(body) == 0 {
		fmt.Println("GetInternalWalletBalFrmGateway: Nil Response")
		fmt.Println("Kindly Check Your Gateway Connection")
		res.Result = "error"
		res.Message = "GetBalanceFailed: No response"
		return res
	}
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)

	//	fmt.Printf("%+v", m)
	res.Result = m["result"].(string)
	res.Message = m["message"].(string)
	details := m["details"].(map[string]interface{})
	res.Details.AccountID = int64(h.GetType(details["account_id"]).(float64))
	res.Details.Available = h.GetType(details["available"]).(float64)
	res.Details.CoinName = h.GetType(details["coin_name"]).(string)
	res.Details.CurrentCapital = h.GetType(details["current_capital"]).(float64)
	res.Details.ExchangeID = int64(h.GetType(details["exchange_id"]).(float64))
	res.Details.MainCapital = h.GetType(details["main_capital"]).(float64)
	res.Details.Pending = h.GetType(details["pending"]).(float64)
	res.Details.Reserved = h.GetType(details["reserved"]).(float64)
	res.Details.Total = h.GetType(details["total"]).(float64)
	res.Details.Used = h.GetType(details["used"]).(float64)
	res.Details.WalletType = h.GetType(details["wallet_type"]).(string)

	return res
}

//UpdateInternalWalletBalToGateway calls gateway and sets the wallet balance.
func UpdateInternalWalletBalToGateway(apiKey string, exchangeID int64, accountID int64, coinName string, walletType string, amount float64, balanceName string, walletAction string) InternalBalanceUpdateResponse {
	var res InternalBalanceUpdateResponse
	// "/updateInternalBalance"
	gwendPoint := fmt.Sprintf("%s/updateInternalBalance?apiKey=%s&eid=%v&aid=%v&coinName=%s&walletType=%s&amount=%v&balanceName=%s&walletAction=%s", GatewayURL, apiKey, exchangeID, accountID, coinName, walletType, amount, balanceName, walletAction)
	body, err := h.GetHTTPRequest(gwendPoint)
	if err != nil {
		fmt.Println("UpdateBalanceFailed due to: ", err)
		res.Result = "error"
		res.Message = "UpdateBalanceFailed"
		return res
	}

	if len(body) == 0 {
		fmt.Println("UpdateInternalWalletBalToGateway: Nil Response")
		fmt.Println("Kindly Check Your Gateway Connection")
		res.Result = "error"
		res.Message = "UpdateBalanceFailed: No response"
		return res
	}

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)

	//	fmt.Printf("%+v", m)
	res.Result = m["result"].(string)
	res.Message = m["message"].(string)
	res.NewBalance = h.GetType(m["new_balance"]).(float64)
	//	fmt.Println(res)
	return res
}
