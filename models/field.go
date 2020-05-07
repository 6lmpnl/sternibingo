package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"time"
)

// Field is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Field struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	Field     slices.Int `json:"field" db:"field"`
}

// String is not required by pop and may be deleted
func (f Field) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

func (f Field) Rows() [][]int {
	return [][]int{f.Field[0:5],
		f.Field[5:10],
		f.Field[10:15],
		f.Field[15:20],
		f.Field[20:25]}
}

// Fields is not required by pop and may be deleted
type Fields []Field

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (f *Field) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.FuncValidator{
			Fn: func() bool {
				if len(f.Field) != 25 {
					return false
				}
				for c := range f.Field {
					if c <= 0 || c >= 100 {
						return false
					}
				}
				return true
			},
			Field:   "Field",
			Name:    "ValidField",
			Message: "Field is not Valid",
		}), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (f *Field) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (f *Field) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
