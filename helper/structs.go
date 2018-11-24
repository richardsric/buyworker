package helper

type ModifyDb struct {
	AffectedRows int64
	ErrorMsg     string
}

type RowSelect struct {
	Columns  map[string]interface{}
	ErrorMsg string
}

type psCoinInfo struct {
	P string
	S string
}
