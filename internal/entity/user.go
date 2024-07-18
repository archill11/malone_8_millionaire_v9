package entity

type User struct {
	Id              int                     `json:"id"`
	Username        string                  `json:"username"`
	Firstname       string                  `json:"firstname"`
	IsAdmin         int                     `json:"is_admin"`
	BotState        string                  `json:"bot_state"`
	Email           string                  `json:"email"`
	Ref             string                  `json:"ref"`
	Lichka          string                  `json:"lichka"`
	Lives           int                     `json:"lives"`
	Step            string                  `json:"step"`
	LatsActiontime  string                  `json:"lats_action_time"`
	IsLastPush      int                     `json:"is_last_push"`
	IsSendPush      int                     `json:"is_send_push"`
	IsFinal         int                     `json:"is_final"`
	FeedbackCnt     int                     `json:"feedback_cnt"`
	FeedbackTime    string                  `json:"feedback_time"`
	CreatedAt       string                  `json:"created_at"`
	NotDelMessId    int                     `json:"not_del_mess_id"`
	InstLink        string                  `json:"inst_link"`
	IsInstPush      int                     `json:"is_inst_push"`
}
