package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

//SendJSON .
func SendJSON(url string, headers map[string]string, bodyTosend []byte) ([]byte, int) {
	client := http.Client{}
	requets, err := http.NewRequest("POST", url,  bytes.NewBuffer(bodyTosend))
	ValidError(err)
	for key, value := range headers {
		requets.Header.Set(key, value)
		requets.Header.Set(key, value)
	}
	response, err := client.Do(requets)
	ValidError(err)
	defer response.Body.Close()
	bodyRS, err := ioutil.ReadAll(response.Body)
	ValidError(err)
	return bodyRS, response.StatusCode
}
