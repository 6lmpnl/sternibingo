package actions

import (
	"github.com/6lmpnl/sternibingo/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

//UsersNew renders the users form
func UsersNew(c buffalo.Context) error {
	u := models.User{}
	c.Set("user", u)
	return c.Render(200, r.HTML("users/new.plush.html"))
}

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(200, r.HTML("users/new.plush.html"))
	}

	err = u.SendValidationEmail(tx)
	if err != nil {
		verrs.Add("Email", "Failed to send Activation-Mail.")
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(500, r.HTML("users/new.plush.html"))
	}

	//c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "I've sent you an activation link. You need to open it before you can log-in.")

	return c.Redirect(302, "/")
}

func UsersActivate(c buffalo.Context) error {
	u := []models.User{}
	tx := c.Value("tx").(*pop.Connection)

	code := c.Param("code")

	query := tx.Q()
	err := query.Where("validation_code = ?", code).All(&u)
	if err != nil {
		return err
	}

	for x := range u {
		u[x].ValidateEmail(tx)
	}

	c.Flash().Add("success", "Account Activated!")
	return c.Redirect(301, "/auth/new")
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/auth/new")
		}
		return next(c)
	}
}
