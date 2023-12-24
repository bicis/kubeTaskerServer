package k8srolebinding

import (
	"context"
	"github.com/kubeTasker/kubeTaskerServer/rpc/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"

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
	//RoleBinding := &v1.RoleBinding{
	//	TypeMeta: metav1.TypeMeta{},
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name: "default-admin",
	//	},
	//	Subjects: nil,
	//	RoleRef:  v1.RoleRef{},
	//}
	//l.svcCtx.K8s.RbacV1().RoleBindings(in.Namespace).Create()
	return &core.CreateRoleBindingResp{}, nil
}
