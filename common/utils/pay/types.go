package pay

type Order struct {
	Body       string
	OutNo      string
	Total      float32
	OpenId     string
	NotifyPath string
}

type PayResponse struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type PayAmoount struct {
	Total int
}

type PayNotifyResponse struct {
	OutTradeNo    string
	TransactionId string
	Attach        string
	SuccessTime   string
	Amount        PayAmoount
}
