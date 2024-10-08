package service

import (
	"database/sql"
	"errors"
	"fmt"
	"test-golang/internal/helpers"
	"test-golang/internal/models"
)

type ArticleService struct {
	db *sql.DB
}

func NewArticleService(db *sql.DB) *ArticleService {
	return &ArticleService{
		db: db,
	}
}

func (s *ArticleService) GetAll() ([]models.Article, error) {
	rows, err := s.db.Query(`
		SELECT id, title, content, created_at, updated_at 
		FROM articles
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var createdAt, updatedAt []uint8

		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		article.CreatedAt, err = helpers.ParseTimestamp(createdAt)
		if err != nil {
			return nil, err
		}

		article.UpdatedAt, err = helpers.ParseTimestamp(updatedAt)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}
	return articles, nil
}

func (s *ArticleService) GetByID(id int) (models.Article, error) {
	if s.db == nil {
		return models.Article{}, fmt.Errorf("nil pointer: ArticleService.db")
	}

	var article models.Article
	var createdAt, updatedAt []uint8

	err := s.db.QueryRow(`
        SELECT id, title, content, created_at, updated_at 
        FROM articles 
        WHERE id = ?
    `, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Article{}, fmt.Errorf("article with ID %d not found", id)
		}
		return models.Article{}, fmt.Errorf("database error: %w", err)
	}

	// Convert []uint8 to time.Time
	article.CreatedAt, err = helpers.ParseTimestamp(createdAt)
	if err != nil {
		return models.Article{}, fmt.Errorf("error parsing created_at: %w", err)
	}

	article.UpdatedAt, err = helpers.ParseTimestamp(updatedAt)
	if err != nil {
		return models.Article{}, fmt.Errorf("error parsing updated_at: %w", err)
	}

	return article, nil
}

func (s *ArticleService) Create(article models.Article) (models.Article, error) {
	result, err := s.db.Exec(`
		INSERT INTO articles (title, content) 
		VALUES (?, ?)
	`, article.Title, article.Content)
	if err != nil {
		return models.Article{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Article{}, err
	}

	return s.GetByID(int(id))
}

func (s *ArticleService) Update(article models.Article) (models.Article, error) {
	result, err := s.db.Exec(`
		UPDATE articles 
		SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, article.Title, article.Content, article.ID)
	if err != nil {
		return models.Article{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Article{}, err
	}
	if rowsAffected == 0 {
		return models.Article{}, errors.New("article not found")
	}

	return s.GetByID(article.ID)
}

func (s *ArticleService) Delete(id int) error {
	result, err := s.db.Exec(`
		DELETE FROM articles 
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("article not found")
	}
	return nil
}
