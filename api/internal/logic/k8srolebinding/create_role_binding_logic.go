package k8srolebinding

import (
	"context"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"

	"github.com/kubeTasker/kubeTaskerServer/api/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleBindingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleBindingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleBindingLogic {
	return &CreateRoleBindingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

func (l *CreateRoleBindingLogic) CreateRoleBinding(req *types.CreateRoleBindingReq) (resp *types.CreateRoleBindingResp, err error) {
	// todo: add your logic here and delete this line
	result, err := l.svcCtx.CoreRpc.CreateRoleBinding(l.ctx, &core.CreateRoleBindingReq{
		Namespace: req.Namespace,
	})
	if err != nil {
		return &types.CreateRoleBindingResp{
			Msg: err.Error(),
		}, nil
	}
	return &types.CreateRoleBindingResp{
		Msg: result.Msg,
	}, nil
}
