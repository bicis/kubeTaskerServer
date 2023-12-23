package k8snamespace

import (
	"context"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeTasker/kubeTaskerServer/rpc/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNamespaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNamespaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNamespaceLogic {
	return &DeleteNamespaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNamespaceLogic) DeleteNamespace(in *core.DeleteNamespaceReq) (*core.DeleteNamespaceResp, error) {
	// todo: add your logic here and delete this line
	err := l.svcCtx.K8s.CoreV1().Namespaces().Delete(l.ctx, in.NamespaceName, metav1.DeleteOptions{})
	if err != nil {
		l.Logger.Error(errors.New("获取Namespace详情失败, " + err.Error()))
		return &core.DeleteNamespaceResp{
			Msg:  "获取Namespace详情失败," + err.Error(),
			Data: "",
		}, nil
	}
	return &core.DeleteNamespaceResp{
		Msg:  "成功！",
		Data: "成功删除namespace:" + in.NamespaceName,
	}, nil
}
