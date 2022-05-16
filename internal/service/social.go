package service

import (
	"context"
	v1 "demo/api/realworld/v1"
	"demo/internal/biz"
)

// 创建文章
func (s *RealworldService) CreateArticle(ctx context.Context, req *v1.CreateArticleRequest) (*v1.SingleArticlesReply, error) {
	ar, err := s.sc.CreateArticle(ctx, &biz.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		TagList:     req.Article.TagList,
		Username:    ctx.Value("loginUser").(string),
	})
	if err != nil {
		return nil, err
	}

	return &v1.SingleArticlesReply{
		Article: formatArticleReply(ar),
	}, nil
}

// 列表文章
func (s *RealworldService) ListArticles(ctx context.Context, req *v1.ListArticlesRequest) (*v1.MultipleArticlesReply, error) {
	filter := make(map[string]string)
	if req.Tag != "" {
		filter["tag"] = req.Tag
	}
	if req.Author != "" {
		filter["author"] = req.Author
	}
	if req.Favorited != "" {
		filter["favorited"] = req.Favorited
	}
	ars, err := s.sc.ListArticles(ctx, biz.ListLimit(req.Limit), biz.ListLimit(req.Offset), biz.ListFilter(filter))
	if err != nil {
		return nil, err
	}
	articles := []*v1.Article{}
	for _, v := range ars {
		articles = append(articles, formatArticleReply(v))
	}
	return &v1.MultipleArticlesReply{
		Articles: articles,
	}, nil
}

// 文章列表
func (s *RealworldService) FeedArticles(ctx context.Context, req *v1.FeedArticlesRequest) (*v1.MultipleArticlesReply, error) {
	ars, err := s.sc.ListArticles(ctx, biz.ListLimit(req.Limit), biz.ListLimit(req.Offset))
	if err != nil {
		return nil, err
	}
	articles := []*v1.Article{}
	for _, v := range ars {
		articles = append(articles, formatArticleReply(v))
	}
	return &v1.MultipleArticlesReply{
		Articles: articles,
	}, nil

}

// 格式化文章
func formatArticleReply(ar *biz.Article) *v1.Article {
	return &v1.Article{
		Title:          ar.Title,
		Body:           ar.Body,
		Description:    ar.Description,
		CreatedAt:      ar.CreatedAt.String(),
		UpdatedAt:      ar.UpdatedAt.String(),
		Favorited:      ar.Favorited,
		FavoritesCount: uint32(ar.FavoritesCount),
		Author: &v1.Author{
			Bio:       ar.Author.Bio,
			Username:  ar.Author.Username,
			Image:     ar.Author.Image,
			Following: ar.Author.Following,
		},
	}
}

// 获取文章详情
func (s *RealworldService) GetArticle(ctx context.Context, req *v1.GetArticleRequest) (*v1.SingleArticlesReply, error) {
	do, err := s.sc.GetArticle(ctx, int(req.ArticleId))
	if err != nil {
		return nil, err
	}
	return &v1.SingleArticlesReply{
		Article: formatArticleReply(do),
	}, nil
}

// 更新文章详情
func (s *RealworldService) UpdateArticle(ctx context.Context, req *v1.UpdateArticleRequest) (*v1.SingleArticlesReply, error) {
	do, err := s.sc.UpdateArticle(ctx, int(req.ArticleId), &biz.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body})
	if err != nil {
		return nil, err
	}
	return &v1.SingleArticlesReply{
		Article: formatArticleReply(do),
	}, nil

}

// 删除文章
func (s *RealworldService) DeleteArticle(ctx context.Context, req *v1.DeleteArticleRequest) (*v1.SingleArticlesReply, error) {
	do, err := s.sc.DeleteArticle(ctx, int(req.ArticleId))
	if err != nil {
		return nil, err
	}
	return &v1.SingleArticlesReply{
		Article: formatArticleReply(do),
	}, nil
}

// func (s *RealworldService) FavoriteArticle(ctx context.Context, req *v1.FavoriteArticleRequest) (*v1.SingleArticlesReply, error) {

// 	return &v1.SingleArticlesReply{}, nil

// }

// func (s *RealworldService) UnfavoriteArticle(ctx context.Context, req *v1.UnfavoriteArticleRequest) (*v1.SingleArticlesReply, error) {

// 	return &v1.SingleArticlesReply{}, nil

// }

// 添加评论
func (s *RealworldService) AddComments(ctx context.Context, req *v1.AddCommentsRequest) (*v1.SingleCommentReply, error) {

	return &v1.SingleCommentReply{}, nil

}

// func (s *RealworldService) GetComments(ctx context.Context, req *v1.GetCommentsRequest) (*v1.MultipleCommentsReply, error) {

// 	return &v1.MultipleCommentsReply{}, nil

// }

// func (s *RealworldService) DeleteComment(ctx context.Context, req *v1.DeleteCommentRequest) (*v1.SingleCommentReply, error) {

// 	return &v1.SingleCommentReply{}, nil

// }

// func (s *RealworldService) GetTags(ctx context.Context, req *v1.GetTagsRequest) (*v1.ListTagsReply, error) {

// 	return &v1.ListTagsReply{}, nil

// }
