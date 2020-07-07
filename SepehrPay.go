package sepehrPay

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/tiaguinho/gosoap"
	"strconv"
	"strings"
)

const (
	tokenRESTService  = "https://mabna.shaparak.ir:8081/V1/PeymentApi/GetToken"
	verifyRESTService = "https://mabna.shaparak.ir:8081/V1/PeymentApi/Advice"
	tokenService      = "https://mabna.shaparak.ir:8082/Token.svc?wsdl"
	verifyService     = "https://mabna.shaparak.ir:8082/ipg.svc?wsdl"
	gatewayURL        = "https://mabna.shaparak.ir:8080/Pay"
)

type Client struct {
	TerminalId int64
}

type getTokenRequest struct {
	Amount      int64  `json:"Amount"`
	CallbackUrl string `json:"CallbackUrl"`
	InvoiceId   string `json:"InvoiceId"`
	Payload     string `json:"Payload"`
	TerminalId  int64  `json:"TerminalId"`
}

type verifyRequest struct {
	DigitalReceipt string `json:"digitalreceipt"`
	TID            int64  `json:"Tid"`
}

type getTokenResponse struct {
	Status      int    `json:"Status",xml:"Status"`
	AccessToken string `json:"Accesstoken,omitempty",xml:"AccessToken"`
}

type verifyResponse struct {
	Status   string `json:"Status,omitempty",xml:"Status"`
	ReturnId string `json:"ReturnId,omitempty",xml:"ReturnId"`
	Message  string `json:"Message,omitempty",xml:"Message"`
}

type ParseResponse struct {
	RespCode       int
	RespMsg        string
	Amount         int64
	InvoiceId      string
	Payload        string
	TerminalId     int64
	TraceNumber    int64
	RRN            int64
	DatePaid       string
	DigitalReceipt string
	IssuerBank     string
	CardNumber     string
}

func (c *Client) getError(err int) error {
	switch err {
	case -1:
		return errors.New("تراکنش پیدا نشد")
	case -2:
		return errors.New("تراکنش قبلا Reverse شده است")
	case -3:
		return errors.New("خطای عمومی")
	case -4:
		return errors.New("امکان انجام درخواست برای این تراکنش وجود ندارد")
	case -5:
		return errors.New("آدرس IP نامعتبر میباشد (IP در لیست آدرسهای معرفی شده توسط پذیرنده موجود نمیباشد)")
	case -6:
		return errors.New("عدم فعال بودن سرویس برگشت تراکنش برای پذیرنده")
	default:
		return nil
	}
}

func (c *Client) GetError(respCode int) error {
	switch respCode {
	case -1:
		return errors.New("کاربر دکمه انصراف را در صفحه پرداخت فشرده است")
	case -2:
		return errors.New("زمان انجام تراکنش برای کاربر به اتمام رسیده است")
	default:
		return nil
	}
}

