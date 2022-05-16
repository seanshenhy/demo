package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Article struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	Slug           string    `json:"slug"`
	TagList        []string  `json:"tagList"`
	Username       string    `json:"username"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         Author    `json:"author"`
}

type Comment struct {
	ID        uint
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ArticleID uint
	Username  string
	Article   *Article
	Author    *Author
}

type ArticleRepo interface {
	Create(ctx context.Context, ar *Article) (*Article, error)
	List(ctx context.Context, opt ...ListOption) ([]*Article, error)
	Get(ctx context.Context, articleId int) (*Article, error)
	Update(ctx context.Context, articleId int, ar *Article) (*Article, error)
	Delete(ctx context.Context, articleId int) error
}

type CommentRepo interface {
	Create(ctx context.Context, articleId int, c *Comment) (*Comment, error)
	Get(ctx context.Context, commentId uint) (*Comment, error)
	List(ctx context.Context, articleId int) ([]*Comment, error)
	Delete(ctx context.Context, id uint) error
}

type TagRepo interface {
	Create(ctx context.Context, ar *Article) (*Article, error)
	Get(ctx context.Context, ar *Article, arId int) (*Article, error)
	Delete(ctx context.Context, arId int) error
}

type SocialUsecase struct {
	ar  ArticleRepo
	cr  CommentRepo
	tr  TagRepo
	ur  UserRepo
	log *log.Helper
}

func NewSocialUseCase(ar ArticleRepo, cr CommentRepo, tr TagRepo, logger log.Logger) *SocialUsecase {
	return &SocialUsecase{ar: ar, cr: cr, tr: tr, log: log.NewHelper(logger)}
}

func (s *SocialUsecase) CreateArticle(ctx context.Context, ar *Article) (*Article, error) {
	// 创建文章
	arr, err := s.ar.Create(ctx, ar)
	// 创建tag
	if len(arr.TagList) > 0 {
		arr, err = s.tr.Create(ctx, arr)
		if err != nil {
			return arr, err
		}
	}
	return arr, nil
}

func (s *SocialUsecase) ListArticles(ctx context.Context, opt ...ListOption) (rv []*Article, err error) {
	rv, err = s.ar.List(ctx, opt...)
	return rv, err
}

func (s *SocialUsecase) FeedArticles(ctx context.Context, opt ...ListOption) (rv []*Article, err error) {
	rv, err = s.ar.List(ctx, opt...)
	return rv, err
}

func (s *SocialUsecase) GetArticle(ctx context.Context, articleId int) (do *Article, err error) {
	do, err = s.ar.Get(ctx, articleId)
	if err != nil {
		return nil, err
	}
	do, err = s.tr.Get(ctx, do, do.ID)
	return
}

func (s *SocialUsecase) UpdateArticle(ctx context.Context, articleId int, do *Article) (*Article, error) {
	do, err := s.ar.Update(ctx, articleId, do)
	if err != nil {
		return nil, err
	}
	return s.tr.Get(ctx, do, do.ID)
}

func (s *SocialUsecase) DeleteArticle(ctx context.Context, articleId int) (*Article, error) {
	ar, err := s.ar.Get(ctx, articleId)
	if err != nil {
		return nil, err
	}
	if err := s.ar.Delete(ctx, articleId); err != nil {
		return ar, err
	}
	if err := s.tr.Delete(ctx, ar.ID); err != nil {
		return ar, err
	}
	return ar, nil
}
