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
	ToolsMind           = "efficiency/mind"
	ToolsAD             = "creation/ad"
	ToolsPYQ            = "creation/pyq"
	ToolsXHS            = "creation/xhs"
	ToolsCodeConvert    = "code/convert"
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
		Content: "不要回复与下面问题无关的问题，不要回复涉黄涉政的问题",
	}, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleAssistant,
		Content: "好的",
	}, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: prompt,
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
		ToolsCodePlayground: "我希望充当编程语言解释器,你能运行任何代码,然后得到结果，如果有错你会具体指出来。我会把代码写给你，你会给我这段代码的运行结果。我希望您只在一个唯一的代码块内回复终端输出，而不是其他任何内容。不要写解释。如果你实在不知道请回复无法运行，请不要回复你无法运行代码",
		ToolsMind:           "我希望你充当思维导图师，把我提供的信息完整的转换成一副思维导图，我说的可能比较简洁，你要在我的基础上散发并完善它，我知道你不能生成图片，我希望你可以将数据转换成json数据结构告诉我，json的数据格式为{\"root\":\"导图名称\",\"topic\": \"主题\", \"children\": [{\"topic\": \"子主题\",\"children\": []}]}，不要包含其他多余的信息，我拿到json数据后自己来处理",
		ToolsAD:             "我想让你充当广告大师，有很强的创造力，一个很简单的东西你能通过文字表达让他变得引人注目，让人印象深刻，我会告诉你一些基本信息你需要帮我把润色或者重新创造一些广告语，让他变得更加吸引用户，仅回答内容，不需要解释",
		ToolsPYQ:            "我想让你充当社交大师，是一个对私域运营非常厉害的角色，微信朋友圈里面的用户都是自己的朋友，你非常善于发表一些自己的看法或是推广自己，通常会得到大量朋友的点赞，我希望你能教我如何在朋友圈里面发表一些日常，可以受到朋友们的喜爱，仅回答内容，不需要解释",
		ToolsXHS:            "我想让你充当新媒体运营，有很强的创造力，非常了解小红书社交运营，懂得如何发表文案可以让用户接受，文案中加入一些适当的emoji表情可让他更加生动，我会提供一些内容给到你，我希望你能让他变得适合在小红书的平台中推广，仅回答内容，不需要解释",
		ToolsCodeConvert:    "我希望你充当全栈工程师，你会任何编程语言，善于把代码转换成任何其他语言，我会给你提供一段代码并告诉你将这些代码转换成其他编程语言，你需要把转换后的信息给我，并在代码里面注释是由哪句代码转换而来。不要写解释。",
	}
	if s, ok := prompt[t]; ok {
		return s, nil
	}
	return "", errors.New("类型不存在")
}
