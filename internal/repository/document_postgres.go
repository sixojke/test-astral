package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sixojke/test-astral/domain"
	"github.com/sixojke/test-astral/pkg/logger"
)

type DocumentPostgres struct {
	db *sqlx.DB
}

func NewDocumentPostgres(db *sqlx.DB) *DocumentPostgres {
	return &DocumentPostgres{
		db: db,
	}
}

type docsByUserIdHelp []docByUserIdHelp

type docByUserIdHelp struct {
	Id           string    `db:"id"`
	Name         string    `db:"name"`
	Mime         string    `db:"mime"`
	FilePath     string    `db:"file_path"`
	IsFile       bool      `db:"is_file"`
	IsPublic     bool      `db:"is_public"`
	DocumentData string    `db:"document_data"`
	Grants       string    `db:"grants"`
	CreatedAt    time.Time `db:"created_at"`
}

func (d *docsByUserIdHelp) prepare() *[]domain.Document {
	docsDirty := *d
	docs := make([]domain.Document, 0, len(docsDirty))
	for _, doc := range docsDirty {
		docs = append(docs, domain.Document{
			Id:           doc.Id,
			Name:         doc.Name,
			Mime:         doc.Mime,
			FilePath:     doc.FilePath,
			IsFile:       doc.IsFile,
			IsPublic:     doc.IsFile,
			DocumentData: doc.DocumentData,
			Grants:       strings.Split(doc.Grants, ","),
			CreatedAt:    doc.CreatedAt,
		})
	}

	return &docs
}

func (r *DocumentPostgres) Create(document *domain.Document, userId string) error {
	logger.Debugf("create document: params=[%v]", *document)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if rerr := tx.Rollback(); rerr != nil && err == nil {
			err = fmt.Errorf("failed to rollback transaction: %w", rerr)
		}
	}()

	query := `
		INSERT INTO documents (
		   	name,
		   	mime,
		   	file_path,
			is_file,
		   	is_public,
		   	document_data,
		   	user_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
	  	) RETURNING
			id
	`

	var documentId string
	if err := tx.QueryRow(query, document.Name, document.Mime, document.FilePath, document.IsFile,
		document.IsPublic, document.DocumentData, userId).Scan(&documentId); err != nil {
		logger.Errorf("failed to insert document: %v", err)
		return err
	}

	query = `
		INSERT INTO access_grants (
			document_id,
			user_id
		) VALUES (
			$1, $2
		)
	`

	if _, err := tx.Exec(query, documentId, userId); err != nil {
		logger.Errorf("failed to insert grant: grant=%v, documentId=%v: %v", userId, documentId, err)
		return err
	}

	if len(document.Grants) > 0 {
		for _, grant := range document.Grants {
			query = `
		  		WITH user_info AS (
					SELECT id
					FROM users
					WHERE login = $1 AND id != $2
		  		)
		  		INSERT INTO access_grants (document_id, user_id)
		  		SELECT $3, id
		  		FROM user_info;
			`

			if _, err := tx.Exec(query, grant, userId, documentId); err != nil {
				logger.Errorf("failed to insert grant: grant=%v, documentId=%v: %v", grant, documentId, err)
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *DocumentPostgres) GetCurrentUserDocuments(currentUserId string, params *domain.FilterParams) (*[]domain.Document, error) {
	logger.Debugf("get current user documents: params=[currentUserId=%v params=%v]", currentUserId, *params)

	query := `
	  SELECT 
		  d.id,
		  d.name,
		  d.mime,
		  d.file_path,
		  d.is_file,
		  d.is_public,
		  COALESCE(STRING_AGG(u.login, ','), '') AS grants,
		  d.created_at
	  FROM documents d
	  LEFT JOIN access_grants ag ON d.id = ag.document_id
	  LEFT JOIN users u ON ag.user_id = u.id
	  WHERE d.user_id = $1
	`

	args := []interface{}{currentUserId}

	if params.Key != "" && params.Value != "" && isValidField(params.Key) {
		query += " AND d." + params.Key + " LIKE $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, params.Value)
	}

	query += `
	  GROUP BY d.id, d.name, d.mime, d.file_path, d.is_file, d.is_public, d.document_data, d.created_at
	  ORDER BY d.created_at ASC
	  LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2) + `;
	`

	args = append(args, params.Limit, params.Offset)

	var docsDirty docsByUserIdHelp
	if err := r.db.Select(&docsDirty, query, args...); err != nil {
		logger.Errorf("failed to get current user documents: %v", err)
		return nil, err
	}

	return docsDirty.prepare(), nil
}

func (r *DocumentPostgres) GetOtherUserDocuments(userId string, currentUserId string, params *domain.FilterParams) (*[]domain.Document, error) {
	logger.Debugf("get other user documents: params=[userId=%v currentUserId=%v params=%v]", userId, currentUserId, *params)

	query := `
	  SELECT 
		d.id,
		d.name,
		d.mime,
		d.file_path,
		d.is_file,
		d.is_public,
		COALESCE(STRING_AGG(u.login, ','), '') AS grants
	  FROM documents d
	  JOIN access_grants ag ON d.id = ag.document_id
	  JOIN users u ON ag.user_id = u.id
	  WHERE 
		d.user_id = $1
		AND d.id IN (
		  SELECT ag.document_id
		  FROM access_grants ag
		  WHERE ag.user_id = $2
		) 
		OR (
			d.user_id = $1
			AND d.is_public = true
		)
	`

	args := []interface{}{userId, currentUserId}

	if params.Key != "" && params.Value != "" && isValidField(params.Key) {
		query += " AND d." + params.Key + " LIKE $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, params.Value)
	}

	query += `
	  GROUP BY d.id, d.name, d.mime, d.file_path, d.is_file, d.is_public, d.document_data
	  ORDER BY d.created_at ASC
	  LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2) + `;
	`

	args = append(args, params.Limit, params.Offset)

	var docsDirty docsByUserIdHelp
	if err := r.db.Select(&docsDirty, query, args...); err != nil {
		logger.Errorf("failed to get other user documents: %v", err)
		return nil, err
	}

	return docsDirty.prepare(), nil
}

func (r *DocumentPostgres) GetById(documentId, userId string) (*domain.Document, error) {
	logger.Debugf("get document by id: params=[documentId=%v userId=%v]", documentId, userId)

	query := `
		SELECT 
			d.id,
			d.name,
			d.mime,
			d.file_path,
			d.is_file,
			d.is_public,
			d.document_data,
			d.created_at
  		FROM documents d
  		WHERE 
	  		d.id = $1
	  		AND (
				d.is_public = TRUE 
	  			OR d.user_id = $2
				OR EXISTS (
					SELECT 1
					FROM access_grants ag
					WHERE ag.document_id = d.id
					AND ag.user_id = $3
				)
	  		)
	`

	var document domain.Document
	if err := r.db.Get(&document, query, documentId, userId, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDocumentNotFound
		}

		logger.Errorf("failed to get document: %v", err)
		return nil, err
	}

	query = `
		SELECT u.login
		FROM access_grants ag
		JOIN users u ON ag.user_id = u.id
		WHERE ag.document_id = $1;
	`

	var grants []string
	if err := r.db.Select(&grants, query, documentId); err != nil {
		logger.Errorf("failed to get grants: %v", err)
		return nil, err
	}

	document.Grants = grants

	return &document, nil
}

