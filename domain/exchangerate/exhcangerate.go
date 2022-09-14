package exchangerate

import (
	"encoding/csv"
	"errors"
	"net/http"
	"time"
)

var ErrCurrencyNotExist = errors.New("currency not exist")

type exchangeRateMap map[string]rate

var exchangeRates exchangeRateMap

func init() {
	exchangeRates = make(map[string]rate)
	setExchangeRate()

	ticker := time.NewTicker(2 * time.Minute)
	go func() {
		for range ticker.C {
			setExchangeRate()
		}
	}()
}

type rate struct {
	SellRate string
	BuyRate  string
}

func setExchangeRate() error {
	resp, err := http.Get("https://rate.bot.com.tw/xrt/flcsv/0/day")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rows := csv.NewReader(resp.Body)
	rows.FieldsPerRecord = -1
	records, err := rows.ReadAll()
	if err != nil {
		return err
	}
	for i, record := range records {
		if i == 0 {
			continue
		}
		exchangeRates[record[0]] = rate{
			BuyRate:  record[2],
			SellRate: record[12],
		}
	}
	return nil
}

func GetRate(currency string) (string, string, error) {
	r, ok := exchangeRates[currency]
	if !ok {
		return "", "", ErrCurrencyNotExist
	}
	return r.BuyRate, r.SellRate, nil
}
