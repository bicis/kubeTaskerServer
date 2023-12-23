package k8sdeployment

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

type Deployment struct{}

// GetDeployments 获取deployment列表，支持过滤、排序、分页
func (d *Deployment) GetDeployments(l *GetDeploymentsLogic, in *core.GetDeploymentsReq) (resp *core.GetDeploymentsResp, err error) {
	//获取deploymentList类型的deployment列表
	deploymentList, err := l.svcCtx.K8s.AppsV1().Deployments(in.Namespace).List(l.ctx, metav1.ListOptions{})
	if err != nil {
		l.Error(errors.New("获取Deployment列表失败, " + err.Error()))
		return &core.GetDeploymentsResp{
			Meg:  err.Error(),
			Data: nil,
		}, errors.New("获取Deployment列表失败, " + err.Error())
	}
	//将deploymentList中的deployment列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
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
	dataSort := filtered.Sort().Paginate()
	// println("Deployment total: ", total)

	//将[]DataCell类型的deployment列表转为appsv1.deployment列表
	deployments := d.fromCells(dataSort.GenericDataList)
	items := make([]*v1.Deployment, 0)
	for i := range deployments {
		items = append(items, &deployments[i])
	}
	return &core.GetDeploymentsResp{
		Meg: "获取Deployment列表成功",
		Data: &core.GetDeploymentsData{
			Items: items,
			Total: int64(total),
		},
	}, nil
}

// 获取deployment详情
func (d *Deployment) GetDeploymentDetail(l *GetDeploymentDetailLogic, in *core.GetDeploymentDetailReq) (resp *core.GetDeploymentDetailResp, err error) {
	deployment, err := l.svcCtx.K8s.AppsV1().Deployments(in.Namespace).Get(l.ctx, in.DeploymentName, metav1.GetOptions{})
	if err != nil {
		l.Error(errors.New("获取Deployment详情失败, " + err.Error()))
		return &core.GetDeploymentDetailResp{
			Msg:  err.Error(),
			Data: nil,
		}, errors.New("获取Deployment详情失败, " + err.Error())
	}
	return &core.GetDeploymentDetailResp{
		Msg:  "获取Deployment详情成功",
		Data: deployment,
	}, nil
}

// 设置deployment副本数
func (d *Deployment) ScaleDeployment(l *ScaleDeploymentLogic, in *core.ScaleDeploymentReq) (resp *core.ScaleDeploymentResp, err error) {
	//获取autoscalingv1.Scale类型的对象，能点出当前的副本数
	scale, err := l.svcCtx.K8s.AppsV1().Deployments(in.Namespace).GetScale(l.ctx, in.DeploymentName, metav1.GetOptions{})
	if err != nil {
		l.Error(errors.New("获取Deployment副本数信息失败, " + err.Error()))
		return &core.ScaleDeploymentResp{
			Msg:  err.Error(),
			Data: "",
		}, errors.New("获取Deployment副本数信息失败, " + err.Error())
	}
	//修改副本数
	scale.Spec.Replicas = int32(in.ScaleNum)
	//更新副本数，传入scale对象
	newScale, err := l.svcCtx.K8s.AppsV1().Deployments(in.Namespace).UpdateScale(l.ctx, in.DeploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		l.Error(errors.New("更新Deployment副本数信息失败, " + err.Error()))
		return &core.ScaleDeploymentResp{
			Msg:  err.Error(),
			Data: "",
		}, errors.New("更新Deployment副本数信息失败, " + err.Error())
	}

	return &core.ScaleDeploymentResp{
		Msg:  "设置Deployment副本数成功",
		Data: fmt.Sprintf("最新副本数: %d", newScale.Spec.Replicas),
	}, nil
}

// 创建deployment,接收DeployCreate对象
func (d *Deployment) CreateDeployment(l *CreateDeploymentLogic, in *core.CreateDeploymentReq) (resp *core.CreateDeploymentResp, err error) {
	//将data中的属性组装成appsv1.Deployment对象,并将入参数据放入
	deployment := &v1.Deployment{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      in.Name,
			Namespace: in.Namespace,
			Labels:    in.Label,
		},
		//Spec中定义副本数、选择器、以及pod属性
		Spec: v1.DeploymentSpec{
			Replicas: &in.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: in.Label,
			},
			Template: corev1.PodTemplateSpec{
				//定义pod名和标签
				ObjectMeta: metav1.ObjectMeta{
					Name:   in.Name,
					Labels: in.Label,
				},
				//定义容器名、镜像和端口
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  in.Name,
							Image: in.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: in.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		//Status定义资源的运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status: v1.DeploymentStatus{},
	}
	//判断是否打开健康检查功能，若打开，则定义ReadinessProbe和LivenessProbe
	if in.HealthCheck {
		//设置第一个容器的ReadinessProbe，因为我们pod中只有一个容器，所以直接使用index 0即可
		//若pod中有多个容器，则这里需要使用for循环去定义了
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: in.HealthPath,
					//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
					//Type=0则表示表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
					//Type=1则表示表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: in.ContainerPort,
					},
				},
			},
			//初始化等待时间
			InitialDelaySeconds: 5,
			//超时时间
			TimeoutSeconds: 5,
			//执行间隔
			PeriodSeconds: 5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: in.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: in.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		//定义容器的limit和request资源
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(in.Cpu),
			corev1.ResourceMemory: resource.MustParse(in.Memory),
		}
		deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(in.Cpu),
			corev1.ResourceMemory: resource.MustParse(in.Memory),
		}
	}
	//调用sdk创建deployment
	_, err = l.svcCtx.K8s.AppsV1().Deployments(in.Namespace).
		Create(l.ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		l.Error(errors.New("创建Deployment失败, " + err.Error()))
		return &core.CreateDeploymentResp{
			Msg: err.Error(),
		}, errors.New("创建Deployment失败, " + err.Error())
	}
	return &core.CreateDeploymentResp{
		Msg: "创建Deployment成功",
	}, nil
}

