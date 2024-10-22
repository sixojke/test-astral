package domain

import "time"

type Document struct {
	Id           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Mime         string `json:"mime" db:"mime"`
	FilePath     string `db:"file_path"`
	IsFile       bool   `json:"is_file" db:"is_file"`
	IsPublic     bool   `json:"is_public" db:"is_public"`
	DocumentData string `json:"json,omitempty" db:"document_data"`
	Grants       []string
	CreatedAt    time.Time `json:"created" db:"created_at"`
}
