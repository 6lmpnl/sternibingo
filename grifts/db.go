package grifts

import (
	"errors"
	"github.com/6lmpnl/sternibingo/models"
	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			fields := [][]int {
				{3,27,36,89,70,
				54,90,60,11,48,
				19,45,2,73,24,
				67,30,52,26,98,
				71,13,84,05,33},

				{14,77,38,2,85,
				61,23,96,49,54,
				6,88,13,31,75,
				42,26,65,59,97,
				70,8,51,12,20},

				{91,70,39,63,7,
				25,46,82,18,54,
				79,68,13,22,86,
				2,57,35,94,40,
				43,1,72,26,80},
			}

			for f := 0; f<len(fields); f++ {
				field := &models.Field{Field: fields[f]}


				verrs, err := tx.ValidateAndCreate(field)
				if err != nil {
					return err;
				}
				if verrs.HasAny() {
					return errors.New(verrs.String())
				}
			}

			return nil
		})
	})

})
