package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Send(addr string, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil || len(buf) == 0 {
		fmt.Println("send json marshal err,or data len is 0 , data is:", data)
		return
	}
	fmt.Printf("send :%s,data :%s", addr, string(buf))
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
