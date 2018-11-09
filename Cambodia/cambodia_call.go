package main

import (
	"fmt"
	"gopkg.in/resty.v1"
	"encoding/json"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"

	"casino/cambodia"
	"time"
)

type TokenInfo struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data TokenData `json:"data"`
}

type TokenData struct {
	Access_token string `json:"access_token"`
	Expires int `json:"expires"`
	Refresh_expires_in string `json:"refresh_expires_in"`
	User_info UserInfo `json:"user_info"`

}
type UserInfo struct {
	Uid int `json:"uid"`
	Username string `json:"username"`
	Balance int `json:"balance"`
}



func authcodetotoken(host,appkey,authorization_code string) (string,string,error) {
	addr := "/v1/token/refresh/"
	requesturl := host + addr

	params := make(map[string]string)
	params["authorization_code"]=authorization_code
	resp, err := resty.R().SetHeader("Authorization", appkey).SetQueryParam("authorization_code",authorization_code).Post(requesturl)
	//headers := make(map[string]string)
	//headers["Content-Type"]="application/json"
	//headers["Authorization"]=appkey
	//resp, err := resty.R().SetHeaders(headers).SetBody(params).Post(requesturl)
	//fmt.Println(resp.String())
	if err != nil {
		fmt.Println("resp err: ",err)
		return "","",err
	}
	var tokeninfo TokenInfo
	err = json.Unmarshal(resp.Body(),&tokeninfo)
	if err != nil {
		fmt.Println("json unmarshal err: ",err)
		return "","",err

	}
	//fmt.Println(tokeninfo)
	return tokeninfo.Data.Access_token,tokeninfo.Data.User_info.Username,nil

}

type TokenAndUid struct {
	Token string
	Uid string
}

func Auth(c *gin.Context) {
	host := "http://test-dev-dabao-api.win.town"
	appkey := "13940263647092369a024cbf287300fc17161d527277be9b0256e44a"
	authorization_code := c.Query("authorization_code")
	//token,uid,err :=authcodetotoken(host,appkey,authorization_code)
	token,uid,err :=cambodia.Authcodetotoken(host,appkey,authorization_code)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(uid)
	var tokenanduid TokenAndUid
	tokenanduid.Token= token
	tokenanduid.Uid=uid
	///////////md5加密uid//////////
	/*data := []byte(strconv.Itoa(uid))
	has := md5.New()
	has.Write(data)
	tokenanduid.Uid = hex.EncodeToString(has.Sum(nil))*/
	c.JSON(http.StatusOK,&tokenanduid)

}

