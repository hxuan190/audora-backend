package context

import (
	goflakeid "github.com/capy-engineer/go-flakeid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type serviceContext struct {
	DBContext   *gorm.DB
	Router      *gin.Engine
	IDGenerator *goflakeid.Generator
}

func NewerviceContext(DBContext *gorm.DB,
	Router *gin.Engine,
	IDGenerator *goflakeid.Generator) *serviceContext {
	return &serviceContext{
		DBContext:   DBContext,
		Router:      Router,
		IDGenerator: IDGenerator,
	}
}
