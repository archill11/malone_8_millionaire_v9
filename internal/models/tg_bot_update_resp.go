package models

type BotErrResp struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type GetUpdatesResp struct {
	Result     []Update `json:"result"`
	Parameters struct {
		MigrateToChatID int `json:"migrate_to_chat_id"`
		RetryAfter      int `json:"retry_after"`
	} `json:"parameters"`

	BotErrResp
}

type GetChatResp struct {
	Result Chat `json:"result"`
	BotErrResp
}

type GetMeResp struct {
	Result User `json:"result"`
	BotErrResp
}

type GetChatMemberResp struct {
	Result struct {
		Status string `json:"status"`
	} `json:"result"`
	BotErrResp
}

type GetFileResp struct {
	Result struct {
		File_id        string `json:"file_id"`
		File_unique_id string `json:"file_unique_id"`
		File_path      string `json:"file_path"`
	} `json:"result"`

	BotErrResp
}

type SendMediaGroupResp struct {
	Result []SendMediaRespResult `json:"result"`
	BotErrResp
}

type SendMediaResp struct {
	Result SendMediaRespResult `json:"result"`
	BotErrResp
}

type SendMediaRespResult struct {
	MessageId int         `json:"message_id"`
	Caption   string      `json:"caption"`
	Chat      Chat        `json:"chat"`
	Video     Video       `json:"video"`
	Photo     []PhotoSize `json:"photo"`
}

type SendMessageResp struct {
	Result SendMessageRespResult `json:"result"`
	BotErrResp
}

type SendMessageRespResult struct {
	MessageId int    `json:"message_id"`
	Text      string `json:"text"`
	Date      int    `json:"date"`
}
