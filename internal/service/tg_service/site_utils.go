package tg_service

import (
	"encoding/base64"
	"fmt"
)

func (srv *TgService) CreateBase64UserData(tg_id int, tg_username, tg_name string) string {
	user, _ := srv.Db.GetUserById(tg_id)
	lichka := user.Lichka
	lichkaId := 6405739421
	if srv.DelAt(lichka) == "markodinncov" {
		lichkaId = 6328098519
	}

	data := fmt.Sprintf(`{"tg_id": %d, "tg_username": "%s", "tg_name": "%s", "lichka_username": "%s", "lichka_tg_id": %d}`, tg_id, tg_username, tg_name, lichka, lichkaId)
	base64Str := base64.URLEncoding.EncodeToString([]byte(data))
	return base64Str
}

