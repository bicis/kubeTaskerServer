package k8srolebinding

import (
	"context"
	"github.com/kubeTasker/kubeTaskerServer/rpc/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleBindingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleBindingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleBindingLogic {
	return &CreateRoleBindingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleBindingLogic) CreateRoleBinding(in *core.CreateRoleBindingReq) (*core.CreateRoleBindingResp, error) {
	// todo: add your logic here and delete this line
	roleBinding := &v1.RoleBinding{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default-admin",
		},
		Subjects: []v1.Subject{v1.Subject{
			Kind:      "ServiceAccount",
			APIGroup:  "",
			Name:      "default",
			Namespace: in.Namespace,
		}},
		RoleRef: v1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "admin",
		},
	}
	_, err := l.svcCtx.K8s.RbacV1().RoleBindings(in.Namespace).Create(l.ctx, roleBinding, metav1.CreateOptions{})
	if err != nil {
		return &core.CreateRoleBindingResp{
			Msg: "创建rolebinding失败:" + err.Error(),
		}, nil
	}
	return &core.CreateRoleBindingResp{
		Msg: "创建rolebinding成功",
	}, nil
}
