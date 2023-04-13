package model

import (
	"errors"
	gogpt "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"time"
)

const (
	ToolsReportDay      = "report/day"
	ToolsReportWeek     = "report/week"
	ToolsReportWork     = "report/work"
	ToolsReportPlot     = "report/plot"
	ToolsSuanMing       = "divination/suanming"
	ToolsJieMeng        = "divination/jiemeng"
	ToolsIntroduce      = "chat/introduce"
	ToolsPursue         = "chat/pursue"
	ToolsReject         = "chat/reject"
	ToolsTranslate      = "convert/translate"
	ToolsActivity       = "creation/activity"
	ToolsDiary          = "creation/diary"
	ToolsArticle        = "creation/article"
	ToolsCodeGenerate   = "code/generate"
	ToolsCodeRegular    = "code/regular"
	ToolsCodePlayground = "code/playground"
)

type Record struct {
	ID          uint32    `gorm:"primary_key" json:"id"`
	Uid         uint32    `json:"uid" gorm:"index:idx_uid"`
	Type        string    `gorm:"type:varchar(20)" json:"type"`
	ChatId      string    `gorm:"type:varchar(50)" json:"chat_id"`
	Title       string    `gorm:"type:varchar(50)" json:"title"`
	Content     string    `json:"content" gorm:"type:text"`
	ShowContent string    `json:"show_content" gorm:"type:text"`
	Result      string    `json:"result" gorm:"type:text"`
	ShowResult  string    `json:"show_result" gorm:"type:text"`
	Model       string    `json:"model" gorm:"type:varchar(20)"`
	Platform    string    `json:"platform" gorm:"type:varchar(20)"`
	IsDelete    bool      `json:"is_delete" gorm:"default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;index:idx_created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt    time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (r Record) GetMessage(db *gorm.DB, user AIUser) ([]gogpt.ChatCompletionMessage, bool, error) {
	var message []gogpt.ChatCompletionMessage
	prompt, err := getContent(r.Type)
	if err != nil {
		return nil, false, err
	}
	message = append(message, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: "不要回复与下面问题无关的问题",
	}, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: prompt,
	}, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleAssistant,
		Content: "好的",
	})

	isFirst := true
	if r.ChatId != "" {
		maxToken := 300
		strLen := 100
		if user.IsVip() {
			maxToken = 800
			strLen = 200
		}
		var records []Record
		db.Raw("select id, content, LEFT(result, ?) as result from gpt_record where uid = ? and chat_id = ? and is_delete = 0 order by id desc limit ?", strLen, r.Uid, r.ChatId, 10).Scan(&records)
		if len(records) > 0 {
			isFirst = false
			totalLen := 0
			for i := len(records) - 1; i >= 0; i-- {
				message = append(message, gogpt.ChatCompletionMessage{
					Role:    "user",
					Content: records[i].Content,
				})
				message = append(message, gogpt.ChatCompletionMessage{
					Role:    "assistant",
					Content: records[i].Result,
				})

				totalLen += len([]rune(records[i].Content)) + len([]rune(records[i].Result))
				if totalLen > maxToken {
					break
				}
			}
		}
	}

	return message, isFirst, nil
}

func getContent(t string) (string, error) {
	prompt := map[string]string{
		ToolsReportWeek:     "请帮我把以下的工作内容填充为一篇完整的周报包含本周内容、下周计划、本周总结,用 markdown 格式以分点叙述的形式输出",
		ToolsReportDay:      "请帮我把以下的工作内容填充为一篇完整的日报，包含今日工作内容、明天工作计划以及总结,用 markdown 格式以分点叙述的形式输出",
		ToolsReportPlot:     "从现在开始你要充当一名编剧，想一些很创意的剧本，从想出有趣的角色、故事的背景、角色之间的对话等开始。一旦你的角色发展完成——创造一个充满曲折的激动人心的故事情节，让观众一直悬念到最后",
		ToolsReportWork:     "从现在开始你要充当一名职业导师，根据我的基本情况，帮助我完成一份述职报告，需要包含个人信息、工作职责、工作成果、工作总结、个人总结、工作计划、对公司的建议等",
		ToolsSuanMing:       "从现在开始你要充当一名占卜大师，结合我的情况给一份详细的算命报告，包含八字分析、五行分析、命理分析、事业分析、爱情分析、财运分析等相关内容，请用markdown格式输出",
		ToolsJieMeng:        "从现在开始你要充当周公，结合我的梦境，给我详细解释一下这个梦所预示的含义",
		ToolsIntroduce:      "从现在开始你要充当演讲大师，结合我的情况，给我写一份自我介绍，希望可以让大家很快记住我",
		ToolsPursue:         "从现在开始你化身为恋爱大师，根据我的情况，教我一步步去追求TA，需要包含具体的计划、步骤、行动等以及如何应对被拒绝的尴尬",
		ToolsReject:         "从现在开始你化身为沟通大师，教我如何去拒绝别人的请求，能够很好的缓解尴尬",
		ToolsTranslate:      "我希望你能担任翻译官、拼写校对和修辞改进的角色。我会用很多语言和你交流，你会识别语言，将其翻译并用更为优美和精炼的语句回答我。请将我简单的词汇和句子替换成更为优美和高雅的表达方式，确保意思不变，但使其更具文学性。请仅回答更正和改进的部分，不要写解释",
		ToolsActivity:       "从现在开始你化身为首席运营官，根据我的要求，为我详细设计一个运营方案，包括但不限于前期准备、活动的实施方案、活动过程跟踪、效果不及预期的方案、活动效果、需要的支持等",
		ToolsDiary:          "从现在开始你化身为文学大家，根据我提供的信息完成一篇日记",
		ToolsArticle:        "从现在开始你化身为文学大家，根据我提供的信息完成一篇佳作",
		ToolsCodeGenerate:   "我希望你能担任全栈开发工程师，可以结合我的需求一步步教我如何实现以及如何编写代码",
		ToolsCodeRegular:    "我希望你充当正则表达式生成器。您的角色是生成匹配文本中特定模式的正则表达式。您应该以一种可以轻松复制并粘贴到支持正则表达式的文本编辑器或编程语言中的格式提供正则表达式。不要写正则表达式如何工作的解释或例子；只需提供正则表达式本身",
		ToolsCodePlayground: "我希望你表现得像一个编程语言解释器。我会把代码写给你，你会用解释器的输出来响应。我希望您只在一个唯一的代码块内回复终端输出，而不是其他任何内容。不要写解释。除非我指示您这样做，否则不要键入命令",
	}
	if s, ok := prompt[t]; ok {
		return s, nil
	}
	return "", errors.New("类型不存在")
}
