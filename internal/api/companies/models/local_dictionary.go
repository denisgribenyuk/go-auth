package models

type LocalDictionary struct {
	KeyID       int64  `json:"key_id" db:"key_id"`
	LanguageISO string `json:"language_iso" db:"language_iso"`
	Value       string `json:"value" db:"value"`
}
