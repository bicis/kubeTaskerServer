package k8spod

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct{}

// 获取pod列表,支持过滤,排序,分页
func (p *Pod) GetPods(l *GetPodsLogic, in *core.GetPodsReq) (getPodsResp *core.GetPodsResp, err error) {
	podList, err := l.svcCtx.K8s.CoreV1().Pods(in.Namespace).List(l.ctx, metav1.ListOptions{})
	if err != nil {
		l.Error("获取Pod列表失败," + err.Error())
		// 返回给上一层,最终返回给前端,前端捕获到后打印出来
		return &core.GetPodsResp{
			Msg:  err.Error(),
			Data: nil,
		}, errors.New("获取Pod列表失败," + err.Error())
	}
	// 实例化dataselector结构体,组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: in.FilterName},
			Paginate: &PaginateQuery{
				Limit: int(in.Limit),
				Page:  int(in.Page),
			},
		},
	}
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	// 排序和分页
	data := filtered.Sort().Paginate()
	//println("Pod total: ", total)
	// 将DataCell类型转成Pod
	pods := p.fromCells(data.GenericDataList)
	items := make([]*v1.Pod, 0)
	for i := range pods {
		items = append(items, &pods[i])
	}
	return &core.GetPodsResp{
		Msg: "获取Pod列表成功",
		Data: &core.GetPodsData{
			Items: items,
			Total: int64(total),
		},
	}, nil
}

// 获取pod详情
func (p *Pod) GetPodDetail(l *GetPodDetailLogic, in *core.GetPodDetailReq) (getPodDetailResp *core.GetPodDetailResp, err error) {
	pod, err := l.svcCtx.K8s.CoreV1().Pods(in.Namespace).Get(l.ctx, in.PodName, metav1.GetOptions{})
	if err != nil {
		l.Error("获取Pod详情失败," + err.Error())
		return &core.GetPodDetailResp{
			Msg:  err.Error(),
			Data: nil,
		}, errors.New("获取Pod详情失败," + err.Error())
	}

	return &core.GetPodDetailResp{
		Msg:  "获取Pod列表成功",
		Data: pod,
	}, nil
}

// 删除Pod
func (p *Pod) DeletePod(l *DeletePodLogic, in *core.DeletePodReq) (deletePodResp *core.DeletePodResp, err error) {
	err = l.svcCtx.K8s.CoreV1().Pods(in.Namespace).Delete(l.ctx, in.PodName, metav1.DeleteOptions{})
	if err != nil {
		l.Error("删除Pod详情失败," + err.Error())
		return &core.DeletePodResp{
			Msg: err.Error(),
		}, errors.New("删除Pod详情失败," + err.Error())
	}
	return &core.DeletePodResp{
		Msg: "删除Pod成功",
	}, nil
}

// 更新Pod
func (p *Pod) UpdatePod(l *UpdatePodLogic, in *core.UpdatePodReq) (updatePodResp *core.UpdatePodResp, err error) {
	var pod = &v1.Pod{}
	// 反序列化为Pod对象
	err = json.Unmarshal([]byte(in.Content), pod)
	if err != nil {
		l.Error("反序列化失败," + err.Error())
		return &core.UpdatePodResp{
			Msg: err.Error(),
		}, errors.New("反序列化失败," + err.Error())
	}
	// 更新pod
	_, err = l.svcCtx.K8s.CoreV1().Pods(in.Namespace).Update(l.ctx, pod, metav1.UpdateOptions{})
	if err != nil {
		l.Error("更新Pod失败," + err.Error())
		return &core.UpdatePodResp{
			Msg: err.Error(),
		}, errors.New("更新Pod失败," + err.Error())
	}
	return &core.UpdatePodResp{
		Msg: "更新Pod成功",
	}, nil
}

