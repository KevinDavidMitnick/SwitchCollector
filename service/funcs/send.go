package funcs

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func PushToFalcon(addr string, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil || len(buf) == 0 {
		log.Println("send json marshal err,or data len is 0 , data is:", data)
		return
	}
	log.Printf("send :%s,data :%s", addr, string(buf))
	request, _ := http.NewRequest("POST", addr, bytes.NewBuffer(buf))
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode/100 != 2 {
			log.Println("reponse err")
		}
	}
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
			log.Println("reponse err")
		}
		return ioutil.ReadAll(resp.Body)
	}
	return nil, err
}