// 删除deployment
func (d *Deployment) DeleteDeployment(l *DeleteDeploymentLogic, in *core.DeleteDeploymentReq) (resp *core.DeleteDeploymentResp, err error) {
	err = l.svcCtx.K8s.AppsV1().Deployments(in.Namespace).Delete(l.ctx, in.DeploymentName, metav1.DeleteOptions{})
	if err != nil {
		l.Error(errors.New("删除Deployment失败, " + err.Error()))
		return &core.DeleteDeploymentResp{
			Msg: err.Error(),
		}, errors.New("删除Deployment失败, " + err.Error())
	}

	return &core.DeleteDeploymentResp{
		Msg: "删除Deployment成功",
	}, nil
}

// 重启deployment
func (d *Deployment) RestartDeployment(l *RestartDeploymentLogic, in *core.RestartDeploymentReq) (resp *core.RestartDeploymentResp, err error) {
	//此功能等同于一下kubectl命令
	//kubectl deployment ${service} -p \
	//'{"spec":{"template":{"spec":{"containers":[{"name":"'"${service}"'","env":[{"name":"RESTART_","value":"'$(date +%s)'"}]}]}}}}'

	//使用patchData Map组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": in.DeploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	//序列化为字节，因为patch方法只接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		l.Error(errors.New("json序列化失败, " + err.Error()))
		return &core.RestartDeploymentResp{
			Msg: err.Error(),
		}, errors.New("json序列化失败, " + err.Error())
	}
	//调用patch方法更新deployment
	_, err = l.svcCtx.K8s.AppsV1().Deployments(in.Namespace).Patch(l.ctx, in.DeploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		l.Error(errors.New("重启Deployment失败, " + err.Error()))
		return &core.RestartDeploymentResp{
			Msg: err.Error(),
		}, errors.New("重启Deployment失败, " + err.Error())
	}

	return &core.RestartDeploymentResp{
		Msg: "重启Deployment成功",
	}, nil
}

// 更新deployment
func (d *Deployment) UpdateDeployment(l *UpdateDeploymentLogic, in *core.UpdateDeploymentReq) (resp *core.UpdateDeploymentResp, err error) {
	var deploy = &v1.Deployment{}

	err = json.Unmarshal([]byte(in.Content), deploy)
	if err != nil {
		l.Error(errors.New("反序列化失败, " + err.Error()))
		return &core.UpdateDeploymentResp{
			Msg: err.Error(),
		}, errors.New("反序列化失败, " + err.Error())
	}

	_, err = l.svcCtx.K8s.AppsV1().Deployments(in.Namespace).Update(l.ctx, deploy, metav1.UpdateOptions{})
	if err != nil {
		l.Error(errors.New("更新Deployment失败, " + err.Error()))
		return &core.UpdateDeploymentResp{
			Msg: err.Error(),
		}, errors.New("更新Deployment失败, " + err.Error())
	}
	return &core.UpdateDeploymentResp{
		Msg: "更新Deployment成功",
	}, nil
}

// 获取每个namespace的deployment数量
func (d *Deployment) GetDeployNumPerNp(l *GetDeployNumPerNpLogic, in *core.GetDeployNumPerNpReq) (resp *core.GetDeployNumPerNpResp, err error) {
	namespaceList, err := l.svcCtx.K8s.CoreV1().Namespaces().List(l.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	deploysNps := make([]*core.GetDeployNumPerNpData, 0)
	for _, namespace := range namespaceList.Items {
		deploymentList, err := l.svcCtx.K8s.AppsV1().Deployments(namespace.Name).List(l.ctx, metav1.ListOptions{})
		if err != nil {
			return &core.GetDeployNumPerNpResp{
				Msg:  err.Error(),
				Data: nil,
			}, err
		}

		deploysNp := &core.GetDeployNumPerNpData{
			Namespace:     namespace.Name,
			DeploymentNum: int64(len(deploymentList.Items)),
		}

		deploysNps = append(deploysNps, deploysNp)
	}
	return &core.GetDeployNumPerNpResp{
		Msg:  "获取每个namespace的deployment数量成功",
		Data: deploysNps,
	}, nil
}

func (d *Deployment) toCells(std []v1.Deployment) []DataCell {
	cells := make([]DataCell, 0)
	for _, v := range std {
		cells = append(cells, deploymentCell(v))
	}
	return cells
}

func (d *Deployment) fromCells(cells []DataCell) []v1.Deployment {
	deployments := make([]v1.Deployment, 0)
	for _, v := range cells {
		deployments = append(deployments, v1.Deployment(v.(deploymentCell)))
	}
	return deployments
}
