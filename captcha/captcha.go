package captcha

import (
	"net/url"
	"strconv"
	"errors"
	"io/ioutil"
	"encoding/base64"
	"strings"
	"encoding/json"
	"regexp"
	"math/rand"
	"fmt"
	"time"
	"net/http"
)

var (
	pattern = "(?m:.*\\((.*)\\).*)"
)

const (
	width = 75 
	height = 80
	baseX = 0
	baseY = 28
)

var compile *regexp.Regexp
func init(){
	rand.Seed(time.Now().UnixNano())
    compile = regexp.MustCompile(pattern)
}

type CaptchaRequest struct{
	base string		`map:"base"`
	loginSite string  `map:"login_site"`
	module string	`map:"module"` 	
	rand  string    `map:"rand"`
	callback string  `map:"callback"`
	timestamp int64  `map:"_"`
}

func (c *CaptchaRequest)Name()string{
	return fmt.Sprintf("%s.jpg", c.callback)
}

type CaptchaResponse struct {
	ResultCode string  `json:"result_code"`
	ResultMessage string `json:"result_message"`
	Image string   `json:"image"`
}

func NewCaptchaResponse(message string) ( *CaptchaResponse, error){
		 result := compile.FindAllStringSubmatch(message, -1)
		 var resp CaptchaResponse
		 err := json.NewDecoder(strings.NewReader(result[0][1])).Decode(&resp)
		return &resp, err
}

func (c *CaptchaResponse)Succ()bool{
	return c.ResultCode == "0"
}

func (c *CaptchaResponse)Photo(name string){
	if !c.Succ(){
		panic("获取验证图片失败")
	}
	img, err := base64.StdEncoding.DecodeString(c.Image)
	if err !=nil {
		fmt.Println("error", err)
	}
	ioutil.WriteFile(name, img, 0755)   //图片写入指定文件
}


func NewCaptchaRequest(base string)*CaptchaRequest{
	now := Now()
	rand := rand.Int31n(1000000000)
	callback := fmt.Sprintf("%s%d_%d", "jQuery", rand, now)
	return &CaptchaRequest{
		base: base,
		loginSite: "E",
		module: "login",
		rand: "sjrand",
		callback: callback,
		timestamp: now,
	}
}

func (c *CaptchaRequest)Get(client *http.Client)(string, error){
	resp, err  := client.Get(c.Encode())
	if err != nil {
		fmt.Println("request error", err)
		return "", nil
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		fmt.Println("err read data", err)
		return "", err
	}
	capthaResp, err := NewCaptchaResponse(string(content))
	if err!=nil{
		fmt.Println("error" ,err)
		return "", err
	}
	 capthaResp.Photo(c.Name())
	 return c.Name(), nil
}

func (c *CaptchaRequest)Encode()string{
	return fmt.Sprintf("%s?login_site=%s&module=%s&rand=%s&callback=%s&_=%d",c.base, 
	 c.loginSite, c.module, c.rand, c.callback, c.timestamp)	
}


func Choose(pics string)(Answer, error){
	result := strings.Split(pics, ",")
	if len(result)==0 || len(result)>8{
		fmt.Println("验证码选取数量 %d有误！", len(pics))
		return nil, errors.New("invalid choose size!")
	}
	fmt.Println(result)
	var answer Answer
	for _, point := range result {
		iPoint, err := strconv.Atoi(point)
		if err !=nil {
			fmt.Println("error input", err)
			return nil, err
		}
		answer  = append(answer, NewPointer(iPoint))
	}
	return answer, nil
}

type Pointer struct{
	X, Y int
}

type Answer []Pointer

func(s Answer)String()string{
	var l string
	first := true 
	for _, point := range s{
		if first{
			l += fmt.Sprintf("%d,%d", point.X, point.Y) 
			first = false
			}else{
			l += fmt.Sprintf(",%d,%d", point.X, point.Y) 
		}
	}
	return l
}


func NewPointer(point int)Pointer{
	return Pointer{
		X:75/2 + ((point-1) % 4)*75,
		Y: 40+((point-1)/4)*80,
	}
}


type CaptchaCheck struct{
	base string
	callback string
	answer Answer
	rand string
	loginSite string
    timestamp int64
}


func NewCaptchaCheck(base string,  req *CaptchaRequest, answer Answer)*CaptchaCheck{
	return &CaptchaCheck{
		base: base,
		callback: req.callback,
		answer: answer,
		rand: req.rand,
		loginSite: req.loginSite,
		timestamp: req.timestamp +1,
	}	
}

func (c *CaptchaCheck)encode()string{
	return fmt.Sprintf("%s?callback=%s&answer=%s&rand=%s&login_site=%s&_=%d", c.base, c.callback, url.QueryEscape(c.answer.String()), c.rand, c.loginSite, c.timestamp)
}

func(c *CaptchaCheck)Check(client *http.Client)bool{
	var request  = c.encode()
	resp, err := client.Get(request)
	if err!= nil{
		fmt.Println("error request url", err)
		return false
	}

	defer resp.Body.Close()
	content , err := ioutil.ReadAll(resp.Body)
	fmt.Println("check response content:", string(content))
	if err !=nil{
		fmt.Println("read data error ", err)
		return false
	}
	checkResp, err := NewCaptchaCheckResponse(string(content))
	fmt.Println("checkresp", checkResp)
	if err !=nil {
		fmt.Println("error read response ", err)
		return false
	}
	return checkResp.Succ()
}

func Now()int64{
	return time.Now().UnixNano()/1e6
}

type CaptchaCheckResponse struct {
	 ResultMessage string  `json:"result_message"`
	 ResultCode string    `json:"result_code"`
}

func NewCaptchaCheckResponse(msg string)( *CaptchaCheckResponse, error){
	 result  := compile.FindAllStringSubmatch(msg, -1)
	 var resp CaptchaCheckResponse
	 err := json.NewDecoder(strings.NewReader(result[0][1])).Decode(&resp)
	 return &resp, err	
}

func(r *CaptchaCheckResponse)Succ()bool{
	return r.ResultCode =="4"
}
