package main

import (
	"gopkg.in/resty.v1"
	"fmt"
	"encoding/json"
	"sync"
)

type AuthCodeInfo struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data DataInfo `json:"data"`
}

type DataInfo struct {
	Authorization_code string `json:"authorization_code"`
	Game_authorization interface{} `json:"game_authorization"`
}

func get_auth_code() string{
	url:= "http://dev-dabao-api.win.town/testonly?agent=10002&username=ge1000004433"
	resp,err := resty.R().Get(url)
	if err != nil {
		fmt.Println("get auth code err :", err)
	}
	fmt.Println(resp.String())
	var authinfo AuthCodeInfo
	//authinfo := AuthCodeInfo{}
	err = json.Unmarshal(resp.Body(),&authinfo)
	if err != nil {
		fmt.Println("get ayth code json unmarshal err :",err)

	}
	fmt.Println(authinfo.Data.Authorization_code)
	return authinfo.Data.Authorization_code
}

type TokenInfo struct {
	Token string `json:"token"`
	Code int `json:"code"`
	Userid string `json:"userid"`
}

func auth(authcode ,machineid,queryurl string) (string,string) {
	m := make(map[string]string)
	m["authorization_code"] = authcode
	m["enc"] = "false"
	host := fmt.Sprintf("%s/%s/auth",queryurl,machineid)
	resp,err := resty.R().SetQueryParams(m).Get(host)
	if err != nil {
		fmt.Println("auth err :",err)
	}
	var t *TokenInfo
	fmt.Println(resp.String())
	err = json.Unmarshal(resp.Body(),&t)
	if err != nil {
		fmt.Println("auth json unmarshal err :" ,err)
	}

	fmt.Println(t)
	return t.Token, t.Userid
}

func enter(token,userid ,machineid,queryurl string) {
	m := make(map[string]string)
	m["token"]=token
	m["userid"]=userid
	m["enc"]="false"
	m["currency"]="cny"
	m["machine_id"] = machineid
	host := fmt.Sprintf("%s/%s/enter",queryurl,machineid)
	resp,err := resty.R().SetQueryParams(m).Get(host)
	if err != nil {
		fmt.Println("enter resp err :",err)
	}
	fmt.Println(resp.String())
}

type Parmes struct {
	Token string `json:"token"`
	Userid string	`json:"userid"`
	Enc string	`json:"enc"`
	Currency string	`json:"currency"`
	Machine_id string	`json:"machine_id"`
	Bet_index string	`json:"bet_index"`
}
func play(token,userid,machineid,queryurl string)  {
	host := fmt.Sprintf("%s/%s/play",queryurl,machineid)
	m := make(map[string]string)
	m["token"] = token
	m["userid"] = userid
	m["enc"] = "false"
	m["currency"] = "cny"
	m["machine_id"] = machineid
	m["bet_index"] = "0"
	//resp,err := resty.R().SetHeader("Content-Type","application/x-www-form-urlencoded").SetBody(p).Post(host)
	resp,err := resty.R().SetQueryParams(m).Post(host)
	if err != nil {
		fmt.Println("play err:",err)
	}
	fmt.Println(resp.StatusCode(),resp.String())

}

func main()  {
	var wg sync.WaitGroup
	for j:= 0 ;j<1;j++{
		wg.Add(1)
		go func (){
			defer wg.Done()
			queryurl := "http://192.168.2.237:10005"
			machineid := "Fafafa"
			auth_code := get_auth_code()
			token,userid := auth(auth_code,machineid,queryurl)
			fmt.Println(token,userid)
			enter(token,userid,machineid,queryurl)
			for i:=0; i<1;i++ {
				//fmt.Println(i)
				play(token,userid,machineid,queryurl)
			}

		}()
	}

	wg.Wait()
}