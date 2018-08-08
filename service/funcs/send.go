package funcs

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func PushToFalcon(addr string, buf []byte) error {
	log.Printf("send :%s,data :%s", addr, bytes.NewBuffer(buf).String())
	request, _ := http.NewRequest("POST", addr, bytes.NewBuffer(buf))
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func GetData(addr string) ([]byte, error) {
	log.Printf("send :%s", addr)
	request, _ := http.NewRequest("GET", addr, nil)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode/100 != 2 {
			return nil, fmt.Errorf("reponse err")
		}
		return ioutil.ReadAll(resp.Body)
	}
	return nil, err
}
