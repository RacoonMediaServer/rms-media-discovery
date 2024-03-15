package mocks

import (
	"context"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"go-micro.dev/v4/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

type usersAllAllowedService struct {
}

func (u usersAllAllowedService) GetPermissions(ctx context.Context, in *rms_users.GetPermissionsRequest, opts ...client.CallOption) (*rms_users.GetPermissionsResponse, error) {
	resp := rms_users.GetPermissionsResponse{
		Perms: []rms_users.Permissions{
			rms_users.Permissions_Search, rms_users.Permissions_AccountManagement,
		},
	}
	return &resp, nil
}

func (u usersAllAllowedService) RegisterUser(ctx context.Context, in *rms_users.User, opts ...client.CallOption) (*rms_users.RegisterUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u usersAllAllowedService) GetUserByTelegramId(ctx context.Context, in *rms_users.GetUserByTelegramIdRequest, opts ...client.CallOption) (*rms_users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u usersAllAllowedService) GetAdminUsers(ctx context.Context, in *emptypb.Empty, opts ...client.CallOption) (*rms_users.GetAdminUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewMockUsersAllAllowed() rms_users.RmsUsersService {
	return &usersAllAllowedService{}
}
