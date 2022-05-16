package biz

import (
	"context"
	v1 "demo/api/realworld/v1"
	"demo/internal/conf"
	"demo/internal/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
)

// DO
type User struct {
	UserID     int
	Email      string
	Username   string
	Token      string
	Bio        string
	Image      string
	PasswdHash string
}
type Follow struct {
	Following string
	Username  string
}
type Author struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}
type UserLogin struct {
	UserID   int
	Email    string
	Username string
	Token    string
	Bio      string
	Image    string
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	VerifyUserExistByEmail(ctx context.Context, email string) bool
	GetUserByUserID(ctx context.Context, id int) (*User, error)
	GetUserByUserName(ctx context.Context, name string) (*User, error)
	UpdateUser(ctx context.Context, user_id int, user *User) (*User, error)
}

type ProfileRepo interface {
	GetFollowByUserID(ctx context.Context, userId int) (*Follow, error)
	FollowUser(ctx context.Context, myUserId, userId int) (bool, error)
	UnfollowUser(ctx context.Context, myUserId, userId int) (bool, error)
}

type UserUsecase struct {
	ur   UserRepo
	pr   ProfileRepo
	jwtc *conf.JWT
	log  *log.Helper
}

func NewUserUseCase(ur UserRepo, pr ProfileRepo, jwtc *conf.JWT, logger log.Logger) *UserUsecase {
	return &UserUsecase{ur: ur, pr: pr, jwtc: jwtc, log: log.NewHelper(logger)}
}

func hashpassword(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

func verifyPassword(hashpwd, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashpwd), []byte(pwd)); err != nil {
		return false
	}
	return true
}

// 生成token
func (uc *UserUsecase) generateToken(u *User) string {
	token, err := auth.GenerateToken([]byte(uc.jwtc.Secret), u.Email, u.Username, u.UserID)
	if err != nil {
		panic("create token error")
	}
	return token
}

// 解析登录信息
func (uc *UserUsecase) ParseLoginInfo(ctx context.Context) auth.LoginUser {
	ctx, err := auth.ParseTokenByCtx(ctx, []byte(uc.jwtc.Secret))
	if err != nil {
		panic("parse token error")
	}
	return ctx.Value("loginUser").(auth.LoginUser)
}

// 注册
func (uc *UserUsecase) Register(ctx context.Context, username, email, password string) (*UserLogin, error) {
	// 注册前判断用户email是否存在
	if uc.ur.VerifyUserExistByEmail(ctx, email) {
		return nil, errors.New(422, "email", "has exist")
	}
	// 注册
	u := &User{
		Email:      email,
		Username:   username,
		PasswdHash: hashpassword(password),
	}
	if err := uc.ur.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return &UserLogin{
		UserID:   u.UserID,
		Email:    email,
		Username: username,
		Token:    uc.generateToken(u),
	}, nil
}

// 登录
func (uc *UserUsecase) Login(ctx context.Context, email, passwd string) (*UserLogin, error) {
	if len(email) == 0 {
		return nil, errors.New(422, "email", "cannot empty")
	}
	u, err := uc.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !verifyPassword(u.PasswdHash, passwd) {
		return nil, errors.Forbidden("user", "account or passwrod error")
	}
	return &UserLogin{
		UserID:   u.UserID,
		Username: u.Username,
		Email:    u.Email,
		Token:    uc.generateToken(u),
	}, nil
}

// 获取当前登录用户
func (uc *UserUsecase) GetCurrentUser(ctx context.Context) (*UserLogin, error) {
	// 通过jwt token解密得到username信息
	loginUser := uc.ParseLoginInfo(ctx)
	// 通过username查询用户
	u, err := uc.ur.GetUserByUserName(ctx, loginUser.Username)
	if err != nil {
		return nil, errors.NotFound("user", "not exist")
	}
	return &UserLogin{
		UserID:   u.UserID,
		Username: u.Username,
		Email:    u.Email,
		Token:    uc.generateToken(u),
		Bio:      u.Bio,
		Image:    u.Image,
	}, nil
}

// 更新用户简介
func (uc *UserUsecase) UpdateUser(ctx context.Context, uur *v1.UpdateUserRequest) (*UserLogin, error) {
	loginUser := uc.ParseLoginInfo(ctx)
	u := &User{
		Email:    uur.User.Email,
		Username: uur.User.Username,
		Bio:      uur.User.Bio,
		Image:    uur.User.Image,
	}
	if uur.User.Password != "" {
		u.PasswdHash = hashpassword(uur.User.Password)
	}
	u, err := uc.ur.UpdateUser(ctx, loginUser.UserID, u)
	if err != nil {
		return nil, err
	}
	return &UserLogin{
		UserID:   u.UserID,
		Username: u.Username,
		Email:    u.Email,
		Token:    uc.generateToken(u),
		Bio:      u.Bio,
		Image:    u.Image,
	}, nil
}

// 获取用户简介
func (uc *UserUsecase) GetProfile(ctx context.Context, userId int) (*Author, error) {
	// 获取用户信息
	u, err := uc.ur.GetUserByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &Author{
		UserID:    u.UserID,
		Username:  u.Username,
		Following: true,
		Bio:       u.Bio,
		Image:     u.Image,
	}, nil
}

// 关注用户
func (uc *UserUsecase) FollowUser(ctx context.Context, userId int) (*Author, error) {
	// 获取当前用户
	loginUser := uc.ParseLoginInfo(ctx)
	// 获取用户信息
	u, err := uc.ur.GetUserByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	if loginUser.UserID == userId {
		return nil, errors.New(422, "user", "cannot follow self")
	}
	// 获取关注信息
	_, err = uc.pr.FollowUser(ctx, loginUser.UserID, userId)
	if err != nil {
		return nil, err
	}
	return &Author{
		UserID:    u.UserID,
		Username:  u.Username,
		Following: true,
		Bio:       u.Bio,
		Image:     u.Image,
	}, nil
}

// 取消关注用户
func (uc *UserUsecase) UnFollowUser(ctx context.Context, userId int) (*Author, error) {
	// 获取当前用户
	loginUser := uc.ParseLoginInfo(ctx)
	// 获取用户信息
	u, err := uc.ur.GetUserByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	// 取消关注信息
	_, err = uc.pr.UnfollowUser(ctx, loginUser.UserID, userId)
	if err != nil {
		return nil, err
	}
	return &Author{
		Username:  u.Username,
		Following: false,
		Bio:       u.Bio,
		Image:     u.Image,
	}, nil
}
