package grifts

import (
	"github.com/6lmpnl/sternibingo/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
