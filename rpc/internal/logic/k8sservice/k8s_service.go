package k8sservice

import (
	"encoding/json"
	"errors"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Service struct{}

// 获取service列表，支持过滤、排序、分页
func (s *Service) GetServices(l *GetServicesLogic, in *core.GetServicesReq) (resp *core.GetServicesResp, err error) {
	//获取serviceList类型的service列表
	serviceList, err := l.svcCtx.K8s.CoreV1().Services(in.Namespace).List(l.ctx, metav1.ListOptions{})
	if err != nil {
		l.Error(errors.New("获取Service列表失败, " + err.Error()))
		return &core.GetServicesResp{
			Msg:  err.Error(),
			Data: nil,
		}, errors.New("获取Service列表失败, " + err.Error())
	}
	//将serviceList中的service列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: s.toCells(serviceList.Items),
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

	//将[]DataCell类型的service列表转为v1.service列表
	services := s.fromCells(data.GenericDataList)
	items := make([]*v1.Service, 0)
	for i := range services {
		items = append(items, &services[i])
	}
	return &core.GetServicesResp{
		Msg: "获取Servicet列表成功",
		Data: &core.GetServicesData{
			Items: items,
			Total: int64(total),
		},
	}, nil
}

// 获取service详情
func (s *Service) GetServiceDetail(l *GetServiceDetailLogic, in *core.GetServiceDetailReq) (resp *core.GetServiceDetailResp, err error) {
	service, err := l.svcCtx.K8s.CoreV1().Services(in.Namespace).Get(l.ctx, in.ServiceName, metav1.GetOptions{})
	if err != nil {
		l.Error(errors.New("获取Service详情失败, " + err.Error()))
		return &core.GetServiceDetailResp{
			Msg:  err.Error(),
			Data: nil,
		}, errors.New("获取Service详情失败, " + err.Error())
	}

	return &core.GetServiceDetailResp{
		Msg:  "获取Service详情成功",
		Data: service,
	}, nil
}

// 创建service,,接收ServiceCreate对象
func (s *Service) CreateService(l *CreateServiceLogic, in *core.CreateServiceReq) (resp *core.CreateServiceResp, err error) {
	//将data中的数据组装成corev1.Service对象
	service := &v1.Service{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      in.Name,
			Namespace: in.Namespace,
			Labels:    in.Label,
		},
		//Spec中定义类型，端口，选择器
		Spec: v1.ServiceSpec{
			Type: v1.ServiceType(in.Type),
			Ports: []v1.ServicePort{
				{
					Name:     "http",
					Port:     int32(in.Port),
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: int32(in.ContainerPort),
					},
				},
			},
			Selector: in.Label,
		},
	}
	//默认ClusterIP,这里是判断NodePort,添加配置
	if in.NodePort != 0 && in.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = int32(in.NodePort)
	}
	//创建Service
	_, err = l.svcCtx.K8s.CoreV1().Services(in.Namespace).Create(l.ctx, service, metav1.CreateOptions{})
	if err != nil {
		l.Error(errors.New("创建Service失败, " + err.Error()))
		return &core.CreateServiceResp{
			Msg: err.Error(),
		}, errors.New("创建Service失败, " + err.Error())
	}

	return &core.CreateServiceResp{
		Msg: "创建Service成功",
	}, nil
}

// 删除service
func (s *Service) DeleteService(l *DeleteServiceLogic, in *core.DeleteServiceReq) (resp *core.DeleteServiceResp, err error) {
	err = l.svcCtx.K8s.CoreV1().Services(in.Namespace).Delete(l.ctx, in.ServiceName, metav1.DeleteOptions{})
	if err != nil {
		l.Error(errors.New("删除Service失败, " + err.Error()))
		return &core.DeleteServiceResp{
			Msg: err.Error(),
		}, errors.New("删除Service失败, " + err.Error())
	}

	return &core.DeleteServiceResp{
		Msg: "删除Service成功",
	}, nil
}

// 更新service
func (s *Service) UpdateService(l *UpdateServiceLogic, in *core.UpdateServiceReq) (resp *core.UpdateServiceResp, err error) {
	var service = &v1.Service{}

	err = json.Unmarshal([]byte(in.Content), service)
	if err != nil {
		l.Error(errors.New("反序列化失败, " + err.Error()))
		return &core.UpdateServiceResp{
			Msg: err.Error(),
		}, errors.New("反序列化失败, " + err.Error())
	}

	_, err = l.svcCtx.K8s.CoreV1().Services(in.Namespace).Update(l.ctx, service, metav1.UpdateOptions{})
	if err != nil {
		l.Error(errors.New("更新service失败, " + err.Error()))
		return &core.UpdateServiceResp{
			Msg: err.Error(),
		}, errors.New("更新service失败, " + err.Error())
	}
	return &core.UpdateServiceResp{
		Msg: "更新Service成功",
	}, nil
}

func (s *Service) toCells(std []v1.Service) []DataCell {
	cells := make([]DataCell, 0)
	for _, v := range std {
		cells = append(cells, serviceCell(v))
	}
	return cells
}

func (s *Service) fromCells(cells []DataCell) []v1.Service {
	services := make([]v1.Service, 0)
	for _, v := range cells {
		services = append(services, v1.Service(v.(serviceCell)))
	}
	return services
}
