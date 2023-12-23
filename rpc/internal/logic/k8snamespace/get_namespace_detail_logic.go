package k8snamespace

import (
	"context"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeTasker/kubeTaskerServer/rpc/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNamespaceDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNamespaceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNamespaceDetailLogic {
	return &GetNamespaceDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNamespaceDetailLogic) GetNamespaceDetail(in *core.GetNamespaceDetailReq) (*core.GetNamespaceDetailResp, error) {
	// todo: add your logic here and delete this line
	namespace, err := l.svcCtx.K8s.CoreV1().Namespaces().Get(l.ctx, in.NamespaceName, metav1.GetOptions{})
	if err != nil {
		l.Logger.Error(errors.New("获取Namespace详情失败, " + err.Error()))
		return &core.GetNamespaceDetailResp{
			Msg: "获取Namespace详情失败, " + err.Error(),
			Data: &core.GetNamespaceDetailData{
				Namespace: nil,
			},
		}, nil
	}
	return &core.GetNamespaceDetailResp{
		Msg: "成功！",
		Data: &core.GetNamespaceDetailData{
			Namespace: namespace,
		},
	}, nil
}
