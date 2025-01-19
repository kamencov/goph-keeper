package sqlite

import (
	"errors"
	"goph-keeper/internal/services/workers"
	"time"
)

var (
	errParseTime = errors.New("failed to parse updated_at")
)

// GetAllSync - возвращает все данные из базы данных.
func (s *Storage) GetAllSync() ([]*workers.SyncModel, error) {
	var syncs []*workers.SyncModel
	var updatedAt string

	query := `SELECT id, user_id, table_name, task_id, action, updated_at FROM sync_client`

	rows, err := s.storage.Query(query)
	if err != nil {
		s.log.Error("failed to get all sync", "error", err)
		return nil, err
	}

	for rows.Next() {
		var sync workers.SyncModel

		if err := rows.Scan(&sync.ID, &sync.UserID, &sync.TableName, &sync.TaskID, &sync.Action, &updatedAt); err != nil {
			s.log.Error("failed to scan sync row", "error", err)
			return nil, err
		}

		// Парсинг строки в time.Time
		parsedTime, err := time.Parse("2006-01-02 15:04:05", updatedAt)
		if err != nil {
			s.log.Error("failed to parse updated_at", "error", err)
			return nil, errParseTime
		}
		sync.UpdatedAt = parsedTime

		syncs = append(syncs, &sync)
	}

	return syncs, nil
}

// GetDataCredentials - возвращает данные из базы данных.
func (s *Storage) GetDataCredentials(userID, taskID int) (*workers.Credentials, error) {
	var cred workers.Credentials
	var updatedAt string

	query := `SELECT * FROM credentials WHERE user_id = $1 AND id = $2`

	if err := s.storage.QueryRow(query, userID, taskID).Scan(&cred.ID, &cred.UserID, &cred.Resource, &cred.Login, &cred.Password, &updatedAt, &cred.DeleteTask); err != nil {
		s.log.Error("failed to get data credentials", "error", err)
		return nil, err
	}

	// Парсинг строки в time.Time
	parsedTime, err := time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		s.log.Error("failed to parse updated_at", "error", err)
		return nil, errParseTime
	}
	cred.UpdatedAt = parsedTime

	return &cred, nil
}

// GetDataTextData - возвращает данные из базы данных.
func (s *Storage) GetDataTextData(userID, taskID int) (*workers.TextData, error) {
	var cred workers.TextData
	var updatedAt string

	query := `SELECT * FROM text_data WHERE user_id =$1 AND id = $2`

	if err := s.storage.QueryRow(query, userID, taskID).Scan(&cred.ID, &cred.UserID, &cred.Text, &updatedAt, &cred.DeleteTask); err != nil {
		s.log.Error("failed to get data text data", "error", err)
		return nil, err
	}

	// Парсинг строки в time.Time
	parsedTime, err := time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		s.log.Error("failed to parse updated_at", "error", err)
		return nil, errParseTime
	}
	cred.UpdatedAt = parsedTime

	return &cred, nil
}

// GetDataBinaryData - возвращает данные из базы данных.
func (s *Storage) GetDataBinaryData(userID, taskID int) (*workers.BinaryData, error) {
	var cred workers.BinaryData
	var updatedAt string

	query := `SELECT * FROM binary_data WHERE user_id = $1 AND id = $2`

	if err := s.storage.QueryRow(query, userID, taskID).Scan(&cred.ID, &cred.UserID, &cred.Binary, &updatedAt, &cred.DeleteTask); err != nil {
		s.log.Error("failed to get data binary data", "error", err)
		return nil, err
	}

	// Парсинг строки в time.Time
	parsedTime, err := time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		s.log.Error("failed to parse updated_at", "error", err)
		return nil, errParseTime
	}
	cred.UpdatedAt = parsedTime

	return &cred, nil
}

// GetDataCards - возвращает данные из базы данных.
func (s *Storage) GetDataCards(userID, taskID int) (*workers.Cards, error) {
	var cred workers.Cards
	var updatedAt string

	query := `SELECT * FROM cards WHERE user_id = $1 AND id = $2`

	if err := s.storage.QueryRow(query, userID, taskID).Scan(&cred.ID, &cred.UserID, &cred.Cards, &updatedAt, &cred.DeleteTask); err != nil {
		s.log.Error("failed to get data cards", "error", err)
		return nil, err
	}

	// Парсинг строки в time.Time
	parsedTime, err := time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		s.log.Error("failed to parse updated_at", "error", err)
		return nil, errParseTime
	}
	cred.UpdatedAt = parsedTime

	return &cred, nil
}

// ClearSyncCredentials - удаляет все данные из базы данных.
func (s *Storage) ClearSyncCredentials() error {

	query := `DELETE FROM sync_client WHERE table_name = 'credentials'`

	_, err := s.storage.Exec(query)
	if err != nil {
		s.log.Error("failed to clear sync credentials", "error", err)
		return err
	}

	return nil
}

// ClearSyncTextData - удаляет все данные из базы данных.
func (s *Storage) ClearSyncTextData() error {

	query := `DELETE FROM sync_client WHERE table_name = 'text_data'`

	_, err := s.storage.Exec(query)
	if err != nil {
		s.log.Error("failed to clear sync text data", "error", err)
		return err
	}

	return nil
}

// ClearSyncBinaryData - удаляет все данные из базы данных.
func (s *Storage) ClearSyncBinaryData() error {

	query := `DELETE FROM sync_client WHERE table_name = 'binary_data'`

	_, err := s.storage.Exec(query)
	if err != nil {
		s.log.Error("failed to clear sync binary data", "error", err)
		return err
	}

	return nil
}

// ClearSyncCards - удаляет все данные из базы данных.
func (s *Storage) ClearSyncCards() error {

	query := `DELETE FROM sync_client WHERE table_name = 'cards'`

	_, err := s.storage.Exec(query)
	if err != nil {
		s.log.Error("failed to clear sync cards", "error", err)
		return err
	}

	return nil
}
