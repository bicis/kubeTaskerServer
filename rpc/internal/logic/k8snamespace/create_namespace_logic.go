package k8snamespace

import (
	"context"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubeTasker/kubeTaskerServer/rpc/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNamespaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNamespaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNamespaceLogic {
	return &CreateNamespaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateNamespaceLogic) CreateNamespace(in *core.CreateNamespaceReq) (*core.CreateNamespaceResp, error) {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.K8s.CoreV1().Namespaces().Create(l.ctx, in.Namespace, metav1.CreateOptions{})
	if err != nil {
		l.Logger.Error(errors.New("创建Namespace失败, " + err.Error()))
		return &core.CreateNamespaceResp{
			Msg: err.Error(),
		}, err
	}
	return &core.CreateNamespaceResp{
		Msg: "成功创建Namespace",
	}, nil
}
