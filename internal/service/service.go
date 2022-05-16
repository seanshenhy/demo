package service

import (
	v1 "demo/api/realworld/v1"
	"demo/internal/biz"

	"github.com/google/wire"

	"github.com/go-kratos/kratos/v2/log"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewRealworldService)

type RealworldService struct {
	v1.UnimplementedRealworldServer

	uc  *biz.UserUsecase
	sc  *biz.SocialUsecase
	log *log.Helper
}

func NewRealworldService(uc *biz.UserUsecase, sc *biz.SocialUsecase, logger log.Logger) *RealworldService {
	return &RealworldService{uc: uc, sc: sc, log: log.NewHelper(logger)}
}
