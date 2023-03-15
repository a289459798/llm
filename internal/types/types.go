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
	OpenId string `json:"openId"`
}

type Task struct {
	Title          string `json:"title"`
	Status         bool   `json:"status"`
	Total          int    `json:"total"`
	CompleteNumber int    `json:"completeNumber"`
	Type           string `json:"type"`
	Amount         int    `json:"amount"`
	Max            int    `json:"max"`
}

type TaskResponse struct {
	Content string `json:"content"`
	List    []Task `json:"list"`
	Max     int    `json:"max"`
	Have    int    `json:"have"`
	Tips    string `json:"tips"`
}

type TaskCompleteResponse struct {
	Total   uint32 `json:"total"`
	Amount  uint32 `json:"amount"`
	Welfare uint32 `json:"welfare"`
}

type TaskRequest struct {
	Type string `json:"type,options=share|ad"`
}

type TaskShareFollowRequest struct {
	OpenId string `json:"openId"`
}

type AIInfoResponse struct {
	Name    string `json:"name"`
	Photo   string `json:"photo"`
	Call    string `json:"call"`
	Welcome string `json:"welcome"`
	Status  bool   `json:"status"`
	RoleId  uint32 `json:"roleId"`
}

type AIEditRequest struct {
	Name   string `form:"name"`
	Call   string `form:"call"`
	RoleId uint32 `form:"roleId"`
	Status bool   `form:"status"`
}

type ReportRequest struct {
	Content string `json:"content"`
}

type WorkRequest struct {
	Use       string `json:"use"`
	Introduce string `json:"introduce"`
	Content   string `json:"content"`
}

type ReportResponse struct {
	Data string `json:"data"`
}

type ImageRequest struct {
	Content string `json:"content"`
}

type ImageResponse struct {
	Url  string `json:"url"`
	Task string `json:"task"`
}

type Pic2picRequest struct {
	Style string `form:"style"`
}

type Pic2picTaskRequest struct {
	Task string `form:"task"`
}

type WatermarkRequest struct {
	Url      string  `json:"url"`
	Content  string  `json:"content"`
	Position string  `json:"position"`
	Opacity  float32 `json:"opacity"`
	FontSize uint    `json:"fontSize"`
	Color    string  `json:"color"`
}

type ImageEditRequest struct {
	Content string `form:"content"`
}

type QiMingRequest struct {
	First    string `json:"first"`
	Number   string `json:"number"`
	Birthday string `json:"birthday"`
	Sex      string `json:"sex"`
	Fix      string `json:"fix,optional"`
	Other    string `json:"other,optional"`
}

type JieMengRequest struct {
	Content string `json:"content"`
}

type SuanMingRequest struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
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

type HoroscopeRequest struct {
	Birthday      string `json:"birthday,optional"`
	Constellation string `json:"constellation,optional"`
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

type ExamRequest struct {
	Content string `json:"content"`
	Type    string `json:"type"`
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

type PlaygroundRequest struct {
	Content string `json:"content"`
	Lang    string `json:"lang,optional"`
}

type CodeResponse struct {
	Data string `json:"data"`
}

type ActivityRequest struct {
	Way     string `json:"way"`
	Target  string `json:"target"`
	Date    string `json:"date"`
	User    string `json:"user"`
	Content string `json:"content"`
}

type DiaryRequest struct {
	City    string `json:"city"`
	Weather string `json:"weather"`
	Content string `json:"content"`
}

type ArticleRequest struct {
	Subject string `json:"subject"`
	Type    string `json:"type"`
	Number  string `json:"number"`
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

type PursueRequest struct {
	Content string `json:"content"`
}

type ChatResponse struct {
	Data string `json:"data"`
}

type ChatRequest struct {
	ChatId     string `json:"chatId"`
	TemplateId uint32 `json:"templateId"`
	Message    string `json:"message"`
}

type ChatHistoryResponse struct {
	ChatId  string        `json:"chatId"`
	History []ChatHistory `json:"history"`
}

type ChatHistory struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type ChatTemplateResponse struct {
	List []ChatTemplate `json:"list"`
}

type ChatTemplate struct {
	TemplateId uint32 `json:"templateId"`
	Message    string `json:"message"`
}

type ChatTemplateRequest struct {
	Type string `form:"type,optional"`
}

type TranslateRequest struct {
	Content string `json:"content"`
	Lang    string `json:"lang"`
}

type ConvertResponse struct {
	Data string `json:"data"`
}

type ValidRequest struct {
	Content string `json:"content"`
}

type ValidResponse struct {
	Data    string `json:"data"`
	ShowAd  bool   `json:"showAd"`
	Consume int    `json:"consume"`
}

type MessageResponse struct {
	Id      uint32 `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
}

type WeChatCallbackResponse struct {
}
