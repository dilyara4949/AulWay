package stripe

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SimulatePayment(amount int, currency string) (bool, error) {
	reqBody := map[string]interface{}{
		"amount":   amount,
		"currency": currency,
		"status":   "succeeded", // "failed"
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:12111/v1/charges", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}
