package main

import (
	"fmt"
	"gopkg.in/resty.v1"
	"encoding/json"
)

type TokenInfo struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Data TokenData `json:"data"`
}

type TokenData struct {
	Access_token string `json:"access_token"`
	Expires string `json:"expires"`
	Refresh_expires_in string `json:"refresh_expires_in"`
	User_info UserInfo `json:"user_info"`

}
type UserInfo struct {
	Uid string `json:"uid"`
	UserInfo

}



func authcodetotoken(authorization_code string) (string,error) {
	host := " http://dev-dabao-api.win.town"
	appkey := "13940263647092369a024cbf287300fc17161d527277be9b0256e44a"
	addr := "/v1/token/refresh"

	requesturl := host + addr
	params := fmt.Sprintf("authorization_code=%v", authorization_code)
	resp, err := resty.R().SetHeader("Authorization", appkey).SetQueryString(params).Post(requesturl)
	if err != nil {
		fmt.Println(err)
		return "",err
	}
	err := json.Unmarshal()
	resp.Body()
}

func main() {
	host := " http://dev-dabao-api.win.town"
	appkey := "13940263647092369a024cbf287300fc17161d527277be9b0256e44a"
	addr := "/v1/token/refresh"

	requesturl := host + addr
	params := fmt.Sprintf("authorization_code=%v", authorization_code)
	resp, err := resty.R().SetHeader("Authorization", appkey).SetQueryString(params).Post(requesturl)
	//r,err :=resty.R().SetQueryString(params).Post(requesturl)
}
