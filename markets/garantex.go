package markets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
			fmt.Println(err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		if err != nil {
			fmt.Println(err)
			return
		}

		var depth Depth

		json.Unmarshal(body, &depth)

		price, err := strconv.ParseFloat(depth.Asks[0].Price, 64)
		prices <- price
	}
	close(prices)
}