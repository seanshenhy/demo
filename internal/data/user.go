package data

import (
	"context"
	"demo/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// 用户表
type User struct {
	ID         int    `gorm:"primarykey;auto_increment" json:"id"`
	CreatedAt  int    `gorm:"type:int(11);not null;default:0;comment:创建时间" json:"created_at"`
	UpdatedAt  int    `gorm:"type:int(11);not null;default:0;comment:更新时间" json:"updated_at"`
	DeletedAt  int    `gorm:"type:int(11);not null;default:0;comment:删除时间" json:"deleted_at"`
	Email      string `gorm:"type:varchar(64);not null;comment:邮箱" json:"email"`
	Username   string `gorm:"type:varchar(64);not null;comment:用户名" json:"username"`
	Bio        string `gorm:"type:varchar(128);not null;comment:简介" json:"bio"`
	Image      string `gorm:"type:varchar(128);not null;comment:图片" json:"image"`
	PasswdHash string `gorm:"type:varchar(255);not null;comment:密码" json:"passwdhash"`
}

// 关注表
type Follow struct {
	ID        int `gorm:"type:int(11);primarykey;auto_increment"`
	CreatedAt int `gorm:"type:int(11);not null;default:0;comment:创建时间" json:"created_at"`
	UpdatedAt int `gorm:"type:int(11);not null;default:0;comment:更新时间" json:"updated_at"`
	DeletedAt int `gorm:"type:int(11);not null;default:0;comment:删除时间" json:"deleted_at"`
	FollowID  int `gorm:"type:int(11);not null;comment:关注人ID" json:"follow_id"`
	UserID    int `gorm:"type:int(11);not null;comment:用户ID" json:"user_id"`
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

type profileRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) error {
	ud := User{
		Email:      u.Email,
		Username:   u.Username,
		Bio:        u.Bio,
		Image:      u.Image,
		PasswdHash: u.PasswdHash,
	}
	rv := r.data.db.Create(&ud)
	u.UserID = int(ud.ID)
	return rv.Error
}

func (r *userRepo) VerifyUserExistByEmail(ctx context.Context, email string) bool {
	var count int64
	r.data.db.Model(&User{}).Where("email=?", email).Count(&count)
	if count > 0 {
		return true
	}
	return false
}
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	u := new(User)
	res := r.data.db.Where("email=?", email).First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("user", "not found by email")
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &biz.User{
		UserID:     int(u.ID),
		Email:      u.Email,
		Username:   u.Username,
		Bio:        u.Bio,
		Image:      u.Image,
		PasswdHash: u.PasswdHash,
	}, nil
}

func (r *userRepo) GetUserByUserID(ctx context.Context, userId int) (*biz.User, error) {
	u := new(User)
	res := r.data.db.Where("id=?", userId).First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("user", "not found by user id")
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &biz.User{
		UserID:     int(u.ID),
		Email:      u.Email,
		Username:   u.Username,
		Bio:        u.Bio,
		Image:      u.Image,
		PasswdHash: u.PasswdHash,
	}, nil
}

func (r *userRepo) GetUserByUserName(ctx context.Context, username string) (*biz.User, error) {
	u := new(User)
	res := r.data.db.Where("username=?", username).First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("user", "not found by username")
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &biz.User{
		UserID:     int(u.ID),
		Email:      u.Email,
		Username:   u.Username,
		Bio:        u.Bio,
		Image:      u.Image,
		PasswdHash: u.PasswdHash,
	}, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, userId int, bu *biz.User) (*biz.User, error) {
	u := &User{
		Username:   bu.Username,
		Email:      bu.Email,
		PasswdHash: bu.PasswdHash,
		Image:      bu.Image,
		Bio:        bu.Bio,
	}
	res := r.data.db.Model(&User{}).Where("id=?", userId).Updates(u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("user", "not found by email")
	}
	return &biz.User{
		UserID:     userId,
		Email:      u.Email,
		Username:   u.Username,
		Bio:        u.Bio,
		Image:      u.Image,
		PasswdHash: u.PasswdHash,
	}, nil
}

func NewProfileRepo(data *Data, logger log.Logger) biz.ProfileRepo {
	return &profileRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 获取关注用户
func (p *profileRepo) GetFollowByUserID(ctx context.Context, userId int) (*biz.Follow, error) {
	f := new(Follow)
	res := p.data.db.Where("follow_id=? and user_id=?", userId).First(&f)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("follow", "not found by username")
	}
	return &biz.Follow{}, nil
}

// 关注用户
func (p *profileRepo) FollowUser(ctx context.Context, myUserId, userId int) (bool, error) {
	f := new(Follow)
	res := p.data.db.Model(&Follow{}).Where("follow_id=? and user_id=?", userId, myUserId).First(&f)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		f.FollowID = userId
		f.UserID = myUserId
		tx := p.data.db.Create(&f)
		if tx.Error != nil {
			return false, tx.Error
		}
	}
	return true, nil
}

// 取消关注
func (p *profileRepo) UnfollowUser(ctx context.Context, myUserId, userId int) (bool, error) {
	res := p.data.db.Model(&Follow{}).Where("follow_id=? and user_id=?", userId, myUserId).Update("deleted_at", time.Now())
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, nil
}