func (c *Client) GetToken(invoiceId string, amount int64, callBackURL, payLoad string) (string, error) {
	soap, err := gosoap.SoapClient(tokenService)
	if err != nil {
		return "", err
	}

	err = soap.Call("GetToken", gosoap.Params{
		"amount":      amount,
		"callbackURL": callBackURL,
		"invoiceID":   invoiceId,
		"terminalID":  c.TerminalId,
		"payload":     payLoad,
	})
	if err != nil {
		return "", err
	}

	res := new(getTokenResponse)
	err = soap.Unmarshal(res)
	if err != nil {
		return "", err
	}

	err = c.getError(res.Status)
	if err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

func (c *Client) GetTokenREST(invoiceId string, amount int64, callBackURL, payLoad string) (string, error) {
	client := resty.New()
	jsObject := jsoniter.ConfigCompatibleWithStandardLibrary

	reqBody, err := jsoniter.Marshal(&getTokenRequest{
		Amount:      amount,
		CallbackUrl: callBackURL,
		InvoiceId:   invoiceId,
		Payload:     payLoad,
		TerminalId:  c.TerminalId,
	})
	if err != nil {
		return "", err
	}

	res, err := client.R().
		SetBody(string(reqBody)).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Post(tokenRESTService + "/")
	if err != nil {
		return "", err
	}

	result := getTokenResponse{}
	err = jsObject.Unmarshal(res.Body(), &result)
	if err != nil {
		return "", err
	}

	err = c.getError(result.Status)
	if err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (c *Client) Verify(digitalReceipt string) (string, error) {
	soap, err := gosoap.SoapClient(verifyService)
	if err != nil {
		return "", err
	}

	err = soap.Call("Advice", gosoap.Params{
		"digitalreceipt": digitalReceipt,
		"Tid":            c.TerminalId,
	})
	if err != nil {
		return "", err
	}

	res := new(verifyResponse)
	err = soap.Unmarshal(res)
	if err != nil {
		return "", err
	}

	if res.Status == "NOk" {
		return "", errors.New(res.Message)
	}

	return res.ReturnId, nil
}

func (c *Client) VerifyREST(digitalReceipt string) (string, error) {
	client := resty.New()
	jsObject := jsoniter.ConfigCompatibleWithStandardLibrary

	reqBody, err := jsoniter.Marshal(&verifyRequest{
		DigitalReceipt: digitalReceipt,
		TID:            c.TerminalId,
	})
	if err != nil {
		return "", err
	}

	res, err := client.R().
		SetBody(string(reqBody)).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Post(verifyRESTService + "/")
	if err != nil {
		return "", err
	}

	result := verifyResponse{}
	err = jsObject.Unmarshal(res.Body(), &result)
	if err != nil {
		return "", err
	}

	if strings.ToLower(result.Status) == "nok" {
		return "", errors.New(result.Message)
	}

	if strings.ToLower(result.Status) == "ok" {
		return result.ReturnId, nil
	}

	return "", errors.New(string(res.Body()))
}

func (c *Client) MakeForm(token, body string) string {
	return fmt.Sprintf(
		`<form id="payForm" name="payForm" action="%s" method="post" target="_self">
					<input type="hidden" id="TerminalID" name="TerminalID" value='%s' />
					<input type="hidden" id="token" name="token" value='%s' />
					<div>%s</div>
					<input type="submit" value="پرداخت" />
				</form>`,
		gatewayURL,
		strconv.FormatInt(c.TerminalId, 10),
		token,
		body,
	)
}

type Peeker interface {
	Peek(key string) []byte
}

func (c *Client) ParseCallBack(req Peeker) (ParseResponse, error) {
	resp := ParseResponse{}
	var err error
	resp.RespCode, err = strconv.Atoi(string(req.Peek("respcode")))
	if err != nil {
		return resp, err
	}

	resp.RespMsg = string(req.Peek("respmsg"))

	resp.Amount, err = strconv.ParseInt(string(req.Peek("amount")), 10, 64)
	if err != nil {
		return resp, err
	}

	resp.InvoiceId = string(req.Peek("invoiceid"))

	resp.Payload = string(req.Peek("payload"))

	resp.TerminalId, err = strconv.ParseInt(string(req.Peek("terminalid")), 10, 64)
	if err != nil {
		return resp, err
	}

	resp.TraceNumber, err = strconv.ParseInt(string(req.Peek("tracenumber")), 10, 64)
	if err != nil {
		return resp, err
	}

	resp.RRN, err = strconv.ParseInt(string(req.Peek("rrn")), 10, 64)
	if err != nil {
		return resp, err
	}

	resp.DatePaid = string(req.Peek("datepaid"))

	resp.DigitalReceipt = string(req.Peek("digitalreceipt"))

	resp.IssuerBank = string(req.Peek("issuerbank"))

	resp.CardNumber = string(req.Peek("cardnumber"))

	return resp, nil
}
