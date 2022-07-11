package markets

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Ask struct {
	Price  string `json:"price"`
	Volume string
}

type Depth struct {
	Asks []Ask `json:"asks"`
	Bids []Ask
}

func MonitorGarantexPrice(orders chan<- Order) {
	for {
		resp, err := http.Get("https://garantex.io/api/v2/depth?market=usdtrub")

		if err != nil {
			log.Println(err)
			time.Sleep(60 * time.Second)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		var depth Depth

		json.Unmarshal(body, &depth)
		resp.Body.Close()

		// handle garantex strange behavior
		if len(depth.Asks) == 0 {
			continue
		}

		buyPrice, err := strconv.ParseFloat(depth.Asks[0].Price, 64)
		if err != nil {
			log.Fatal(err)
		}

		sellPrice, err := strconv.ParseFloat(depth.Bids[0].Price, 64)
		if err != nil {
			log.Fatal(err)
		}

		bestOrder := Order{
			"Null",
			"Garantex",
			"Market",
			buyPrice,
			sellPrice,
			depth.Bids[0].Volume,
			"-",
			"-",
			"https://garantex.io/trading/usdtrub",
		}

		orders <- bestOrder

		time.Sleep(2 * time.Second)
	}
}
