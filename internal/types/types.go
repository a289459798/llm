// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	Code    string `json:"code"`
	Channel string `json:"channel,optional"`
}

type InfoRequest struct {
}

type InfoResponse struct {
	Token     string `json:"token"`
	Amount    uint32 `json:"amount"`
	Uid       uint32 `json:"uid"`
	OpenId    string `json:"openId"`
	Vip       bool   `json:"vip"`
	Code      string `json:"code"`
	Group     bool   `json:"group"`
	VipName   string `json:"vipName"`
	VipGive   uint32 `json:"vipGive"`
	VipExpiry string `json:"vipExpiry"`
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
	Type   string `json:"type,options=share|ad|group"`
	Status bool   `json:"status,optional"`
}

type TaskShareFollowRequest struct {
	OpenId string `json:"openId"`
}

type AIInfoResponse struct {
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	Call     string `json:"call"`
	Welcome  string `json:"welcome"`
	Status   bool   `json:"status"`
	RoleId   uint32 `json:"roleId"`
	ShowAd   bool   `json:"showAd"`
	RoleName string `json:"roleName"`
}

type AIEditRequest struct {
	Name   string `form:"name"`
	Call   string `form:"call"`
	RoleId uint32 `form:"roleId"`
	Status bool   `form:"status"`
}

type LoginWXQrcodeResponse struct {
	Qrcode   string `json:"qrcode"`
	SceneStr string `json:"sceneStr"`
}

type LoginCheckRequest struct {
	SceneStr string `form:"sceneStr"`
}

type WxqrcodeRequest struct {
	Channel string `form:"channel,optional"`
}

type ChatHistoryListResponse struct {
	Pagination Pagination        `json:"pagination"`
	Data       []ChatHistoryData `json:"data"`
}

type ChatHistoryData struct {
	Q      string `json:"q"`
	ChatId string `json:"chatId"`
	Time   string `json:"time"`
}

type SuanliHistoryListResponse struct {
	Pagination Pagination          `json:"pagination"`
	Data       []SuanliHistoryData `json:"data"`
}

type SuanliHistoryData struct {
	Amount int    `json:"amount"`
	Desc   string `json:"desc"`
	Time   string `json:"time"`
	Way    uint8  `json:"way"`
	Type   string `json:"type"`
}

type ToolsHistoryListResponse struct {
	Data []ToolsHistoryData `json:"data"`
}

type ToolsHistoryData struct {
	Key string `json:"key"`
}

type HashRateExchangeListResponse struct {
	Data []HashRateExchange `json:"data"`
}

type HashRateExchange struct {
	Date   string `json:"date"`
	Amount uint32 `json:"amount"`
	Use    uint32 `json:"use"`
	Expiry string `json:"expiry"`
	Status uint8  `json:"status"`
}

type UserNotifyUnreadResponse struct {
	Status bool `json:"status"`
}

type UserNotifyListResponse struct {
	Data []UserNotifyResponse `json:"data"`
}

type UserNotifyResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  bool   `json:"status"`
	Time    string `json:"time"`
}

type Response struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PageRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type ReportRequest struct {
	ChatId  string `json:"chatId,optional"`
	Content string `json:"content"`
}

type WorkRequest struct {
	ChatId    string `json:"chatId,optional"`
	Use       string `json:"use"`
	Introduce string `json:"introduce"`
	Content   string `json:"content"`
}

type ReportResponse struct {
	Data string `json:"data"`
}

type ImageRequest struct {
	Content string `json:"content"`
	Model   string `json:"model,optional,options=DALL-E|GPT-PLUS|StableDiffusion|Midjourney"`
	Number  int    `json:"number,optional,options=1|2|4"`
	Clarity string `json:"clarity,optional,options=standard|high|superhigh"`
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

type ImageMultiResponse struct {
	Url []string `json:"url"`
}

type ImageMultiAsyncResponse struct {
	Task  string   `json:"task"`
	Model string   `json:"model"`
	Url   []string `json:"url"`
}

type ImageEditRequest struct {
	Content string `form:"content"`
}

type PicRepairRequest struct {
	Image string `json:"image"`
}

type ImageTaskRequest struct {
	Task  string `json:"task"`
	Model string `json:"model"`
}

type Image2TextRequest struct {
	Image string `json:"image"`
}

type Image2TextResponse struct {
	Data string `json:"data"`
}

type QiMingRequest struct {
	First    string `json:"first"`
	Number   string `json:"number"`
	Birthday string `json:"birthday"`
	Sex      string `json:"sex"`
	Fix      string `json:"fix,optional"`
	Other    string `json:"other,optional"`
	ChatId   string `json:"chatId,optional"`
}

type JieMengRequest struct {
	Content string `json:"content"`
	ChatId  string `json:"chatId,optional"`
}

type SuanMingRequest struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Sex      string `json:"sex,optional"`
	Content  string `json:"content,optional"`
	ChatId   string `json:"chatId,optional"`
}

type GSQiMingRequest struct {
	Industry string `json:"industry"`
	Range    string `json:"range"`
	Culture  string `json:"culture,optional"`
	Other    string `json:"other,optional"`
	ChatId   string `json:"chatId,optional"`
}

