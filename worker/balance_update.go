package worker

import (
	"strconv"
	"strings"

	h "github.com/richardsric/buyworker/helper"
)

// WalletBalanceUpdate is use to update the internal wallet balance. It returns the updated value on success and -1 on failure Example data:
// ("btc", "capital", "SIMULATION", "credit", 1, 1, 0.2258)
func WalletBalanceUpdate(coiName, balanceName, walletType, action string, accountID int64, exchangeID int64, amount float64) float64 {

	coiName = strings.ToUpper(coiName)
	balanceName = strings.ToLower(balanceName)
	walletType = strings.ToUpper(walletType)
	action = strings.ToLower(action)

	//	fmt.Println(balanceName)
	queryStrig := "SELECT " + balanceName + " FROM internal_wallet_balances WHERE account_id = $1 AND wallet_type = $2 AND coin_name = $3 AND exchange_id = $4"
	//fmt.Println(queryStrig)
	sql := h.DBSelectRow(queryStrig, accountID, walletType, coiName, exchangeID)
	if sql.Columns[balanceName] == nil || sql.ErrorMsg != "" {
		return -1
	}
	bal := sql.Columns[balanceName].([]uint8)
	balance, _ := strconv.ParseFloat(string(bal), 64)
	if action == "credit" {
		finalBalance := balance + amount
		creditQuery := "UPDATE internal_wallet_balances SET " + balanceName + " = $1 WHERE account_id = $2 AND wallet_type = $3 AND coin_name = $4 AND exchange_id = $5"
		credit := h.DBModify(creditQuery, finalBalance, accountID, walletType, coiName, exchangeID)
		if credit.AffectedRows > 0 {
			return finalBalance
		}
	}

	if action == "debit" {
		finalBalance := balance - amount
		debitQuery := "UPDATE internal_wallet_balances SET " + balanceName + " = $1 WHERE account_id = $2 AND wallet_type = $3 AND coin_name = $4 AND exchange_id = $5"
		debit := h.DBModify(debitQuery, finalBalance, accountID, walletType, coiName, exchangeID)
		if debit.AffectedRows > 0 {
			return finalBalance
		}
	}

	return -1
}

// GetWalletBalance returns all the balance
func GetWalletBalance(coiName, walletType string, accountID, exchangeID int) balances {

	coiName = strings.ToUpper(coiName)
	//balanceName = strings.ToLower(balanceName)
	walletType = strings.ToUpper(walletType)

	//fmt.Println(balanceName)
	queryStrig := "SELECT capital, pending, available, reserved, used FROM internal_wallet_balances WHERE account_id = $1 AND wallet_type = $2 AND coin_name = $3 AND exchange_id = $4"
	//fmt.Println(queryStrig)
	sql := h.DBSelectRow(queryStrig, accountID, walletType, coiName, exchangeID)
	if sql.ErrorMsg != "" {
		return balances{
			Msg: sql.ErrorMsg,
		}

	}
	cap := sql.Columns["capital"].([]uint8)
	capital, _ := strconv.ParseFloat(string(cap), 64)

	pen := sql.Columns["pending"].([]uint8)
	pending, _ := strconv.ParseFloat(string(pen), 64)

	ava := sql.Columns["available"].([]uint8)
	available, _ := strconv.ParseFloat(string(ava), 64)

	res := sql.Columns["reserved"].([]uint8)
	reserved, _ := strconv.ParseFloat(string(res), 64)

	use := sql.Columns["used"].([]uint8)
	used, _ := strconv.ParseFloat(string(use), 64)
	return balances{
		Capital:   capital,
		Pending:   pending,
		Available: available,
		Reserved:  reserved,
		Used:      used,
	}
}
