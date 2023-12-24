package user

import (
	"context"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/logic/k8snamespace"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/types"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.UUIDsReq) (resp *types.BaseMsgResp, err error) {
	for i := range req.Ids {
		lu := NewGetUserByIdLogic(l.ctx, l.svcCtx)
		respLu, err := lu.GetUserById(&types.UUIDReq{
			Id: req.Ids[i],
		})
		if err != nil {
			continue
		}
		ln := k8snamespace.NewDeleteNamespaceLogic(l.ctx, l.svcCtx)
		_, _ = ln.DeleteNamespace(&types.DeleteNamespaceReq{NamespaceName: *respLu.Data.Username})
	}
	result, err := l.svcCtx.CoreRpc.DeleteUser(l.ctx, &core.UUIDsReq{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, result.Msg)}, nil
}