func (r *DocumentPostgres) CheckById(documentId, userId string) (bool, error) {
	logger.Debugf("check document by id: params=[documentId=%v userId=%v]", documentId, userId)

	query := `
	SELECT 
		1
	FROM documents d
	WHERE 
		  d.id = $1
		  AND (
			d.is_public = TRUE 
			OR EXISTS (
				  SELECT 1
				  FROM access_grants ag
				  WHERE ag.document_id = d.id
				  AND ag.user_id = $2
			)
		  )
`

	var exists bool
	if err := r.db.Get(&exists, query, documentId, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.ErrDocumentNotFound
		}

		logger.Errorf("failed to get document: %v", err)
		return false, err
	}

	return exists, nil
}

func (r *DocumentPostgres) Delete(documentId, userId string) (filePath string, err error) {
	logger.Debugf("delete document: params=[documentId=%v]", documentId)
	query := `
		SELECT file_path
		FROM documents
		WHERE id = $1 AND user_id = $2
	`

	if err := r.db.Get(&filePath, query, documentId, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrDocumentNotFound
		}

		logger.Errorf("failed to get file_path: %v", err)
		return "", err
	}

	query = `
		DELETE FROM documents
		WHERE id = $1 AND user_id = $2
	`

	_, err = r.db.Exec(query, documentId, userId)
	if err != nil {
		logger.Errorf("failed to delete document: %v", err)
		return "", err
	}

	return filePath, nil
}

func isValidField(field string) bool {
	validFields := []string{"name", "mime", "file_path", "is_file", "is_public", "document_data", "created_at"}
	for _, v := range validFields {
		if v == field {
			return true
		}
	}
	return false
}
