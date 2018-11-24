package worker

import (
	"encoding/json"
	"errors"
	"fmt"

	h "github.com/richardsric/buyworker/helper"
)

//GetOrderInfo gets order info
func GetOrderInfo(apiKey, exchangeID, accountID, orderID string) (orderDetail, error) {

	var oInfo orderDetail
	body, err := h.GetHTTPRequest(fmt.Sprintf("%s/getOrderInfo?apiKey=%s&uuid=%s&eid=%v&aid=%v", GatewayURL, apiKey, orderID, exchangeID, accountID))
	if err != nil {
		fmt.Println("GetOrderInfo: Error On Bittrex GetTicker Func due to ", err)
		return oInfo, err
	}
	var v OrderInfoResponse
	var exchangeInfo ExchangeInfo
	err = json.Unmarshal(body, &v)
	if err != nil {
		return oInfo, err
	}
	if v.Result == "error" {
		return oInfo, errors.New(v.Message)
	}
	err = json.Unmarshal(v.Details, &oInfo)
	if err != nil {
		return oInfo, err
	}
	err = json.Unmarshal(v.ExchangeInfo, &exchangeInfo)
	if err != nil {
		return oInfo, err
	}
	cRes := orderDetail{
		Market:            oInfo.Market,
		OrderType:         oInfo.OrderType,
		ActualQuantity:    oInfo.ActualQuantity,
		QuantityRemaining: oInfo.QuantityRemaining,
		ActualRate:        oInfo.ActualRate,
		OrderStatus:       oInfo.OrderStatus,
		Fee:               oInfo.Fee,
		OrderDate:         oInfo.OrderDate,
		Price:             oInfo.Price,
		PricePerUnit:      oInfo.PricePerUnit,
		Reserved:          oInfo.Reserved,
		Exchange:          exchangeInfo,
	}
	return cRes, nil
}
