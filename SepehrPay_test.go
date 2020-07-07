package sepehrPay

import (
	"testing"
)

func TestGetToken(t *testing.T) {
	client := Client{
		TerminalId: 0,
	}

	res, err := client.GetToken("123123aa", 1_000_000, "http://localhost/order", "none")
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestGetTokenREST(t *testing.T) {
	client := Client{
		TerminalId: 0,
	}

	res, err := client.GetTokenREST("123123124aa", 1_000_000, "http://localhost/order", "none")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestVerify(t *testing.T) {
	client := Client{
		TerminalId: 0,
	}

	res, err := client.Verify("123")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestVerifyREST(t *testing.T) {
	client := Client{
		TerminalId: 0,
	}

	res, err := client.VerifyREST("123")
	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}
