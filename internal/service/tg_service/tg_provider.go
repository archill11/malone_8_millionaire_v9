package tg_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"myapp/internal/models"
	"myapp/pkg/files"
	"net/http"
	"strconv"
)

func (srv *TgService) GetUpdates(offset, timeout int, token string) ([]models.Update, error) {
	json_data, err := json.Marshal(map[string]any{
		"offset":  offset,
		"timeout": timeout,
	})
	if err != nil {
		return []models.Update{}, fmt.Errorf("GetUpdates Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, token, "getUpdates"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return []models.Update{}, fmt.Errorf("GetUpdates Post err: %v", err)
	}
	defer resp.Body.Close()
	var cAny models.GetUpdatesResp
	if err := json.NewDecoder(resp.Body).Decode(&cAny); err != nil {
		return cAny.Result, fmt.Errorf("GetUpdates Decode err: %v", err)
	}
	if cAny.ErrorCode != 0 {
		return cAny.Result, fmt.Errorf("GetUpdates errResp: %+v", cAny.BotErrResp)
	}
	return cAny.Result, nil
}

func (srv *TgService) GetMe(token string) (models.GetMeResp, error) {
	resp, err := http.Get(fmt.Sprintf(srv.Cfg.TgEndp, token, "getMe"))
	if err != nil {
		return models.GetMeResp{}, fmt.Errorf("GetMe Get err: %v", err)
	}
	defer resp.Body.Close()
	var cAny models.GetMeResp
	if err := json.NewDecoder(resp.Body).Decode(&cAny); err != nil {
		return models.GetMeResp{}, fmt.Errorf("GetMe Decode err: %v", err)
	}
	if cAny.ErrorCode != 0 {
		return cAny, fmt.Errorf("GetMe errResp: %+v", cAny)
	}
	return cAny, nil
}

func (srv *TgService) GetChat(chatId int, token string) (models.GetChatResp, error) {
	json_data, err := json.Marshal(map[string]any{
		"chat_id": strconv.Itoa(chatId),
	})
	if err != nil {
		return models.GetChatResp{}, fmt.Errorf("GetChat Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, token, "getChat"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return models.GetChatResp{}, fmt.Errorf("GetChat Post err: %v", err)
	}
	defer resp.Body.Close()
	var cAny models.GetChatResp
	if err := json.NewDecoder(resp.Body).Decode(&cAny); err != nil {
		return models.GetChatResp{}, fmt.Errorf("GetChat Decode err: %v", err)
	}
	if cAny.ErrorCode != 0 {
		return cAny, fmt.Errorf("GetChat errResp: %+v", cAny)
	}
	return cAny, nil
}

func (srv *TgService) GetFile(fileId string) (models.GetFileResp, error) {
	url := fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, fmt.Sprintf("getFile?file_id=%s", fileId))
	resp, err := http.Get(url)
	if err != nil {
		return models.GetFileResp{}, fmt.Errorf("GetFile Get file_id: %s err: %v", fileId, err)
	}
	defer resp.Body.Close()
	var cAny models.GetFileResp
	if err := json.NewDecoder(resp.Body).Decode(&cAny); err != nil {
		return models.GetFileResp{}, fmt.Errorf("GetFile Decode err: %v", err)
	}
	if cAny.ErrorCode != 0 {
		return cAny, fmt.Errorf("GetFile errResp: %+v", cAny)
	}
	return cAny, nil
}

func (srv *TgService) GetChatMember(user_id, chatId int) (models.GetChatMemberResp, error) {
	json_data, err := json.Marshal(map[string]any{
		"chat_id": strconv.Itoa(chatId),
		"user_id": user_id,
	})
	if err != nil {
		return models.GetChatMemberResp{}, fmt.Errorf("GetChat Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "getChatMember"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return models.GetChatMemberResp{}, fmt.Errorf("GetChatMember Get err: %v", err)
	}
	defer resp.Body.Close()
	var cAny models.GetChatMemberResp
	if err := json.NewDecoder(resp.Body).Decode(&cAny); err != nil {
		return models.GetChatMemberResp{}, fmt.Errorf("GetChatMember Decode err: %v", err)
	}
	return cAny, nil
}

func (srv *TgService) SendForceReply(chat int, mess string) error {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":      strconv.Itoa(chat),
		"text":         mess,
		"reply_markup": `{"force_reply": true}`,
	})
	if err != nil {
		return fmt.Errorf("SendForceReply Decode err: %v", err)
	}
	err = srv.SendData(json_data, "sendMessage")
	if err != nil {
		return fmt.Errorf("SendForceReply sendData err: %v", err)
	}
	return nil
}

