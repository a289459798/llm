package pay

type Order struct {
	Body   string
	OutNo  string
	Total  float32
	OpenId string
}

type PayResponse struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}
