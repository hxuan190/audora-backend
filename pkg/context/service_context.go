package context

import (
	goflakeid "github.com/capy-engineer/go-flakeid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceContext struct {
	DBContext   *gorm.DB
	Router      *gin.Engine
	IDGenerator *goflakeid.Generator
}

func NewerviceContext(DBContext *gorm.DB,
	Router *gin.Engine,
	IDGenerator *goflakeid.Generator) *ServiceContext {
	return &ServiceContext{
		DBContext:   DBContext,
		Router:      Router,
		IDGenerator: IDGenerator,
	}
}

func (ctx ServiceContext) GetDB() *gorm.DB {
	return ctx.DBContext
}

func (ctx ServiceContext) GetRouter() *gin.Engine {
	return ctx.Router
}

func (ctx ServiceContext) GetIDGenerator() *goflakeid.Generator {
	return ctx.IDGenerator
}
