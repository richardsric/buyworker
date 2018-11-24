package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	w "github.com/richardsric/buyworker/worker"
)

func main() {

	http.HandleFunc("/", index)
	http.ListenAndServe(":6000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "iTradeCoin Buy Order Worker Service Is Running On Port 6000")
}

func init() {

	var name = "iTradeCoin Buy Order Update Worker"
	var version = "0.001 DEVEL"
	var developer = "iYochu Nig LTD"

	fmt.Println("App Name: ", name)
	fmt.Println("App Version: ", version)
	fmt.Println("Developer Name: ", developer)
	//Run the Default service settings inline.
	w.GetServericeSettings()
	go w.BuyUpdateWorker()
}
