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

	title := req.InduceTitle
	if title == "" {
		title = "您有一条未读消息"
	}

	subtitle := req.InduceSubtitle
	imageURL := req.InduceImageURL

	now := time.Now()
	result, err := db.GetDB().Exec(
		"INSERT INTO shortlinks (code, original_url, created_at, updated_at, 诱导标题, 诱导副标题, 诱导图片URL) VALUES (?, ?, ?, ?, ?, ?, ?)",
		code, req.URL, now, now, title, subtitle, imageURL,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.Shortlink{
		ID:             id,
		Code:           code,
		OriginalURL:    req.URL,
		CreatedAt:      now,
		UpdatedAt:      now,
		IsDeleted:      false,
		IsDisabled:     false,
		TotalVisits:    0,
		InduceTitle:    title,
		InduceSubtitle: subtitle,
		InduceImageURL: imageURL,
	}, nil
}

func (s *ShortlinkService) GetByCode(code string) (*models.Shortlink, error) {
	if !utils.IsValidCode(code) {
		return nil, errors.New("invalid code format")
	}

	var shortlink models.Shortlink
	err := db.GetDB().QueryRow(
		"SELECT id, code, original_url, created_at, updated_at, is_deleted, is_disabled, total_visits, total_duration, COALESCE(诱导标题, ''), COALESCE(诱导副标题, ''), COALESCE(诱导图片URL, '') FROM shortlinks WHERE code = ? AND is_deleted = FALSE",
		code,
	).Scan(
		&shortlink.ID, &shortlink.Code, &shortlink.OriginalURL, &shortlink.CreatedAt, &shortlink.UpdatedAt,
		&shortlink.IsDeleted, &shortlink.IsDisabled, &shortlink.TotalVisits, &shortlink.TotalDuration,
		&shortlink.InduceTitle, &shortlink.InduceSubtitle, &shortlink.InduceImageURL,
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
		"SELECT id, code, original_url, created_at, updated_at, is_deleted, is_disabled, total_visits, total_duration, COALESCE(诱导标题, ''), COALESCE(诱导副标题, ''), COALESCE(诱导图片URL, '') FROM shortlinks WHERE is_deleted = FALSE ORDER BY created_at DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shortlinks []models.Shortlink
	for rows.Next() {
		var sl models.Shortlink
		if err := rows.Scan(&sl.ID, &sl.Code, &sl.OriginalURL, &sl.CreatedAt, &sl.UpdatedAt, &sl.IsDeleted, &sl.IsDisabled, &sl.TotalVisits, &sl.TotalDuration, &sl.InduceTitle, &sl.InduceSubtitle, &sl.InduceImageURL); err != nil {
			return nil, 0, err
		}
		shortlinks = append(shortlinks, sl)
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
		"SELECT id, code, original_url, created_at, updated_at, is_deleted, is_disabled, total_visits, total_duration, COALESCE(诱导标题, ''), COALESCE(诱导副标题, ''), COALESCE(诱导图片URL, '') FROM shortlinks WHERE id = ?",
		id,
	).Scan(
		&shortlink.ID, &shortlink.Code, &shortlink.OriginalURL, &shortlink.CreatedAt, &shortlink.UpdatedAt,
		&shortlink.IsDeleted, &shortlink.IsDisabled, &shortlink.TotalVisits, &shortlink.TotalDuration,
		&shortlink.InduceTitle, &shortlink.InduceSubtitle, &shortlink.InduceImageURL,
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
