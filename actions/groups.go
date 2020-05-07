package actions

import (
	"github.com/6lmpnl/sternibingo/models"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"path"

	"github.com/gobuffalo/buffalo"
)

// GroupsShow default implementation.
func GroupsShow(c buffalo.Context) error {
	groups := &models.Groups{}
	group := &models.Group{}
	current_user := c.Value("current_user").(*models.User)

	tx := c.Value("tx").(*pop.Connection)
	err := tx.Where("id IN (SELECT group_id FROM ingroup i WHERE user_id = ?)", current_user.ID).All(groups)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("group", group)
	c.Set("escape", url.PathEscape)
	c.Set("join", path.Join)
	c.Set("groups", groups)

	return c.Render(http.StatusOK, r.HTML("groups/show.html"))
}

func GroupShow(c buffalo.Context) error {
	groupname := c.Param("groupname")
	group := &models.Group{}

	tx := c.Value("tx").(*pop.Connection)
	err := tx.Eager("InGroup.User").Where("name = ?", groupname).First(group)
	if err != nil {
		errors.WithStack(err)
	}

	c.Set("group", group)

	return c.Render(http.StatusOK, r.HTML("groups/showgroup.html"))
}

// GroupsNew default implementation.
func GroupsNew(c buffalo.Context) error {
	c.Set("group", models.Group{})

	return c.Render(http.StatusOK, r.HTML("groups/new.html"))
}

// GroupsCreate
func CreateGroup(c buffalo.Context) error {
	//groupsPath := c.Value("groupsPath")//.(string)
	group := &models.Group{}
	c.Bind(group)

	tx := c.Value("tx").(*pop.Connection)
	u := c.Value("current_user").(*models.User)

	group.AddUser(u.ID, true)

	verrs, err := group.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("group", group)
		c.Set("verrs", verrs)
		return c.Redirect(302, "/groups/new")
	}
	c.Logger().Debug(group)

	return c.Redirect(302, "/groups/")
}
