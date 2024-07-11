package pg

import (
	"encoding/json"
	"fmt"
	"myapp/internal/entity"
	"myapp/pkg/my_time_parser"
	"time"
)

func (s *Database) AddNewUser(id int, username, firstname string) error {
	tn := time.Now().In(my_time_parser.Msk).Format(my_time_parser.MyRfc)
	q := `
		INSERT INTO users (id, username, firstname, created_at)
			VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`
	_, err := s.Exec(q, id, username, firstname, tn)
	if err != nil {
		return fmt.Errorf("AddNewUser Exec err: %s", err)
	}
	return nil
}

func (s *Database) DeleteUserById(id int) error {
	q := `
		DELETE FROM users WHERE id = $1
	`
	_, err := s.Exec(q, id)
	if err != nil {
		return fmt.Errorf("DeleteUserById Exec err: %s", err)
	}
	return nil
}

func (s *Database) DeleteUserByUsername(username string) error {
	q := `
		DELETE FROM users WHERE username = $1
	`
	_, err := s.Exec(q, username)
	if err != nil {
		return fmt.Errorf("DeleteUserByUsername Exec err: %s", err)
	}
	return nil
}

func (s *Database) GetUserById(id int) (entity.User, error) {
	q := `
		SELECT coalesce((
			SELECT to_json(c)
	  		FROM users as c
	  		WHERE id = $1
		), '{}'::json)
	`
	var u entity.User
	var data []byte
	err := s.QueryRow(q, id).Scan(&data)
	if err != nil {
		return u, fmt.Errorf("GetUserById Scan: %v", err)
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return u, fmt.Errorf("GetUserById Unmarshal: %v", err)
	}
	return u, nil
}

func (s *Database) GetUsersByStep(step string) ([]entity.User, error) {
	q := `
		SELECT coalesce((
			SELECT json_agg(c)
	  		FROM users as c
	  		WHERE step = $1
		), '[]'::json)
	`
	u := make([]entity.User, 0)
	var data []byte
	err := s.QueryRow(q, step).Scan(&data)
	if err != nil {
		return u, fmt.Errorf("GetUsersByStep Scan: %v", err)
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return u, fmt.Errorf("GetUsersByStep Unmarshal: %v", err)
	}
	return u, nil
}

func (s *Database) GetUsersByRef(ref string) ([]entity.User, error) {
	q := `
		SELECT coalesce((
			SELECT json_agg(c)
	  		FROM users as c
	  		WHERE ref = $1
		), '[]'::json)
	`
	u := make([]entity.User, 0)
	var data []byte
	err := s.QueryRow(q, ref).Scan(&data)
	if err != nil {
		return u, fmt.Errorf("GetUsersByRef Scan: %v", err)
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return u, fmt.Errorf("GetUsersByRef Unmarshal: %v", err)
	}
	return u, nil
}

func (s *Database) GetUserByUsername(username string) (entity.User, error) {
	q := `
		SELECT coalesce((
			SELECT to_json(c)
	  		FROM users as c
	  		WHERE username = $1
		), '{}'::json)
	`
	u := entity.User{}
	var data []byte
	err := s.QueryRow(q, username).Scan(&data)
	if err != nil {
		return u, fmt.Errorf("GetUsersByUsername Scan: %v", err)
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return u, fmt.Errorf("GetUsersByUsername Unmarshal: %v", err)
	}
	return u, nil
}

func (s *Database) GetAllUsers() ([]entity.User, error) {
	q := `
		SELECT coalesce((
			SELECT json_agg(c)
	  		FROM users as c
		), '[]'::json)
	`
	u := make([]entity.User, 0)
	var data []byte
	err := s.QueryRow(q).Scan(&data)
	if err != nil {
		return u, fmt.Errorf("GetAllUsers Scan: %v", err)
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return u, fmt.Errorf("GetAllUsers Unmarshal: %v", err)
	}
	return u, nil
}

