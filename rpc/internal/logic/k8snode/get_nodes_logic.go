package k8snode

import (
	"context"
	"errors"
	"github.com/kubeTasker/kubeTaskerServer/rpc/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNodesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNodesLogic {
	return &GetNodesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNodesLogic) GetNodes(in *core.GetNodesReq) (*core.GetNodesResp, error) {
	// todo: add your logic here and delete this line
	//获取nodeList类型的node列表
	nodeList, err := l.svcCtx.K8s.CoreV1().Nodes().List(l.ctx, metav1.ListOptions{})
	if err != nil {
		l.Logger.Error(errors.New("获取Node列表失败, " + err.Error()))
		return &core.GetNodesResp{
			Msg: "获取Node列表失败, " + err.Error(),
			Data: &core.GetNodesData{
				Items: nil,
				Total: 0,
			},
		}, nil
	}
	//将nodeList中的node列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: toCells(nodeList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: in.FilterName},
			Paginate: &PaginateQuery{
				Limit: int(in.Limit),
				Page:  int(in.Page),
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	//将[]DataCell类型的node列表转为v1.node列表
	nodes := fromCells(data.GenericDataList)
	items := make([]*corev1.Node, 0)
	for i := range nodes {
		items = append(items, &nodes[i])
	}

	return &core.GetNodesResp{
		Msg: "成功！",
		Data: &core.GetNodesData{
			Items: items,
			Total: int64(total),
		},
	}, nil
}
