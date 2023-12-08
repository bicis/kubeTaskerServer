package workflow

import (
	"context"

	"github.com/kubeTasker/kubeTaskerServer/rpc/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkflowLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWorkflowLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkflowLogsLogic {
	return &WorkflowLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WorkflowLogsLogic) WorkflowLogs(in *core.WorkflowLogRequest, stream core.Core_WorkflowLogsServer) error {
	// todo: add your logic here and delete this line

	// Stream must supplement other logic by itself, and the logic is incomplete.
	// /root/TeamProject/kubeTaskerServer/rpc/types/core/core_grpc.pb.go

	return nil
}
