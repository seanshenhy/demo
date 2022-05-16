package service

import (
	"context"
	v1 "demo/api/realworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
)

// 登录
func (s *RealworldService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.UserReply, error) {
	ud, err := s.uc.Login(ctx, req.User.Email, req.User.Password)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			UserId:   int64(ud.UserID),
			Username: ud.Username,
			Email:    ud.Email,
			Image:    ud.Image,
			Token:    ud.Token,
			Bio:      ud.Bio,
		},
	}, nil
}

// 注册
func (s *RealworldService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.UserReply, error) {
	if req.User.Username == "" || req.User.Email == "" || req.User.Password == "" {
		return nil, errors.New(422, "params", "error")
	}
	u, err := s.uc.Register(ctx, req.User.Username, req.User.Email, req.User.Password)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			UserId:   int64(u.UserID),
			Username: u.Username,
			Token:    u.Token,
			Email:    u.Email,
		},
	}, nil
}

// 获取当前用户
func (s *RealworldService) GetCurrentUser(ctx context.Context, req *v1.GetCurrentUserRequest) (*v1.UserReply, error) {
	u, err := s.uc.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			UserId:   int64(u.UserID),
			Username: u.Username,
			Email:    u.Email,
			Token:    u.Token,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}, nil

}

// 更新用户简介
func (s *RealworldService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UserReply, error) {
	u, err := s.uc.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			UserId:   int64(u.UserID),
			Username: u.Username,
			Email:    u.Email,
			Token:    u.Token,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}, nil
}

// 获取用户简介
func (s *RealworldService) GetProfile(ctx context.Context, req *v1.GetProfileRequest) (*v1.ProfileReply, error) {
	p, err := s.uc.GetProfile(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &v1.ProfileReply{
		Profile: &v1.ProfileReply_Profile{
			UserId:    int64(p.UserID),
			Username:  p.Username,
			Bio:       p.Bio,
			Image:     p.Image,
			Following: p.Following,
		},
	}, nil
}

// 关注用户
func (s *RealworldService) FollowUser(ctx context.Context, req *v1.FollowUserRequest) (*v1.ProfileReply, error) {
	p, err := s.uc.FollowUser(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &v1.ProfileReply{
		Profile: &v1.ProfileReply_Profile{
			UserId:    int64(p.UserID),
			Username:  p.Username,
			Bio:       p.Bio,
			Image:     p.Image,
			Following: p.Following,
		},
	}, nil
}

// 取消用户关注
func (s *RealworldService) UnfollowUser(ctx context.Context, req *v1.UnfollowUserRequest) (*v1.ProfileReply, error) {
	p, err := s.uc.UnFollowUser(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &v1.ProfileReply{
		Profile: &v1.ProfileReply_Profile{
			UserId:    int64(p.UserID),
			Username:  p.Username,
			Bio:       p.Bio,
			Image:     p.Image,
			Following: p.Following,
		},
	}, nil
}
