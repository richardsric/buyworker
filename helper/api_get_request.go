package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetHTTPRequest This is use to make http Get request. It returns byte an error
func GetHTTPRequest(url string) (bs []byte, err error) {
	fmt.Println("Getting HTTP Request From URL: " + url + "")

	res, err := http.Get(url)
	if (err) != nil {

		return nil, err
	}
	defer res.Body.Close()
	bs, err = ioutil.ReadAll(res.Body)
	if (err) != nil {

		return nil, err
	}

	res.Body.Close()
	return bs, nil
}
