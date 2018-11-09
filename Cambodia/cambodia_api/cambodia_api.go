package cambodia_api

import (
"casino/util"
"encoding/json"
"fmt"
"gopkg.in/resty.v1"
)

type TokenInfo struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    *TokenData `json:"data"`
}

type TokenData struct {
	Access_token       string    `json:"access_token"`
	Expires            int       `json:"expires"`
	Refresh_expires_in string    `json:"refresh_expires_in"`
	User_info          *UserInfo `json:"user_info"`
}
type UserInfo struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}

func Authcodetotoken(host, appkey, authorization_code string) (string, string, error) {
	uri := "/v1/token/refresh/"
	requesturl := host + uri

	resp, err := resty.R().SetHeader("Authorization", appkey).SetQueryParam("authorization_code", authorization_code).Post(requesturl)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode() != 200 {
		return "", "", fmt.Errorf("Authcodetotoken Error: resp status code err: %v", resp.StatusCode())
	}

	var tokeninfo TokenInfo
	err = json.Unmarshal(resp.Body(), &tokeninfo)
	if err != nil {
		return "", "", err

	}
	return tokeninfo.Data.Access_token, tokeninfo.Data.User_info.Username, nil
}

////////Get Balance
type ErrorInfo struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type Balance struct {
	ErrorInfo
	Result *BalanceResult `json:"result"`
}
type BalanceResult struct {
	UserId       string  `json:"userId"`
	Balance      float64 `json:"balance"`
	CurrencyCode string  `json:"currencyCode"`
	ResponseDate string  `json:"responseDate"`
}

func GetBalance(host, signature, userId string) (float64, error) {
	parames := make(map[string]string)
	parames["signature"] = signature
	parames["cmd"] = "GetBalance"
	parames["userId"] = userId
	resp, err := resty.R().SetHeader("Content-Type", "application/json").SetBody(&parames).Post(host)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("GetBalance Error: resp status code err: %v", resp.StatusCode())
	}

	var b Balance
	err = json.Unmarshal(resp.Body(), &b)
	if err != nil {
		return 0, err
	}
	return b.Result.Balance, nil

}

//////////////////////Create Bet ///////////////
type CreateBetInfo struct {
	Signature      string  `json:"signature"`
	Cmd            string  `json:"cmd"`
	MemberId       string  `json:"memberId"`
	BetLogId       string  `json:"betLogId"`
	IsTrial        bool    `json:"isTrial"`
	Status         int     `json:"status"`
	PlayType       string  `json:"playType"`
	Result         string  `json:"result"`
	Note           string  `json:"note"`
	Currency       string  `json:"currency"`
	BetAmount      float64 `json:"betAmount"`
	ValidBetAmount float64 `json:"validBetAmount"`
	PayoutAmount   float64 `json:"payoutAmount"`
	WinLossAmount  float64 `json:"winLossAmount"`
	BetAt          string  `json:"betAt"`
	Odds           float64 `json:"odds"`
	IpAddress      string  `json:"IpAddress"`
}

type BetStatus struct {
	Ok   bool       `json:"ok"`
	Data *ErrorInfo `json:"data"`
}

func CreateBet(host, signature, memberId, betLogId, playType, betAt, currency string, status int, betAmount, validBetAmount, odds float64) error {
	var parames CreateBetInfo
	parames.Signature = signature
	parames.Cmd = "CreateBetLog"
	parames.MemberId = memberId
	parames.BetLogId = betLogId
	parames.Status = status
	parames.PlayType = playType
	parames.Currency = currency
	parames.BetAmount = betAmount
	parames.ValidBetAmount = validBetAmount
	parames.BetAt = betAt
	parames.Odds = odds

	fmt.Println(util.JsonFormat(&parames))
	resp, err := resty.R().SetHeader("Content-Type", "application/json").SetBody(&parames).Post(host)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("CreateBet Error: resp status code err: %v", resp.StatusCode())
	}

	var c BetStatus
	err = json.Unmarshal(resp.Body(), &c)
	if err != nil {
		return err
	}
	if !c.Ok {
		return fmt.Errorf("CreateBet Error: Code: %v, Msg: %v", c.Data.ErrorCode, c.Data.ErrorMessage)
	}
	return nil
}

//////////////////Update Bet////////////
type UpdateBetInfo struct {
	Signature string           `json:"signature"`
	Cmd       string           `json:"cmd"`
	Data      []*UpdateBetData `json:"data"`
}
type UpdateBetData struct {
	MemberId       string  `json:"memberId"`
	BetLogId       string  `json:"betLogId"`
	Status         int     `json:"status"`
	PlayType       string  `json:"playType"`
	Result         string  `json:"result"`
	Note           string  `json:"note"`
	Currency       string  `json:"Currency"`
	PayoutAmount   float64 `json:"payoutAmount"`
	WinLossAmount  float64 `json:"winLossAmount"`
	ValidBetAmount float64 `json:"validBetAmount"`
	Odds           float64 `json:"odds"`
	IpAddress      string  `json:"IpAddress"`
}

func UpdateBet(host, signature, memberId, betLogId, playType, currency string, status int, payoutAmount, winLossAmount float64) error {
	var parames UpdateBetInfo
	parames.Signature = signature
	parames.Cmd = "UpdateBetLog"

	var data UpdateBetData
	data.MemberId = memberId
	data.BetLogId = betLogId
	data.Status = status
	data.PlayType = playType
	data.Currency = currency
	data.PayoutAmount = payoutAmount
	data.WinLossAmount = winLossAmount
	parames.Data = []*UpdateBetData{&data}

	resp, err := resty.R().SetHeader("Content-Type", "application/json").SetBody(&parames).Post(host)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("UpdateBet Error: resp status code err: %v", resp.StatusCode())
	}

	var c BetStatus
	err = json.Unmarshal(resp.Body(), &c)
	if err != nil {
		return err
	}

	if !c.Ok {
		return fmt.Errorf("UpdateBet Error: Code: %v, Msg: %v", c.Data.ErrorCode, c.Data.ErrorMessage)
	}
	return nil
}

///////////////////Get Bet Log//////////////
type GetBetInfo struct {
	Signature string `json:"signature"`
	Cmd       string `json:"cmd"`
	BetLogId  string `json:"betLogId"`
	MemberId  string `json:"memberId"`
}
type GetBetStatus struct {
	Ok   bool             `json:"ok"`
	Data GetBetStatusData `json:"data"`
}
type GetBetStatusData struct {
	MemberId       string   `json:"memberId"`
	BetLogId       string   `json:"betLogId"`
	IsTrial        bool     `json:"isTrial"`
	Status         int      `json:"status"`
	PlayType       string   `json:"playType"`
	Result         []string `json:"result"`
	Note           string   `json:"note"`
	Currency       string   `json:"currency"`
	BetAmount      string   `json:"betAmount"`
	ValidBetAmount string   `json:"validBetAmount"`
	PayoutAmount   float64  `json:"payoutAmount"`
	WinLossAmount  string   `json:"winLossAmount"`
	BetAt          string   `json:"betAt"`
	Odds           string   `json:"odds"`
}

func GetBetLog(host, signature, betLogId, memberId string) (*GetBetStatus, error) {
	var parames GetBetInfo
	parames.Signature = signature
	parames.Cmd = "GetBetLog"
	parames.BetLogId = betLogId
	parames.MemberId = memberId

	resp, err := resty.R().SetHeader("Content-Type", "application/json").SetBody(&parames).Post(host)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("CreateBet Error: resp status code err: %v", resp.StatusCode())
	}

	var g GetBetStatus
	err = json.Unmarshal(resp.Body(), &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

