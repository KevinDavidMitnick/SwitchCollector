package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Send(addr string, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil {
		fmt.Println("send err , data is:", data)
		return
	}
	fmt.Println("send :%s,data :%s", addr, string(buf))
	request, _ := http.NewRequest("POST", addr, bytes.NewBuffer(buf))
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("TIMEOUT", "10")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode/100 != 2 {
			fmt.Println("reponse err")
		}
	}
}
