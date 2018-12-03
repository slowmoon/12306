package main

import (
	"fmt"
	"strings"
	"encoding/json"
)



type CaptchaCheckResponse struct {
	 ResultMessage string  `json:"result_message"`
	 ResultCode string    `json:"result_code"`
}


func main(){
  var msg = `{"result_message":"验证码校验成功","result_code":"4"}`
  var resp CaptchaCheckResponse
   json.NewDecoder(strings.NewReader(msg)).Decode(&resp)
   
   fmt.Println(resp)

}











