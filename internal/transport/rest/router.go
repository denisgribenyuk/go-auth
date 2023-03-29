package rest

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Router struct {
	*gin.Engine
}

func New() *Router {
	//log := logrus.New()

	r := gin.New()
	r.Use(gin.Recovery())
	//r.Use(ginlogrus.Logger(log))
	p := ginprometheus.NewPrometheus("gin")

	// Эта штука для того, чтобы в прометеус отправлялась аггрегация по урлам вида /api/v1/users/:id, вместо /api/v1/users/1
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.String()
		for _, p := range c.Params {
			url = strings.Replace(url, p.Value, ":"+p.Key, 1)
		}
		return url
	}

	p.UseWithAuth(r, gin.Accounts{"metrics": os.Getenv("SERVICE_SECRET_KEY")})
	r.Use(gin.Logger())

	return &Router{
		Engine: r,
	}
}
