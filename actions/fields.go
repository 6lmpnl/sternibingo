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

// List gets all Fields. This function is mapped to the path
// GET /fields
func ShowFields(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	fields := &models.Fields{}

	// Retrieve all Fields from the DB
	if err := tx.All(fields); err != nil {
		return err
	}

	u := c.Value("current_user").(*models.User)
	counts, err := models.GetCountsForUser(tx, u)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("fields", fields)
	c.Set("counts", counts)

	return c.Render(200, r.HTML("/fields/index.plush.html"))
}

// Destroy deletes a Field from the DB. This function is mapped
// to the path DELETE /fields/{field_id}
func DestroyField(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Field
	field := &models.Field{}

	// To find the Field the parameter field_id is used.
	if err := tx.Find(field, c.Param("field_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(field); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "field.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/fields")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(field))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(field))
	}).Respond(c)
}
