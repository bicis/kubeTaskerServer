package k8snode

import (
	"context"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeTasker/kubeTaskerServer/rpc/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNodeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNodeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNodeDetailLogic {
	return &GetNodeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNodeDetailLogic) GetNodeDetail(in *core.GetNodeDetailReq) (*core.GetNodeDetailResp, error) {
	node, err := l.svcCtx.K8s.CoreV1().Nodes().Get(l.ctx, in.NodeName, metav1.GetOptions{})
	if err != nil {
		l.Logger.Error(errors.New("获取Node详情失败, " + err.Error()))
		return &core.GetNodeDetailResp{
			Msg: "获取Node详情失败, " + err.Error(),
			Data: &core.GetNodeDetailData{
				Node: nil,
			},
		}, nil
	}

	return &core.GetNodeDetailResp{
		Msg: "成功！",
		Data: &core.GetNodeDetailData{
			Node: node,
		},
	}, nil
}
