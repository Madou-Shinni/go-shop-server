package handle

import (
	"context"
	"errors"
	"go-shop-server/user/global"
	"go-shop-server/user/model"
	pb "go-shop-server/user/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type User struct {
	*pb.UnimplementedUserServer
}

func structToPbStructPoint(user model.User) *pb.UserResponse {
	var role int32
	if user.Role != nil {
		role = int32(*user.Role)
	}
	return &pb.UserResponse{
		Id:        int32(user.ID),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		DeletedAt: user.DeletedAt.Time.String(),
		Mobile:    user.Mobile,
		Password:  user.Password,
		NickName:  user.NickName,
		Birthday:  user.Birthday.Format("2006-10-12"),
		Gender:    user.Gender,
		Role:      role,
	}
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (u *User) GetUserList(ctx context.Context, req *pb.PageInfo) (*pb.UserListResponse, error) {
	var (
		users      []model.User
		count      int64
		pbUserList []*pb.UserResponse
	)
	if err := global.DB.Model(&model.User{}).Count(&count).Error; err != nil {
		return nil, err
	}

	if err := global.DB.Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Find(&users).Error; err != nil {
		return nil, err
	}

	// struct to proto struct
	for _, v := range users {
		pbUserList = append(pbUserList, structToPbStructPoint(v))
	}

	res := &pb.UserListResponse{
		Count:        uint32(count),
		UserResponse: pbUserList,
	}

	return res, nil
}

func (u *User) GetUserMobile(ctx context.Context, req *pb.GetUserMobileRequest) (*pb.UserResponse, error) {
	var user model.User
	if err := global.DB.Where("mobile = ?", req.Mobile).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "用户不存在")
		}
		return nil, err
	}

	return structToPbStructPoint(user), nil
}

func (u *User) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*emptypb.Empty, error) {
	var user model.User
	if err := global.DB.First(&user, req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "用户不存在")
		}
		return nil, err
	}

	if err := global.DB.Where("id = ?", user.ID).Save(&req).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