type YYQiMingRequest struct {
	Name   string `json:"name"`
	Sex    string `json:"sex"`
	Other  string `json:"other,optional"`
	ChatId string `json:"chatId,optional"`
}

type HoroscopeRequest struct {
	Birthday      string `json:"birthday,optional"`
	Constellation string `json:"constellation,optional"`
	ChatId        string `json:"chatId,optional"`
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
	ChatId  string `json:"chatId,optional"`
}

type ExamRequest struct {
	Content string `json:"content"`
	Type    string `json:"type"`
	ChatId  string `json:"chatId,optional"`
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
	ChatId  string `json:"chatId,optional"`
}

type PlaygroundRequest struct {
	Content string `json:"content"`
	Lang    string `json:"lang,optional"`
	ChatId  string `json:"chatId,optional"`
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
	ChatId  string `json:"chatId,optional"`
}

type DiaryRequest struct {
	Content string `json:"content"`
	ChatId  string `json:"chatId,optional"`
}

type ArticleRequest struct {
	Subject string `json:"subject"`
	Type    string `json:"type"`
	Number  string `json:"number"`
	Content string `json:"content"`
	ChatId  string `json:"chatId,optional"`
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
	ChatId   string `json:"chatId,optional"`
}

type SalaryRequest struct {
	Content string `json:"content"`
	ChatId  string `json:"chatId,optional"`
}

type RejectRequest struct {
	Type    string `json:"type"`
	Way     string `json:"way"`
	Content string `json:"content"`
	ChatId  string `json:"chatId,optional"`
}

type PursueRequest struct {
	Content string `json:"content"`
	ChatId  string `json:"chatId,optional"`
}

type ChatResponse struct {
	Data string `json:"data"`
}

type ChatRequest struct {
	ChatId     string `json:"chatId"`
	TemplateId uint32 `json:"templateId"`
	Message    string `json:"message"`
	Model      string `json:"model,optional,options=GPT-3.5|GPT-4"`
	Image      string `json:"image,optional"`
}

type ChatHistoryRequest struct {
	ChatId string `path:"chatId"`
}

type ChatHistoryResponse struct {
	ChatId  string        `json:"chatId"`
	Model   string        `json:"model"`
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
	ChatId  string `json:"chatId,optional"`
}

type ConvertResponse struct {
	Data string `json:"data"`
}

type ValidRequest struct {
	Content string `json:"content"`
	Params  string `json:"params,optional"`
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

type QrCodeResponse struct {
	Data string `json:"data"`
}

type QrCodeRequest struct {
	Path  string `form:"path"`
	Scene string `form:"scene"`
}

type ShortLinkRequest struct {
	Page  string `json:"page"`
	Title string `json:"title"`
}

type UploadTokenResponse struct {
	Token string `json:"token"`
}

type ValidImageRequest struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type VipPriceResponse struct {
	Data []VipDataResponse `json:"data"`
}

type VipDataResponse struct {
	ID       uint32  `json:"id"`
	Name     string  `json:"name"`
	Origin   float32 `json:"origin"`
	Price    float32 `json:"price"`
	Amount   uint32  `json:"amount"`
	Day      uint32  `json:"day"`
	Discount float32 `json:"discount"`
}

type VipGiveResponse struct {
	Day    int    `json:"day"`
	Expiry string `json:"expiry"`
}

type VipCxchangeRequest struct {
	Code string `json:"code"`
}

type VipCxchangeResponse struct {
}

type VipCodeGenerateRequest struct {
	VipId  uint32 `json:"vipId"`
	Day    uint32 `json:"day"`
	AICode string `json:"aiCode"`
}

type VipCodeGenerateResponse struct {
	Code string `json:"code"`
}

type VipPrivilegeListResponse struct {
	Data []VipPrivilegeResponse `json:"data"`
}

type VipPrivilegeResponse struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

type HashRatePriceResponse struct {
	Data []HashRateResponse `json:"data"`
}

type HashRateResponse struct {
	ID       uint32  `json:"id"`
	Origin   float32 `json:"origin"`
	Price    float32 `json:"price"`
	VipPrice float32 `json:"vip_price"`
	Amount   uint32  `json:"amount"`
	Day      uint32  `json:"day"`
}

type HashRateCodeGenerateRequest struct {
	Day    uint32 `json:"day"`
	Amount uint32 `json:"amount"`
	AICode string `json:"aiCode"`
}

type HashRateCodeGenerateResponse struct {
	Code string `json:"code"`
}

type HashRateCxchangeRequest struct {
	Code string `json:"code"`
}

type HashRateCxchangeResponse struct {
}

type VipPayRequest struct {
	Platform string `json:"platform,options=wechat|alipay"`
}

type PayResponse struct {
	Data string `json:"data"`
}

type WeChatCallbackResponse struct {
}

type WechatValidateRequest struct {
	AppKey    string `path:"appkey"`
	Signature string `form:"signature,optional"`
	Timestamp string `form:"timestamp,optional"`
	Nonce     string `form:"nonce,optional"`
	Echostr   string `form:"echostr,optional"`
}

type WechatPayResponse struct {
	Data string `json:"data"`
}

type PayRequest struct {
	Type     string `path:"type"`
	Merchant string `path:"merchant"`
}

type CronResponse struct {
}