///////////////////Get Balance///////////////
type Balance struct {
	ErrorCode int `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Result BalanceResult
}
type BalanceResult struct {
	UserId string `json:"userId"`
	Balance float64 `json:"balance"`
	CurrencyCode string `json:"currencyCode"`
	ResponseDate string `json:"responseDate"`
}
func GetBalance(host ,signature ,userId string) (float64,error) {
	parames := make(map[string]string)
	parames["signature"]=signature
	parames["cmd"]="GetBalance"
	parames["userId"]=userId
	resp,err :=resty.R().SetHeader("Content-Type","application/json").SetBody(&parames).Post(host)
	if err != nil {
		fmt.Println(err)
		return 0,err
	}
	resp.Body()
	var b Balance
	err = json.Unmarshal(resp.Body(),&b)
	if err != nil {
		return 0,err
	}
	return b.Result.Balance,nil

}
////////////////////////////////////////
type Balance2 struct {
	Balance float64 `json:"balance"`
}

func Blance(c *gin.Context)  {
	host :=  "http://sw-test.lx345.net/ks/apis"
	userId := "ge1000004433"
	signature := "CChWk6zMzO1t9DW3RbRoM2mWUKmV5hJP"
	//b,err := GetBalance(host,signature,userId)
	b,err := cambodia.GetBalance(host,signature,userId)
	if err != nil {
		fmt.Println("Blance",err)
	}
	var bb Balance2
	bb.Balance = b
	//fmt.Println(b)
	c.JSON(http.StatusOK,&b)
}
//////////////////////Create Bet ///////////////
type CreateBetInfo struct {
	Signature string
	Cmd string
	MemberId string
	BetLogId string
	IsTrial bool
	Status int
	PlayType string
	Result string
	Note string
	Currency string
	BetAmount float64
	ValidBetAmount float64
	PayoutAmount float64
	WinLossAmount float64
	BetAt string
	Odds float64
	IpAddress string
}
type BetStatus struct {
	Ok bool `json:"ok"`
	Data map[string]interface{}
}
func CreateBet(host,signature,memberId,BetLogId ,playType , betAt string ,status int,betAmount ,validBetAmount,odds float64) BetStatus {
	var parames CreateBetInfo
	parames.Signature=signature
	parames.Cmd="CreateBetLog"
	parames.MemberId=memberId
	parames.BetLogId=BetLogId
	parames.Status=status
	parames.PlayType=playType
	parames.BetAmount=betAmount
	parames.ValidBetAmount=validBetAmount
	parames.BetAt = betAt
	parames.Odds = odds
	resp,err := resty.R().SetHeader("Content-Type","application/json").SetBody(parames).Post(host)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(resp.String())
	var c BetStatus
	err = json.Unmarshal(resp.Body(),&c)
	if err != nil {
		fmt.Println(err)
	}
	return c
}

///////////////////////////////////////




//////////////////Update Bet////////////
type UpdateBetInfo struct {
	Signature string
	Cmd string
	Data []UpdateBetData
}
type UpdateBetData struct {
	MemberId string
	BetLogId string
	Status int
	PlayType string
	Result string
	Note string
	Currency string
	PayoutAmount float64
	WinLossAmount float64
	ValidBetAmount float64
	Odds float64
	IpAddress string
}

func UpdateBet(host,signature,memberId,betLogId,playType,Currency string,status int ,payoutAmount, winLossAmount float64) BetStatus {
	var parames UpdateBetInfo
	parames.Signature =signature
	parames.Cmd="UpdateBetLog"

	var data UpdateBetData
	data.MemberId=memberId
	data.BetLogId=betLogId
	data.Status=status
	data.PlayType=playType
	data.Currency=Currency
	data.PayoutAmount=payoutAmount
	data.WinLossAmount=winLossAmount
	parames.Data=[]UpdateBetData{data}
	resp,err := resty.R().SetHeader("Content-Type","application/json").SetBody(parames).Post(host)
	if err != nil {
		fmt.Println(err)
	}
	var c BetStatus
	err = json.Unmarshal(resp.Body(),&c)
	if err != nil {
		fmt.Println(err)
	}
	return c
}

///////////////////Get Bet Log//////////////
type GetBetInfo struct {
	Signature string
	Cmd string
	BetLogId string
	MemberId string
}
type GetBetStatus struct {
	Ok bool `json:"ok"`
	Data GetBetStatusData `json:"data"`
}
type GetBetStatusData struct {
	MemberId string `json:"memberId"`
	BetLogId string `json:"betLogId"`
	IsTrial bool `json:"isTrial"`
	Status int `json:"status"`
	PlayType string `json:"playType"`
	Result []string `json:"result"`
	Note string `json:"note"`
	Currency string `json:"currency"`
	BetAmount string `json:"betAmount"`
	ValidBetAmount string `json:"validBetAmount"`
	PayoutAmount float64 `json:"payoutAmount"`
	WinLossAmount string `json:"winLossAmount"`
	BetAt string `json:"betAt"`
	Odds string `json:"odds"`
}

func GetBetLog(host,signature,betLogId,memberId string) GetBetStatus {
	var parames GetBetInfo
	parames.Signature=signature
	parames.Cmd="GetBetLog"
	parames.BetLogId=betLogId
	parames.MemberId=memberId

	resp,err := resty.R().SetHeader("Content-Type","application/json").SetBody(parames).Post(host)
	if err != nil {
		fmt.Println("GetBetLog err: ",err)
	}
	var g GetBetStatus
	err = json.Unmarshal(resp.Body(),&g)
	if err != nil {
		fmt.Println("GetBetLog json unmarshal err: ",err)
	}
	return g
}

func main() {
//GetBalance("http://sw-test.lx345.net/ks/apis")
	/*host:= "http://sw-test.lx345.net/ks/apis"
	signature:="CChWk6zMzO1t9DW3RbRoM2mWUKmV5hJP"
	memberId:="ge1000004433"
	BetLogId:="2018092816"
	playType:="10001"
	status := 0
	betAt := "2018-09-29 16:11:00"
	var betAmount float64 = 100
	var validBetAmount float64= 100
	var odds float64 = 1

	err :=cambodia.CreateBet(host,signature,memberId,BetLogId,playType,betAt,"CNY",status,betAmount,validBetAmount,odds)
	if err != nil {
		fmt.Println(err)
	}
	status=1
	err1 := cambodia.UpdateBet(host,signature,memberId,BetLogId,playType,"CNY",status,100,100)
	if err1 != nil {
		fmt.Println(err1)
	}
	g,err:= cambodia.GetBetLog(host,signature,BetLogId,memberId)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(g)
	*/
	a := time.Now().String()
	fmt.Printf("%v,%s",a,a)

	/*
	router := gin.Default()
	router.GET("/Auth",Auth)
	router.GET("/Blance",Blance)
	router.Run("0.0.0.0:12345")
*/
}