func (srv *TgService) SendMessage(chat_id int, text string) (models.SendMessageResp, error) {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":                  strconv.Itoa(chat_id),
		"text":                     text,
		"disable_web_page_preview": true,
	})
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessage Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendMessage"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessage Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMessageResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessage Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendMessage errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendMessageWRM(chat_id int, text, reply_markup string) (models.SendMessageResp, error) {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":                  strconv.Itoa(chat_id),
		"text":                     text,
		"disable_web_page_preview": true,
		"reply_markup": reply_markup,
		"parse_mode": "HTML",
	})
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessage Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendMessage"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessage Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMessageResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessage Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendMessage errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendMessageMarkdown(chat_id int, text string) (models.SendMessageResp, error) {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":                  strconv.Itoa(chat_id),
		"text":                     text,
		"disable_web_page_preview": true,
		"parse_mode":               "MarkdownV2",
	})
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessageMarkdown Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendMessage"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessageMarkdown Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMessageResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessageMarkdown Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendMessageMarkdown errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendMessageHTML(chat_id int, text string) (models.SendMessageResp, error) {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":                  strconv.Itoa(chat_id),
		"text":                     text,
		"disable_web_page_preview": true,
		"parse_mode":               "HTML",
	})
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessageHTML Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendMessage"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessageHTML Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMessageResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendMessageHTML Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendMessageHTML errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) CopyMessage(chat_id, from_chat_id, message_id int) (models.SendMessageResp, error) {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":      strconv.Itoa(chat_id),
		"from_chat_id": strconv.Itoa(from_chat_id),
		"message_id":   message_id,
	})
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("CopyMessage Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "copyMessage"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("CopyMessage Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMessageResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMessageResp{}, fmt.Errorf("CopyMessage Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("CopyMessage errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendVideoNote(body io.Reader, contentType string) (models.SendMediaResp, error) {
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendVideoNote"),
		contentType,
		body,
	)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideoNote Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMediaResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideoNote Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendVideoNote errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendVideoNoteCurrFile(chat_id int, file string) (models.SendMediaResp, error) {
	futureJson := map[string]string{
		"video_note":   fmt.Sprintf("@%s", file),
		"chat_id": strconv.Itoa(chat_id),
	}
	contentType, body, err := files.CreateForm(futureJson)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideoNoteCurrFile CreateForm err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendVideoNote"),
		contentType,
		body,
	)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideoNote Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMediaResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideoNote Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendVideoNote errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendAnimation(body io.Reader, contentType string) (models.SendMediaResp, error) {
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendAnimation"),
		contentType,
		body,
	)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendAnimation Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMediaResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendAnimation Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendAnimation errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendAnimationV2(chatId int, fileNameInServer string) (models.SendMediaResp, error) {
	futureJson := map[string]string{
		"chat_id": strconv.Itoa(chatId),
		"animation": fmt.Sprintf("@%s", fileNameInServer),
	}
	contentType, body, err := files.CreateForm(futureJson)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendAnimationV2 CreateForm err: %v", err)
	}
	resp, err := srv.SendAnimation(body, contentType)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendAnimationV2 Post err: %v", err)
	}
	return resp, nil
}

func (srv *TgService) DeleteMessage(chat_id, message_id int) error {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":    strconv.Itoa(chat_id),
		"message_id": strconv.Itoa(message_id),
	})
	if err != nil {
		return fmt.Errorf("DeleteMessage Marshal err: %v", err)
	}
	err = srv.SendData(json_data, "deleteMessage")
	if err != nil {
		return fmt.Errorf("DeleteMessage sendData err: %v", err)
	}
	return nil
}

func (srv *TgService) EditMessageText(chat_id, message_id int, text string) error {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":    strconv.Itoa(chat_id),
		"message_id": message_id,
		"text":       text,
	})
	if err != nil {
		return fmt.Errorf("EditMessageText Marshal err: %v", err)
	}
	_, err = http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "editMessageText"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return fmt.Errorf("EditMessageText Post err: %v", err)
	}
	return nil
}

func (srv *TgService) EditMessageTextHTML(chat_id, message_id int, text string) error {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":                  strconv.Itoa(chat_id),
		"message_id":               message_id,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	})
	if err != nil {
		return fmt.Errorf("EditMessageTextHTML Marshal err: %v", err)
	}
	_, err = http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "editMessageText"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return fmt.Errorf("EditMessageTextHTML Post err: %v", err)
	}
	return nil
}

