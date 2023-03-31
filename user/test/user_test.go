package test

import (
	"context"
	"fmt"
	"go-shop-server/user/global"
	"go-shop-server/user/model"
	pb "go-shop-server/user/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	MysqlInit("root", "123456", "127.0.0.1", "3306", "go-shop-server")

	var userList []model.User

	for i := 0; i < 10; i++ {
		userList = append(userList, model.User{
			Model: model.Model{
				CreatedAt: &model.LocalTime{Time: time.Now()},
				UpdatedAt: &model.LocalTime{Time: time.Now()},
				DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
			},
			Mobile:   fmt.Sprintf("1111111111%d", i),
			Password: "",
			NickName: fmt.Sprintf("sni%d", i),
			Birthday: &model.LocalTime{Time: time.Now()},
			Gender:   "",
			Role:     new(int),
		})
	}

	global.DB.CreateInBatches(&userList, 10)
}

func TestGetUserList(t *testing.T) {
	Init()
	defer conn.Close()

	list, err := userClient.GetUserList(context.Background(), &pb.PageInfo{
		PageNum:  1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("failed getUserList:%v", err)
		return
	}

	for _, response := range list.UserResponse {
		fmt.Println(response)
	}
}

var (
	userClient pb.UserClient
	conn       *grpc.ClientConn
)

func Init() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	userClient = pb.NewUserClient(conn)
}

// mysql连接初始化
func MysqlInit(username, pwd, host, port, dbname string) {
	// dsn := "root:123456@tcp(192.168.0.6:3306)/gin?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		username, pwd, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true, //打印sql
		//SkipDefaultTransaction: true, //禁用事务
	})

	if err != nil {
		fmt.Println(err)
	}

	// 自动迁移
	db.AutoMigrate(
		// 表
		model.User{},
	)

	global.DB = db
}
