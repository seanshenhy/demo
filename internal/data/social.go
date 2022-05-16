package data

import (
	"context"
	"demo/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/clause"
)

// 文章的po
type Article struct {
	ID             int    `gorm:"type:int(11);primarykey;auto_increment"`
	CreatedAt      int    `gorm:"type:int(11);not null;default:0;comment:创建时间" json:"created_at"`
	UpdatedAt      int    `gorm:"type:int(11);not null;default:0;comment:更新时间" json:"updated_at"`
	DeletedAt      int    `gorm:"type:int(11);not null;default:0;comment:删除时间" json:"deleted_at"`
	Title          string `gorm:"type:varchar(64);not null;comment:文章标题" json:"title"`
	Description    string `gorm:"type:varchar(255);not null;comment:文章描述" json:"description"`
	Body           string `gorm:"type:varchar(511);not null;comment:文章体" json:"body"`
	FavoritesCount int    `gorm:"type:int(11);not null;default:0;comment:赞数量" json:"favorites_count"`
	UserID         string `gorm:"type:int(11);not null;comment:用户ID" json:"userId"`
}

// 评论po
type Comment struct {
	ID        int    `gorm:"type:int(11);primarykey;auto_increment"`
	CreatedAt int    `gorm:"type:int(11);not null;default:0;comment:创建时间" json:"created_at"`
	UpdatedAt int    `gorm:"type:int(11);not null;default:0;comment:更新时间" json:"updated_at"`
	DeletedAt int    `gorm:"type:int(11);not null;default:0;comment:删除时间" json:"deleted_at"`
	Body      string `gorm:"type:varchar(64);not null;comment:评论内容" json:"body"`
	ArticleID int    `gorm:"type:varchar(64);not null;comment:文章ID" json:"article_id"`
	UserID    int    `gorm:"type:int(11);not null;comment:用户ID" json:"user_id"`
}

// 标签po
type Tag struct {
	ID        int    `gorm:"type:int(11);primarykey;auto_increment"`
	CreatedAt int    `gorm:"type:int(11);not null;default:0;comment:创建时间" json:"created_at"`
	UpdatedAt int    `gorm:"type:int(11);not null;default:0;comment:更新时间" json:"updated_at"`
	DeletedAt int    `gorm:"type:int(11);not null;default:0;comment:删除时间" json:"deleted_at"`
	ArticleID int    `gorm:"type:int(11);not null;comment:文章ID" json:"article_id"`
	Tag       string `gorm:"type:varchar(64);not null; comment:标记" json:"tag"`
}

type articleRepo struct {
	data *Data
	log  *log.Helper
}

type commentRepo struct {
	data *Data
	log  *log.Helper
}

type tagRepo struct {
	data *Data
	log  *log.Helper
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (r *articleRepo) Create(ctx context.Context, do *biz.Article) (*biz.Article, error) {
	po := &Article{
		Title:       do.Title,
		Description: do.Description,
		Body:        do.Body,
	}
	rv := r.data.db.Create(po)
	do.ID = int(po.ID)
	return do, rv.Error
}

func (r *articleRepo) List(ctx context.Context, opt ...biz.ListOption) ([]*biz.Article, error) {
	pos := []Article{}
	rv := r.data.db.Find(&pos)
	if rv.Error != nil {
		return nil, rv.Error
	}
	dos := []*biz.Article{}
	for _, v := range pos {
		dos = append(dos, &biz.Article{
			Title:          v.Title,
			Description:    v.Description,
			Body:           v.Body,
			FavoritesCount: v.FavoritesCount,
		})
	}
	return dos, rv.Error
}

func (r *articleRepo) Get(ctx context.Context, articleId int) (*biz.Article, error) {
	po := new(Article)
	rv := r.data.db.First(&po)
	return &biz.Article{
		Title:          po.Title,
		Description:    po.Description,
		Body:           po.Body,
		FavoritesCount: po.FavoritesCount,
	}, rv.Error
}

func (r *articleRepo) Update(ctx context.Context, articleId int, do *biz.Article) (*biz.Article, error) {
	po := &Article{
		Title:          do.Title,
		Description:    do.Description,
		Body:           do.Body,
		FavoritesCount: do.FavoritesCount,
	}
	r.data.db.Model(&po).Where("id=?", articleId).Updates(po)
	return nil, nil
}

func (r *articleRepo) Delete(ctx context.Context, articleId int) error {
	var ars []Article
	tx := r.data.db.Clauses(clause.Returning{}).Where("id = ?", articleId).Delete(&ars)
	if len(ars) > 0 {
		return nil
	}
	return tx.Error
}

func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *commentRepo) Create(ctx context.Context, articleId int, do *biz.Comment) (*biz.Comment, error) {
	po := &Comment{
		Body:      do.Body,
		ArticleID: int(do.ArticleID),
	}
	tx := r.data.db.Create(po)
	return &biz.Comment{
		Body:      do.Body,
		ArticleID: do.ArticleID,
		Username:  do.Username,
	}, tx.Error
}

func (r *commentRepo) Get(ctx context.Context, articleId uint) (*biz.Comment, error) {
	po := &Comment{}
	tx := r.data.db.First(po, articleId)
	return &biz.Comment{
		Body:      po.Body,
		ArticleID: uint(po.ArticleID),
	}, tx.Error
}

func (r *commentRepo) List(ctx context.Context, articleId int) ([]*biz.Comment, error) {
	pos := []Comment{}
	tx := r.data.db.Find(&pos)
	if tx.Error != nil {
		return []*biz.Comment{}, nil
	}
	dos := []*biz.Comment{}
	for _, v := range pos {
		dos = append(dos, &biz.Comment{
			Body:      v.Body,
			ArticleID: uint(v.ArticleID),
		})
	}
	return dos, nil
}

func (r *commentRepo) Delete(ctx context.Context, id uint) error {
	return r.data.db.Delete(&Comment{}, id).Error
}

func NewTagRepo(data *Data, logger log.Logger) biz.TagRepo {
	return &tagRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *tagRepo) Create(ctx context.Context, ar *biz.Article) (*biz.Article, error) {
	tds := []Tag{}
	ttd := []string{}
	for _, v := range ar.TagList {
		tds = append(tds, Tag{ArticleID: int(ar.ID), Tag: v})
		ttd = append(ttd, v)
	}
	rv := r.data.db.Create(tds)
	ar.TagList = ttd
	return ar, rv.Error
}

func (r *tagRepo) Get(ctx context.Context, ar *biz.Article, arId int) (*biz.Article, error) {
	tds := []Tag{}
	rv := r.data.db.Where("article_id=?", arId).Find(&tds)
	ttd := []string{}
	for _, v := range tds {
		ttd = append(ttd, v.Tag)
	}
	ar.TagList = ttd
	return ar, rv.Error
}

func (r *tagRepo) Delete(ctx context.Context, arId int) error {
	tds := []Tag{}
	rv := r.data.db.Where("article_id=?", arId).Delete(&tds)
	if len(tds) > 0 {
		return nil
	}
	return rv.Error
}
