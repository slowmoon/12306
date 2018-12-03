package main

import (
	"io/ioutil"
	"spider/login"
	"net/http/cookiejar"
	"fmt"
	"net/http"
	"spider/captcha"
)

var (
	captchaUrl = "https://kyfw.12306.cn/passport/captcha/captcha-image64"
	captchaCheck = "https://kyfw.12306.cn/passport/captcha/captcha-check"
	loginUrl = "https://kyfw.12306.cn/passport/web/login"
	loginRedirect = "https://kyfw.12306.cn/otn/login/userLogin"

)

func main() {
	

	//打开浏览器

	var client = http.DefaultClient
	 jar , err := cookiejar.New(nil)
	 if err!=nil{
		fmt.Println("initilize cookie jar error", err)
		return
	}
	client.Jar = jar

	//申请验证码
	var req = captcha.NewCaptchaRequest(captchaUrl)
	name ,err := req.Get(client)   //验证码的名称
	fmt.Println(name)

	fmt.Println("请输入参数选择的图片 [1-8]:")
	var images string
	fmt.Scanf("%s", &images)
	answer , err := captcha.Choose(images)
	if err != nil {
		fmt.Println("choose image error!", err)
		return
	}
	//验证验证码
	check := captcha.NewCaptchaCheck(captchaCheck, req, answer)
	checkResult := check.Check(client)
	if !checkResult{
		fmt.Println("captcha check error! please check again")
		return
	}
	//输入用户名账号和密码
	fmt.Println("请输入用户名和密码 :")
	var username, password string
	fmt.Scanf("%s%s", &username, &password)
   loginReq := login.NewLogin(loginUrl, username, password, answer)

   //登录名和密码进行验证
   loginResult, err := loginReq.Post(client)
   if err != nil{
	   fmt.Println("login error !", err)
	   return
   }
   fmt.Println("login result:", loginResult)

   redirectResult, err := client.Get(loginRedirect)
   if err !=nil{
	   fmt.Println("redirect error", err)
		return
	}
   defer redirectResult.Body.Close()

   content, err := ioutil.ReadAll(redirectResult.Body)
   if err != nil {
	   fmt.Println("error", err)
	   return
   }
   fmt.Println(string(content))
}
