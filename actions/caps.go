package actions

import (
	"github.com/6lmpnl/sternibingo/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

func CapsView (c buffalo.Context) error {
	//c.Set("cap", models.Cap{})

	caps := []models.Cap{}

	tx := c.Value("tx").(*pop.Connection)
	u := c.Value("current_user").(*models.User)
	err := tx.Q().Where("userid = ?", u.ID).Order("number ASC").All(&caps)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("caps", caps)

	return c.Render(200, r.HTML("/caps/view.plush.html"))
}

func CapCreate (c buffalo.Context) error {
	cap := &models.Cap{}
	if err := c.Bind(cap); err != nil {
		errors.WithStack(err)
	}


	u := c.Value("current_user").(*models.User)
	cap.BelongsTo = u


	tx := c.Value("tx").(*pop.Connection)
	verrs, err := cap.Create(tx)
	if (err != nil) {
		return errors.WithStack(err)
	}
	if (verrs.HasAny()) {
		c.Set("errors", verrs)
	}

	return c.Redirect(301, c.Request().Referer())
}

func CapRemove (c buffalo.Context) error {
	return nil
}