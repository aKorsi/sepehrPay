<div dir="rtl">

# کتابخانه اتصال به سامانه [ پرداخت الکترونیک سپهر](https://www.sepehrpay.com/) 

روش کار:

<div dir="ltr">

```go

sepehrPayClient := sepehrPay.Client{TerminalId: TerminalId}

// payment request 
token, err := sepehrPayClient.GetTokenREST(orderId, amount, callBackUrl, additionalData)
if err != nil {
	return err
}
code := sepehrPayClient.MakeForm(token, "")
	
// verify request
resp, err := sepehrPayClient.ParseCallBack(req)
if err != nil {
	return err
}
_, err = sepehrPayClient.VerifyREST(resp.DigitalReceipt)
if err != nil {
	return err
}

```

</div>

در صورت خطا یا مشکلی حتما اعلام بفرمایید

<div dir="ltr">

`github.com/tiaguinho/gosoap`

**{SOAP package for Go}**

`github.com/go-resty/resty`

**{Simple HTTP and REST client library for Go}**

`github.com/json-iterator/go`

**{A high-performance 100% compatible drop-in replacement of "encoding/json"}**

</div>

</div>

