package main

import (
	"flag"
	"fmt"
	"go-shop-server/user/global"
	"go-shop-server/user/handle"
	"go-shop-server/user/model"
	pb "go-shop-server/user/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net"
)

func main() {
	// 获取用户输入参数
	ip := flag.String("ip", "0.0.0.0", "ip地址")
	port := flag.Int("port", 50051, "端口号")

	flag.Parse()

	fmt.Printf("ip: %s\nport: %d\n", *ip, *port)

	// 注册grpc服务
	server := grpc.NewServer()
	pb.RegisterUserServer(server, &handle.User{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	// 启动服务
	err = server.Serve(listen)
	if err != nil {
		fmt.Printf("failed to grpc Serve: %v", err)
		return
	}

	// 初始化MySQL
	MysqlInit("root", "123456", "127.0.0.1", "3306", "go-shop-server")
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