func (s *Database) EditAdmin(id, is_admin int) error {
	q := `UPDATE users SET is_admin = $1 WHERE id = $2`
	_, err := s.Exec(q, is_admin, id)
	if err != nil {
		return fmt.Errorf("EditAdmin Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditLives(id, lives int) error {
	q := `UPDATE users SET lives = $1 WHERE id = $2`
	_, err := s.Exec(q, lives, id)
	if err != nil {
		return fmt.Errorf("EditLives Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditBotState(id int, bot_state string) error {
	q := `UPDATE users SET bot_state = $1 WHERE id = $2`
	_, err := s.Exec(q, bot_state, id)
	if err != nil {
		return fmt.Errorf("EditBotState Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditEmail(id int, email string) error {
	q := `UPDATE users SET email = $1 WHERE id = $2`
	_, err := s.Exec(q, email, id)
	if err != nil {
		return fmt.Errorf("EditEmail Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditLichka(id int, lichka string) error {
	q := `UPDATE users SET lichka = $1 WHERE id = $2`
	_, err := s.Exec(q, lichka, id)
	if err != nil {
		return fmt.Errorf("EditLichka Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditStep(id int, step string) error {
	q := `UPDATE users SET step = $1 WHERE id = $2`
	_, err := s.Exec(q, step, id)
	if err != nil {
		return fmt.Errorf("EditStep Exec err: %v", err)
	}
	return nil
}

func (s *Database) UpdateLatsActiontime(id int) error {
	tn := time.Now().In(my_time_parser.Msk).Format(my_time_parser.MyRfc)

	q := `UPDATE users SET lats_action_time = $1 WHERE id = $2`
	_, err := s.Exec(q, tn, id)
	if err != nil {
		return fmt.Errorf("UpdateLatsActiontime Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditLatsActiontime(id int, lats_action_time string) error {
	q := `UPDATE users SET lats_action_time = $1 WHERE id = $2`
	_, err := s.Exec(q, lats_action_time, id)
	if err != nil {
		return fmt.Errorf("EditLatsActiontime Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditIsLastPush(id, is_last_push int) error {
	q := `UPDATE users SET is_last_push = $1 WHERE id = $2`
	_, err := s.Exec(q, is_last_push, id)
	if err != nil {
		return fmt.Errorf("EditIsLastPush Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditIsFinal(id, is_final int) error {
	q := `UPDATE users SET is_final = $1 WHERE id = $2`
	_, err := s.Exec(q, is_final, id)
	if err != nil {
		return fmt.Errorf("EditIsFinal Exec err: %v", err)
	}
	return nil
}

func (s *Database) UpdateFeedbackTime(id int) error {
	tn := time.Now().In(my_time_parser.Msk).Format(my_time_parser.MyRfc)

	q := `UPDATE users SET feedback_time = $1 WHERE id = $2`
	_, err := s.Exec(q, tn, id)
	if err != nil {
		return fmt.Errorf("UpdateFeedbackTime Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditFeedbackTime(id int, feedback_time string) error {
	q := `UPDATE users SET feedback_time = $1 WHERE id = $2`
	_, err := s.Exec(q, feedback_time, id)
	if err != nil {
		return fmt.Errorf("EditFeedbackTime Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditdFeedbackCnt(id, feedback_cnt int) error {
	q := `UPDATE users SET feedback_cnt = $1 WHERE id = $2`
	_, err := s.Exec(q, feedback_cnt, id)
	if err != nil {
		return fmt.Errorf("EditdFeedbackCnt Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditIsSendPush(id, is_send_push int) error {
	q := `UPDATE users SET is_send_push = $1 WHERE id = $2`
	_, err := s.Exec(q, is_send_push, id)
	if err != nil {
		return fmt.Errorf("EditIsSendPush Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditRef(id int, ref string) error {
	q := `UPDATE users SET ref = $1 WHERE id = $2`
	_, err := s.Exec(q, ref, id)
	if err != nil {
		return fmt.Errorf("EditRef Exec err: %v", err)
	}
	return nil
}

func (s *Database) EditNotDelMessId(id int, not_del_mess_id int) error {
	q := `UPDATE users SET not_del_mess_id = $1 WHERE id = $2`
	_, err := s.Exec(q, not_del_mess_id, id)
	if err != nil {
		return fmt.Errorf("EditNotDelMessId Exec err: %v", err)
	}
	return nil
}
