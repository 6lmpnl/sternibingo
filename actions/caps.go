package actions

import (
	"fmt"
	"github.com/6lmpnl/sternibingo/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
	"github.com/pkg/errors"
	"net/http"
)

func CapsView(c buffalo.Context) error {
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

func CapCreate(c buffalo.Context) error {
	cap := &models.Cap{}
	if err := c.Bind(cap); err != nil {
		errors.WithStack(err)
	}

	u := c.Value("current_user").(*models.User)
	cap.BelongsTo = u

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := cap.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
	}

	return c.Redirect(301, c.Request().Referer())
}

func DestroyCap(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Field
	cap := &models.Cap{}

	// To find the Field the parameter field_id is used.
	if err := tx.Find(cap, c.Param("capid")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(cap); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "cap.deleted.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/caps")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(cap))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(cap))
	}).Respond(c)
}