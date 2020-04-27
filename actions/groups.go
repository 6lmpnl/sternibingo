package actions

import (
    "net/http"
    
	"github.com/gobuffalo/buffalo"
)

// GroupsShow default implementation.
func GroupsShow(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("groups/show.html"))
}

// GroupsNew default implementation.
func GroupsNew(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("groups/new.html"))
}

// GroupsCreate
func CreateGroup(c buffalo.Context) error {
	return c.Redirect(302, c.Value("groupsPath").(string));
}