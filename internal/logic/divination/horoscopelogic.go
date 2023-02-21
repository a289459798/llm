package divination

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HoroscopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHoroscopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HoroscopeLogic {
	return &HoroscopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HoroscopeLogic) Horoscope(req *types.HoroscopeRequest, w http.ResponseWriter) (resp *types.DivinationResponse, err error) {
	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	horoscope := req.Constellation

	if horoscope == "" {
		// 计算星座
		if req.Birthday == "" {
			return nil, errors.New("参数错误")
		}
		horoscope = getHoroscope(req.Birthday)
	}

	today := time.Now().Format("2006-01-02")
	record := &model.Record{}
	l.svcCtx.Db.Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").Where("type = ?", "divination/horoscope").Where("content = ?", horoscope).Find(record)
	if record.Result != "" {
		for _, v := range record.Result {
			w.Write([]byte(utils.EncodeURL(fmt.Sprintf("%c", v))))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(100 * time.Millisecond)
		}

		return
	}

	prompt := fmt.Sprintf("详细介绍一下%s，%s的运势，包含爱情、事业、健康、总结等方面的信息，请用markdown格式输出", today, horoscope)
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           prompt,
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}
	// 创建上下文
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

	ch := make(chan struct{})

	stream, err := l.svcCtx.GptClient.CreateCompletionStream(ctx, gptReq)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	result := ""
	go func() {
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Printf("Stream error: %v\n", err)
				break
			}
			if len(response.Choices) > 0 {
				w.Write([]byte(utils.EncodeURL(response.Choices[0].Text)))
				result += response.Choices[0].Text
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}

		}
		w.Write([]byte(utils.EncodeURL("\n\n***\n\n*由于占星学并非科学，所以不要完全相信星座运势的准确性，而是把它们看作是提供给你灵感的一个有趣的方式。*")))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		close(ch)
	}()

	select {
	case <-ch:
		// 处理已完成
		logx.Infof("EventStream logic finished")
	case <-ctx.Done():
		// 处理被取消
		logx.Errorf("EventStream logic canceled")
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "divination/horoscope",
		Content: horoscope,
		Result:  result,
	})
	return
}

func getHoroscope(bir string) string {
	arr := strings.Split(bir, "-")
	month, _ := strconv.Atoi(arr[0])
	day, _ := strconv.Atoi(arr[1])
	month -= 1
	var (
		DAY_ARR = [12]int{20, 19, 21, 20, 21, 22, 23, 23, 23, 24, 23, 22}
		ZODIACS = [13]string{"摩羯座", "水瓶座", "双鱼座", "白羊座", "金牛座", "双子座", "巨蟹座", "狮子座", "处女座", "天秤座", "天蝎座", "射手座", "摩羯座"}
	)

	if day < DAY_ARR[month] {
		return ZODIACS[month]
	} else {
		return ZODIACS[month+1]
	}
}
