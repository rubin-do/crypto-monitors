package markets

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Ask struct {
	Price string `json:"price"`
}

type Depth struct {
	Asks []Ask `json:"asks"`
}

func MonitorGarantexPrice(prices chan float64) {
	for {
		resp, err := http.Get("https://garantex.io/api/v2/depth?market=usdtrub")

		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		var depth Depth

		json.Unmarshal(body, &depth)

		price, err := strconv.ParseFloat(depth.Asks[0].Price, 64)
		if err != nil {
			log.Fatal(err)
		}

		prices <- price

		resp.Body.Close()
	}
}
