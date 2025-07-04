// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package bumflixdb

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Film struct {
	ID               pgtype.UUID
	Title            string
	Slug             string
	Year             pgtype.Int4
	SourceKey        string
	TranscodedPrefix pgtype.Text
	Status           string
	CreatedAt        pgtype.Timestamptz
}
