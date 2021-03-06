package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	"github.com/6lmpnl/sternibingo/models"
	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_sternibingo_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		txm := popmw.Transaction(models.DB)
		app.Use(txm)

		// Setup and use translations:
		app.Use(translations())

		app.Use(useLayout)

		app.GET("/", HomeHandler)
		app.GET("/routes", ListRoutes)

		//AuthMiddlewares
		app.Use(SetCurrentUser)
		app.Use(Authorize)

		//Routes for Auth
		auth := app.Group("/auth")
		auth.GET("/", AuthLanding)
		auth.GET("/new", AuthNew)
		auth.POST("/", AuthCreate)
		auth.DELETE("/", AuthDestroy)
		auth.Middleware.Skip(Authorize, AuthLanding, AuthNew, AuthCreate)
		auth.Middleware.Skip(useLayout, AuthLanding, AuthNew, AuthCreate, AuthDestroy)

		//Routes for User registration
		users := app.Group("/users")
		users.GET("/new", UsersNew)
		users.POST("/", UsersCreate)
		//users.Middleware.Skip(txm, UsersCreate) // skip auto transaction middleware to handle transaction by hand
		users.GET("/activate/{code}", UsersActivate)
		users.Middleware.Remove(Authorize)
		users.Middleware.Remove(useLayout)

		//Routes for Caps
		caps := app.Group("/caps")
		caps.GET("/", CapsView)
		caps.POST("/", CapCreate)
		caps.DELETE("/{capid}", DestroyCap)

		//Routes for Fields
		fields := app.Group("/fields")
		fields.GET("/", ShowFields)

		//Routes for Groups
		app.GET("/groups/", GroupsShow)
		app.GET("/groups/new", GroupsNew)
		app.POST("/groups/", CreateGroup)
		app.GET("/groups/{groupname}", GroupShow)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func useLayout(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Value("route")
		c.Set("useLayout", true)
		return next(c)
	}
}