func (srv *TgService) EditMessageReplyMarkup(chat_id, message_id int) error {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":    strconv.Itoa(chat_id),
		"message_id": message_id,
	})
	if err != nil {
		return fmt.Errorf("EditMessageReplyMarkup Marshal err: %v", err)
	}
	_, err = http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "editMessageReplyMarkup"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return fmt.Errorf("EditMessageReplyMarkup Post err: %v", err)
	}
	return nil
}

func (srv *TgService) SendVideo(body io.Reader, contentType string) (models.SendMediaResp, error) {
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendVideo"),
		contentType,
		body,
	)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideo Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMediaResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideo Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendVideo errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendVideoWCaption(chat_id int, caption, fileNameInServer string) (models.SendMediaResp, error) {
	futureJson := map[string]string{
		"chat_id":    strconv.Itoa(chat_id),
		"caption":    caption,
		"parse_mode": "HTML",
		"video":      fmt.Sprintf("@%s", fileNameInServer),
	}
	contentType, body, err := files.CreateForm(futureJson)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideoWCaption CreateForm err: %v", err)
	}
	resp, err := srv.SendVideo(body, contentType)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendVideoWCaption SendVideo err: %v", err)
	}
	return resp, nil
}

func (srv *TgService) SendPhoto(contentType string, body io.Reader) (models.SendMediaResp, error) {
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendPhoto"),
		contentType,
		body,
	)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendPhoto Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMediaResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendPhoto Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendPhoto errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendPhotoWCaption(chat_id int, caption, fileNameInServer string) (models.SendMediaResp, error) {
	futureJson := map[string]string{
		"chat_id":    strconv.Itoa(chat_id),
		"caption":    caption,
		"parse_mode": "HTML",
		"photo":      fmt.Sprintf("@%s", fileNameInServer),
	}
	contentType, body, err := files.CreateForm(futureJson)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendPhotoWCaption CreateForm err: %v", err)
	}
	resp, err := srv.SendPhoto(contentType, body)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendPhotoWCaption SendPhoto err: %v", err)
	}
	return resp, nil
}

func (srv *TgService) SendPhotoWCaptionWRM(chat_id int, caption, fileNameInServer, reply_markup string) (models.SendMediaResp, error) {
	futureJson := map[string]string{
		"chat_id":    strconv.Itoa(chat_id),
		"caption":    caption,
		"parse_mode": "HTML",
		"photo":      fmt.Sprintf("@%s", fileNameInServer),
		"reply_markup": reply_markup,
	}
	contentType, body, err := files.CreateForm(futureJson)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendPhotoWCaptionWRM CreateForm err: %v", err)
	}
	resp, err := srv.SendPhoto(contentType, body)
	if err != nil {
		return models.SendMediaResp{}, fmt.Errorf("SendPhotoWCaptionWRM SendPhoto err: %v", err)
	}
	return resp, nil
}

func (srv *TgService) SendContact(chat_id int, phone_number, first_name string) (models.SendMessageResp, error) {
	json_data, err := json.Marshal(map[string]any{
		"chat_id":      strconv.Itoa(chat_id),
		"phone_number": phone_number,
		"first_name":   first_name,
	})
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendContact Marshal err: %v", err)
	}
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, "sendContact"),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendContact Post err: %v", err)
	}
	defer resp.Body.Close()
	var j models.SendMessageResp
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return models.SendMessageResp{}, fmt.Errorf("SendContact Decode err: %v", err)
	}
	if j.ErrorCode != 0 {
		return j, fmt.Errorf("SendContact errResp: %+v", j.BotErrResp)
	}
	return j, nil
}

func (srv *TgService) SendData(json_data []byte, method string) error {
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, method),
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return fmt.Errorf("sendData Post err: %v", err)
	}
	defer resp.Body.Close()
	var cAny models.BotErrResp
	if err := json.NewDecoder(resp.Body).Decode(&cAny); err != nil {
		return fmt.Errorf("sendData Decode err: %v", err)
	}
	if cAny.ErrorCode != 0 {
		return fmt.Errorf("sendData ErrResp: %+v", cAny)
	}
	return nil
}

func (srv *TgService) SendDataV2(method string, contentType string, body io.Reader) error {
	resp, err := http.Post(
		fmt.Sprintf(srv.Cfg.TgEndp, srv.Cfg.Token, method),
		contentType,
		body,
	)
	if err != nil {
		return fmt.Errorf("sendDataV2 Post err: %v", err)
	}
	defer resp.Body.Close()
	var cAny models.BotErrResp
	if err := json.NewDecoder(resp.Body).Decode(&cAny); err != nil {
		return fmt.Errorf("sendDataV2 Decode err: %v", err)
	}
	if cAny.ErrorCode != 0 {
		return fmt.Errorf("sendDataV2 ErrResp: %+v", cAny)
	}
	return nil
}
