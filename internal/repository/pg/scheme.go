package pg

import (
	"encoding/json"
	"fmt"
	"myapp/internal/entity"
)

func (s *Database) GetsSchemeById(id string) (entity.Scheme, error) {
	q := `
		SELECT coalesce((
			SELECT to_json(c)
	  		FROM schemes as c
	  		WHERE id = $1
		), '{}'::json)
	`
	var u entity.Scheme
	var data []byte
	err := s.QueryRow(q, id).Scan(&data)
	if err != nil {
		return u, fmt.Errorf("GetsSchemeById Scan: %v", err)
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return u, fmt.Errorf("GetsSchemeById Unmarshal: %v", err)
	}
	return u, nil
}

func (s *Database) GetsSchemeByLichka(lichka string) (entity.Scheme, error) {
	q := `
		SELECT coalesce((
			SELECT to_json(c)
	  		FROM schemes as c
	  		WHERE lichka = $1
		), '{}'::json)
	`
	var u entity.Scheme
	var data []byte
	err := s.QueryRow(q, lichka).Scan(&data)
	if err != nil {
		return u, fmt.Errorf("GetsSchemeByLichka Scan: %v", err)
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return u, fmt.Errorf("GetsSchemeByLichka Unmarshal: %v", err)
	}
	return u, nil
}

func (s *Database) EditSchemeById(id, sc_name string, sc_idx int) error {
	q := `
		UPDATE schemes SET
			sc_name = $1,
			sc_idx = $2
		WHERE id = $3
	`
	_, err := s.Exec(q, sc_name, sc_idx, id)
	if err != nil {
		return fmt.Errorf("EditScheme Exec err: %v", err)
	}
	return nil
}
