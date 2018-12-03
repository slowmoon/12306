package login

import (
	"io/ioutil"
	"strings"
	"encoding/json"
	"net/url"
	"fmt"
	"errors"
	"net/http"
	"spider/captcha"
)

type Login struct {
	base string
	username string
	password string
	appid    string
	answer   captcha.Answer
	method   string
}

func NewLogin(base ,name, pswd  string, answer captcha.Answer)*Login{
	return &Login{
		base: base,
		username: name,
		password: pswd,
		appid: "otn",
		answer: answer,
		method: http.MethodPost,
	}
}

type LoginResponse struct {
	ResultCode int  `json:"result_code"`
	ResultMessage string `json:"result_message"`
}

func NewLoginResponse(msg string)(*LoginResponse, error){
	var resp LoginResponse
	err := json.NewDecoder(strings.NewReader(msg)).Decode(&resp)
	return &resp, err
}	

//post 传输数据
func (r *Login)Post(client *http.Client)(string, error){
	  if r.base == "" {
		  return "", errors.New("login address must not empty")
	  }

	  data := url.Values{}
	  data.Add("username", r.username)
	  data.Add("password", r.password)
	  data.Add("appid", r.appid)
	  data.Add("answer", r.answer.String())

	 result, err :=  client.PostForm(r.base, data)
	 if err != nil {
		 fmt.Println("post login info error", err)
		 return "", err
	 }
	 defer result.Body.Close()
	 content ,err := ioutil.ReadAll(result.Body)
	 if err !=nil {
		 fmt.Println("fail to read message", err)
		 return "", err
	 }
	 resp, err  := NewLoginResponse(string(content))
	 if err !=nil {
		 fmt.Println("parse message error", err)
		 return "", err
	 }
	 return resp.ResultMessage, nil
}



