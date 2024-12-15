package mocks

import (
	"context"

	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"go-micro.dev/v4/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

type usersAllAllowedService struct {
}

// CheckPermissions implements rms_users.RmsUsersService.
func (u *usersAllAllowedService) CheckPermissions(ctx context.Context, in *rms_users.CheckPermissionsRequest, opts ...client.CallOption) (*rms_users.CheckPermissionsResponse, error) {
	return &rms_users.CheckPermissionsResponse{Allowed: true}, nil
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
