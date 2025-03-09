package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Location struct {
	Rs       int    `json:"rs"`       // 响应状态，1表示成功
	Code     int    `json:"code"`     // 状态码，0表示无错误
	Address  string `json:"address"`  // 地理地址
	IP       string `json:"ip"`       // 查询的 IP 地址
	IsDomain int    `json:"isDomain"` // 是否为域名，0表示不是，1表示是
}

func GetLocation(ip string) (*Location, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.ip.cn/api/index?ip=%s&type=1", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var location Location
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil, err
	}

	return &location, nil
}
