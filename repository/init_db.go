package repository

import (
	"context"
	"github.com/CCpro10/micro_douyin/conf"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	initMySQL()
	//initRedis()
	initRepository()
}

var (
	// RedisCtx Redis相关全局变量
	RedisCtx context.Context
	RedisDB  *redis.Client

	// DB gorm全局变量
	DB *gorm.DB
)

func initRedis() {
	RedisCtx = context.Background()
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     conf.Config.Redis.Addr,
		Password: conf.Config.Redis.Password,
		DB:       conf.Config.Redis.DB,
	})
	_, err := RedisDB.Ping(RedisCtx).Result()
	if err != nil {
		panic(err)
	}
}

func initMySQL() {
	// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := conf.Config.MYSQL.Username + ":" +
		conf.Config.MYSQL.Password + "@tcp(" +
		conf.Config.MYSQL.Addr + ")/" +
		conf.Config.MYSQL.Database +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	}) //这里用短变量声明会有歧义
	if err != nil {
		panic(err)
	}

	creatDatabase()
}

func creatDatabase() {
	err := DB.AutoMigrate(&Comment{}, &Favorite{}, &Follow{}, &User{}, &Video{})
	if err != nil {
		panic(err)
	}
}

var (
	userRepo    IUserRepository
	videoRepo   IVideoRepository
	favorite    IFavoriteRepository
	followRepo  IFollowRepository
	commentRepo ICommentRepository
)

func initRepository() {
	userRepo = &UserRepository{}
	videoRepo = &VideoRepository{}
	favorite = &FavoriteRepository{}
	followRepo = &FollowRepository{}
	commentRepo = &CommentRepository{}
}

func GetUserRepository() IUserRepository {
	return userRepo
}

func GetVideoRepository() IVideoRepository {
	return videoRepo
}

func GetFavoriteRepository() IFavoriteRepository {
	return favorite
}

func GetFollowRepository() IFollowRepository {
	return followRepo
}

func GetCommentRepository() ICommentRepository {
	return commentRepo
}
