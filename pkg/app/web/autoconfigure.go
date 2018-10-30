package web

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/kataras/iris"
	irsctx "github.com/kataras/iris/context"
)

const Profile = "web"

type configuration struct {
	at.AutoConfiguration

	Properties properties `mapstructure:"web"`
}

func newWebConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newWebConfiguration)
}

// Context is the instance of context.Context
func (c *configuration) Context(app *iris.Application) context.Context {
	ctx := &Context{
		Context: irsctx.NewContext(app),
	}

	if c.Properties.View.Enabled {
		v := c.Properties.View
		app.RegisterView(iris.HTML(v.ResourcePath, v.Extension))

		app.Get(v.ContextPath, func(ctx iris.Context) {
			ctx.View(v.DefaultPage)
		})
	}

	app.ContextPool.Attach(func() irsctx.Context {
		return &Context{
			// Optional Part 3:
			Context: irsctx.NewContext(app),
		}
	})

	return ctx
}