// 获取Pod中的容器名
func (p *Pod) GetPodContainer(l *GetPodContainerLogic, in *core.GetPodContainerReq) (getPodContainerResp *core.GetPodContainerResp, err error) {
	getPodDetailLogic := NewGetPodDetailLogic(l.ctx, l.svcCtx)
	pod, err := getPodDetailLogic.GetPodDetail(&core.GetPodDetailReq{
		PodName:   in.PodName,
		Namespace: in.Namespace,
	})
	if err != nil {
		return &core.GetPodContainerResp{
			Msg:  err.Error(),
			Data: nil,
		}, err
	}
	containers := make([]string, 0)
	for _, container := range pod.Data.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return &core.GetPodContainerResp{
		Msg:  "获取Pod容器成功",
		Data: containers,
	}, nil
}

// 获取Pod内容器日志
func (p *Pod) GetPodLog(l *GetPodLogLogic, in *core.GetPodLogReq) (getPodLogResp *core.GetPodLogResp, err error) {
	//设置日志配置,容器名,获取内容的配置
	lineLimit := int64(podLogTailLine)
	option := &v1.PodLogOptions{
		Container: in.ContainerName,
		TailLines: &lineLimit,
	}
	// 获取一个request实例
	req := l.svcCtx.K8s.CoreV1().Pods(in.Namespace).GetLogs(in.PodName, option)
	// 发起stream连接,获取到Response.body
	podLogs, err := req.Stream(l.ctx)
	if err != nil {
		l.Error("更新Pod失败," + err.Error())
		return &core.GetPodLogResp{
			Meg: err.Error(),
		}, errors.New("更新Pod失败," + err.Error())
	}
	defer func(podLogs io.ReadCloser) {
		err := podLogs.Close()
		if err != nil {
			logx.Error("Failed to close the pod log stream")
		}
	}(podLogs)
	// 将Response.body 写入到缓存区,目的为了转换成string类型
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		l.Error("复制podLog失败," + err.Error())
		return &core.GetPodLogResp{
			Meg: err.Error(),
		}, errors.New("复制podLog失败," + err.Error())
	}
	return &core.GetPodLogResp{
		Meg:  "获取Pod中容器日志成功",
		Data: buf.String(),
	}, nil
}

// 获取每个namespace中pod的数量
func (p *Pod) GetPodNumPerNp(l *GetPodNumPerNpLogic, in *core.GetPodNumPerNpReq) (getPodNumPerNpResp *core.GetPodNumPerNpResp, err error) {
	namespaceList, err := l.svcCtx.K8s.CoreV1().Namespaces().List(l.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	podsNps := make([]*core.GetPodNumPerNpData, 0)
	for _, namespace := range namespaceList.Items {
		//获取pod列表
		podList, err := l.svcCtx.K8s.CoreV1().Pods(namespace.Name).List(l.ctx, metav1.ListOptions{})
		if err != nil {
			return &core.GetPodNumPerNpResp{
				Msg:  err.Error(),
				Data: nil,
			}, err
		}
		//组装数据
		podsNp := &core.GetPodNumPerNpData{
			Namespace: namespace.Name,
			PodNum:    int64(len(podList.Items)),
		}
		//添加到podsNps数组中
		podsNps = append(podsNps, podsNp)
	}
	return &core.GetPodNumPerNpResp{
		Msg:  "获取每个namespace的pod数量成功",
		Data: podsNps,
	}, nil
}

// 类型转换方法corev1.Pod --> DataCell,DataCell-->corev1.Pod
func (p *Pod) toCells(pods []v1.Pod) []DataCell {
	cells := make([]DataCell, 0)
	for _, v := range pods {
		cells = append(cells, podCell(v))
	}
	return cells
}

func (p *Pod) fromCells(cells []DataCell) []v1.Pod {
	pods := make([]v1.Pod, 0)
	for _, v := range cells {
		// cells[i].(podCell)是将DataCell类型转换成podCell
		pods = append(pods, v1.Pod(v.(podCell)))
	}
	return pods
}
