package nicehash

import (
	"github.com/mineselskabet/go-bitcoin"
)

type RigStats__ struct {
	TimeStamp     int64           `json:"statsTime"`
	Market        string          `json:"market"`
	UnpaidAmount  bitcoin.Bitcoin `json:"unpaidAmount"`
	Difficulty    float64         `json:"difficulty"`
	ProxyID       int             `json:"proxyId"`
	TimeConnected int64           `json:"timeConnected"`
	SpeedAccepted float64         `json:"speedAccepted"`
	Profitability float64         `json:"profitability"`
}

type Rig struct {
	RigID        string          `json:"rigId"`
	StatusTime   int64           `json:"statusTime"`
	UnpaidAmount bitcoin.Bitcoin `json:"unpaidAmount"`
}

func (c *Client) Rigs() ([]Rig, error) {
	URL := "https://api2.nicehash.com/main/api/v2/mining/rigs2"

	var proxy struct {
		MiningRigs []Rig `json:"miningRigs"`
	}

	err := c.exchangeJSON("GET", URL, nil, &proxy)
	if err != nil {
		return nil, err
	}

	return proxy.MiningRigs, nil
}

func (c *Client) MiningAddress() (string, error) {
	URL := "https://api2.nicehash.com/main/api/v2/mining/miningAddress"

	var proxy struct {
		Address string `json:"address"`
	}

	err := c.exchangeJSON("GET", URL, nil, &proxy)
	if err != nil {
		return "", err
	}

	return proxy.Address, nil
}

type AlgoStats struct {
	Unpaid        bitcoin.Bitcoin `json:"unpaid"`
	Profitability bitcoin.Bitcoin `json:"profitability"`
	SpeedAccepted float64         `json:"speedAccepted"`
	SpeedRejected float64         `json:"speedRejected"`
	DisplaySuffix string          `json:"displaySuffix"`
	Active        bool            `json:"isActive"`
}

func (c *Client) AlgoStats() (map[string]AlgoStats, error) {
	URL := "https://api2.nicehash.com/main/api/v2/mining/algo/stats"
	var proxy struct {
		Algorithms map[string]AlgoStats `json:"algorithms"`
	}

	err := c.exchangeJSON("GET", URL, nil, &proxy)
	if err != nil {
		return nil, err
	}

	return proxy.Algorithms, nil
}
