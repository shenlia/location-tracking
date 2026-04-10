package services

import (
	"database/sql"
	"errors"
	"time"

	"location-tracking-shortlink/db"
	"location-tracking-shortlink/models"
	"location-tracking-shortlink/utils"
)

type ShortlinkService struct{}

func NewShortlinkService() *ShortlinkService {
	return &ShortlinkService{}
}

func (s *ShortlinkService) Create(req *models.ShortlinkCreateRequest) (*models.Shortlink, error) {
	code, err := utils.GenerateCode()
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {
		var exists int
		err := db.GetDB().QueryRow("SELECT COUNT(*) FROM shortlinks WHERE code = ?", code).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			break
		}
		code, err = utils.GenerateCode()
		if err != nil {
			return nil, err
		}
	}

	模板 := req.诱导模板
	if 模板 == "" {
		模板 = "court"
	}

	标题 := req.诱导标题
	副标题 := req.诱导副标题
	图片URL := req.诱导图片URL

	for _, t := range models.预设模板库 {
		if t.ID == 模板 {
			if 标题 == "" {
				标题 = t.Title
			}
			if 副标题 == "" {
				副标题 = t.Subtitle
			}
			break
		}
	}

	now := time.Now()
	result, err := db.GetDB().Exec(
		"INSERT INTO shortlinks (code, original_url, created_at, updated_at, 诱导标题, 诱导副标题, 诱导图片URL, 诱导模板) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		code, req.URL, now, now, 标题, 副标题, 图片URL, 模板,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Shortlink{
		ID:          id,
		Code:        code,
		OriginalURL: req.URL,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsDeleted:   false,
		IsDisabled:  false,
		TotalVisits: 0,
		诱导标题:        标题,
		诱导副标题:       副标题,
		诱导图片URL:     图片URL,
		诱导模板:        模板,
	}, nil
}

func (s *ShortlinkService) GetByCode(code string) (*models.Shortlink, error) {
	if !utils.IsValidCode(code) {
		return nil, errors.New("invalid code format")
	}

	var shortlink models.Shortlink
	err := db.GetDB().QueryRow(
		"SELECT id, code, original_url, created_at, updated_at, is_deleted, is_disabled, total_visits, total_duration, COALESCE(诱导标题, ''), COALESCE(诱导副标题, ''), COALESCE(诱导图片URL, ''), COALESCE(诱导模板, 'court') FROM shortlinks WHERE code = ? AND is_deleted = FALSE",
		code,
	).Scan(
		&shortlink.ID, &shortlink.Code, &shortlink.OriginalURL, &shortlink.CreatedAt, &shortlink.UpdatedAt,
		&shortlink.IsDeleted, &shortlink.IsDisabled, &shortlink.TotalVisits, &shortlink.TotalDuration,
		&shortlink.诱导标题, &shortlink.诱导副标题, &shortlink.诱导图片URL, &shortlink.诱导模板,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("shortlink not found")
		}
		return nil, err
	}

	return &shortlink, nil
}

func (s *ShortlinkService) List(page, pageSize int) ([]models.Shortlink, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var total int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM shortlinks WHERE is_deleted = FALSE").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.GetDB().Query(
		"SELECT id, code, original_url, created_at, updated_at, is_deleted, is_disabled, total_visits, total_duration, COALESCE(诱导标题, ''), COALESCE(诱导副标题, ''), COALESCE(诱导图片URL, ''), COALESCE(诱导模板, 'court') FROM shortlinks WHERE is_deleted = FALSE ORDER BY created_at DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shortlinks []models.Shortlink
	for rows.Next() {
		var s models.Shortlink
		if err := rows.Scan(&s.ID, &s.Code, &s.OriginalURL, &s.CreatedAt, &s.UpdatedAt, &s.IsDeleted, &s.IsDisabled, &s.TotalVisits, &s.TotalDuration, &s.诱导标题, &s.诱导副标题, &s.诱导图片URL, &s.诱导模板); err != nil {
			return nil, 0, err
		}
		shortlinks = append(shortlinks, s)
	}

	return shortlinks, total, nil
}

func (s *ShortlinkService) Delete(id int64) error {
	_, err := db.GetDB().Exec("UPDATE shortlinks SET is_deleted = TRUE, updated_at = ? WHERE id = ?", time.Now(), id)
	return err
}

func (s *ShortlinkService) Toggle(id int64) (*models.Shortlink, error) {
	var isDisabled bool
	err := db.GetDB().QueryRow("SELECT is_disabled FROM shortlinks WHERE id = ?", id).Scan(&isDisabled)
	if err != nil {
		return nil, err
	}

	newState := !isDisabled
	_, err = db.GetDB().Exec("UPDATE shortlinks SET is_disabled = ?, updated_at = ? WHERE id = ?", newState, time.Now(), id)
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *ShortlinkService) GetByID(id int64) (*models.Shortlink, error) {
	var shortlink models.Shortlink
	err := db.GetDB().QueryRow(
		"SELECT id, code, original_url, created_at, updated_at, is_deleted, is_disabled, total_visits, total_duration, COALESCE(诱导标题, ''), COALESCE(诱导副标题, ''), COALESCE(诱导图片URL, ''), COALESCE(诱导模板, 'court') FROM shortlinks WHERE id = ?",
		id,
	).Scan(
		&shortlink.ID, &shortlink.Code, &shortlink.OriginalURL, &shortlink.CreatedAt, &shortlink.UpdatedAt,
		&shortlink.IsDeleted, &shortlink.IsDisabled, &shortlink.TotalVisits, &shortlink.TotalDuration,
		&shortlink.诱导标题, &shortlink.诱导副标题, &shortlink.诱导图片URL, &shortlink.诱导模板,
	)
	if err != nil {
		return nil, err
	}
	return &shortlink, nil
}

func (s *ShortlinkService) IncrementVisits(id int64, duration int64) error {
	_, err := db.GetDB().Exec(
		"UPDATE shortlinks SET total_visits = total_visits + 1, total_duration = total_duration + ?, updated_at = ? WHERE id = ?",
		duration, time.Now(), id,
	)
	return err
}
