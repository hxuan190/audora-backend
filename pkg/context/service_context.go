package context

import (
	"music-app-backend/pkg/redis"
	"music-app-backend/pkg/storage"

	goflakeid "github.com/capy-engineer/go-flakeid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceContext struct {
	DBContext      *gorm.DB
	Router         *gin.Engine
	IDGenerator    *goflakeid.Generator
	RedisClient    *redis.Client
	StorageService *storage.MinIOService
}

func NewServiceContext(DBContext *gorm.DB,
	Router *gin.Engine,
	IDGenerator *goflakeid.Generator,
	RedisClient *redis.Client,
	StorageService *storage.MinIOService) *ServiceContext {
	return &ServiceContext{
		DBContext:      DBContext,
		Router:         Router,
		IDGenerator:    IDGenerator,
		RedisClient:    RedisClient,
		StorageService: StorageService,
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
func (ctx ServiceContext) GetRedisClient() *redis.Client {
	return ctx.RedisClient
}

func (ctx ServiceContext) GetStorageService() *storage.MinIOService {
	return ctx.StorageService
}
