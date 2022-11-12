package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/idprm/go-yellowclinic/src/config"
)

/**
 * https://api.golkarekta.com/v7/callback/data/cepat_sehat
 */
type CallbackRequest struct {
	VoucherCode string `json:"voucher_code"`
}

func CallbackVoucher(voucher string) (string, error) {
	url := config.ViperEnv("GOLKAREKTA_URL") + "/v7/callback/data/cepat_sehat/"

	reqBody := CallbackRequest{
		VoucherCode: voucher,
	}

	payload, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=utf8")

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
