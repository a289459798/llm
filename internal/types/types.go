// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	Code string `json:"code"`
}

type InfoRequest struct {
}

type InfoResponse struct {
	Token  string `json:"token"`
	Amount uint32 `json:"amount"`
	Uid    uint32 `json:"uid"`
}

type ReportRequest struct {
	Content string `json:"content"`
}

type ReportResponse struct {
	Data string `json:"data"`
}

type ImageRequest struct {
	Content string `json:"content"`
}

type ImageResponse struct {
	Url string `json:"url"`
}

type WatermarkRequest struct {
	Url      string  `json:"url"`
	Content  string  `json:"content"`
	Position string  `json:"position"`
	Opacity  float32 `json:"opacity"`
	FontSize uint    `json:"fontSize"`
	Color    string  `json:"color"`
}

type QiMingRequest struct {
	First    string `json:"first"`
	Number   uint16 `json:"number"`
	Brithday string `json:"brithday"`
	Sex      string `json:"sex"`
	Fix      string `json:"fix,optional"`
	Other    string `json:"other,optional"`
}

type JieMengRequest struct {
	Content string `json:"content"`
}

type SuanMingRequest struct {
	Name     string `json:"name"`
	Brithday string `json:"brithday"`
	Sex      string `json:"sex,optional"`
	Content  string `json:"content,optional"`
}

type GSQiMingRequest struct {
	Industry string `json:"industry"`
	Range    string `json:"range"`
	Culture  string `json:"culture,optional"`
	Other    string `json:"other,optional"`
}

type YYQiMingRequest struct {
	Name  string `json:"name"`
	Sex   string `json:"sex"`
	Other string `json:"other,optional"`
}

type DivinationResponse struct {
	Data string `json:"data"`
}

type IdiomRequest struct {
	Pre     string `json:"pre,optional"`
	Content string `json:"content"`
}

type TwentyFourRequest struct {
	Content string `json:"content"`
}

type GameResponse struct {
	Code uint   `json:"code"`
	Data string `json:"data"`
}

type RegularRequest struct {
	Content string `json:"content"`
}

type GenerateRequest struct {
	Content string `json:"content"`
	Lang    string `json:"lang,optional"`
	ChatId  string `json:"chat_id,optional"`
}

type NameRequest struct {
	Content string `json:"content"`
	Lang    string `json:"lang"`
	Type    string `json:"type"`
}

type CodeResponse struct {
	Data string `json:"data"`
}

type ActivityRequest struct {
	Way     string `json:"way"`
	Target  string `json:"target"`
	Period  string `json:"period"`
	User    string `json:"user"`
	Content string `json:"content"`
}

type CreationResponse struct {
	Data string `json:"data"`
}

type IntroduceRequest struct {
	Name     string `json:"name"`
	Native   string `json:"native"`
	Interest string `json:"interest"`
	Way      string `json:"way,optional"`
	Content  string `json:"content,optional"`
}

type SalaryRequest struct {
	Content string `json:"content"`
}

type RejectRequest struct {
	Type    string `json:"type"`
	Way     string `json:"way"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Data string `json:"data"`
}

type TranslateRequest struct {
	Content string `json:"content"`
	Lang    string `json:"lang"`
}

type ConvertResponse struct {
	Data string `json:"data"`
}
