package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/idprm/go-yellowclinic/src/config"
)

/**
 * https://api.golkarekta.com/v7/callback/data/cepat_sehat
 */
type CallbackRequest struct {
	Voucher string `form:"voucher" json:"voucher"`
}

func CallbackVoucher(voucher string) (string, error) {
	urlAddress := config.ViperEnv("GOLKAREKTA_URL") + "/v7/callback/data/cepat_sehat/"

	reqBody := CallbackRequest{
		Voucher: voucher,
	}

	var bearer = "Bearer " + "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImV4cCI6MTY5MjE0MTg5OH0.eyJuYW1lIjoidmlzaW9ucGx1cyIsIndoaXRlbGlzdF9pcCI6Ijo6MSIsInR5cGUiOiIiLCJ1c2VybmFtZSI6InZpc2lvbnBsdXMiLCJrZXkiOiJ2aXNpb25wbHVzIiwidHlwZV9leHBpcmVkIjoiY29udGludWUiLCJpc19zdGFydCI6IjIwMjItMDgtMTUgMDg6MDA6MDAiLCJpc19lbmQiOiIyMDIzLTA4LTE1IDIzOjU5OjU5Iiwib25fZGF5cyI6IiIsInRpbWVzIjoiMjAyMjA4MTUwMjI0NTkifQ.4rp_CRiw8axNnMI9gWLU5We12gBcJX3csQ5aBsenSoM"

	param := url.Values{}
	param.Set("voucher_code", reqBody.Voucher)
	payload := bytes.NewBufferString(param.Encode())

	req, err := http.NewRequest("POST", urlAddress, payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", bearer)

	if err != nil {
		return "", errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	resp, err := client.Do(req)

	if err != nil {
		return "", errors.New(err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return string([]byte(body)), nil
}
