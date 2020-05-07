package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"time"
)

// Cap is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Cap struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Number      int       `json:"number" db:"number" form:"Number"`
	BelongsTo   *User     `json:"belongs_to,omitempty" belongs_to:"user"`
	BelongsToID uuid.UUID `json:"-" db:"userid"`
}

// String is not required by pop and may be deleted
func (c Cap) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Caps is not required by pop and may be deleted
type Caps []Cap

func (cap *Cap) Create(tx *pop.Connection) (*validate.Errors, error) {
	//nil, errors.New("error")
	return tx.ValidateAndCreate(cap)
}

func GetCounts(tx *pop.Connection) (map[int]int, error) {
	x := map[int]int{}

	type Countmap struct {
		Number int `db:"number"`
		Count  int `db:"count"`
	}

	type Countmaps []Countmap

	counts := Countmaps{}

	q := tx.Q()
	err := q.RawQuery("SELECT number, Count(Number) AS count FROM caps GROUP BY Number").All(&counts)
	if err != nil {
		return x, errors.WithStack(err)
	}

	for _, count := range counts {
		x[count.Number] = count.Count
	}

	return x, nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Cap) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.IntIsPresent{
			Name:  "number is present",
			Field: c.Number,
		},
		&validators.IntIsGreaterThan{
			Name:     "geq 0",
			Field:    c.Number,
			Compared: -1,
		},
		&validators.IntIsLessThan{
			Name:     "leq 100",
			Field:    c.Number,
			Compared: 100,
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Cap) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Cap) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
