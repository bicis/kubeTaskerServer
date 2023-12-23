package k8sconfigmap

import (
	"github.com/kubeTasker/kubeTaskerServer/rpc/types/core"
	"k8s.io/api/core/v1"
)

import (
	"encoding/json"
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMap struct{}

type ConfigMapsResp struct {
	Items []corev1.ConfigMap `json:"items"`
	Total int                `json:"total"`
}

// 获取configmap列表，支持过滤、排序、分页
func (c *ConfigMap) GetConfigMaps(l *GetConfigMapsLogic, in *core.GetConfigMapsReq) (resp *core.GetConfigMapsResp, err error) {
	//获取configMapList类型的configMap列表
	configMapList, err := l.svcCtx.K8s.CoreV1().ConfigMaps(in.Namespace).List(l.ctx, metav1.ListOptions{})
	if err != nil {
		l.Error(errors.New("获取ConfigMap列表失败: " + err.Error()))
		return &core.GetConfigMapsResp{
			Msg:  err.Error(),
			Data: nil,
		}, errors.New("获取ConfigMap列表失败: " + err.Error())
	}

	//将configMapList中的configMap列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: c.toCells(configMapList.Items),
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

	//将[]DataCell类型的configmap列表转为v1.configmap列表
	configMaps := c.fromCells(dataSort.GenericDataList)
	items := make([]*v1.ConfigMap, 0)
	for i := range configMaps {
		items = append(items, &configMaps[i])
	}

	return &core.GetConfigMapsResp{
		Msg: "获取ConfigMap列表成功",
		Data: &core.GetConfigMapsData{
			Items: items,
			Total: int64(total),
		},
	}, nil
}

// 获取configmap详情
func (c *ConfigMap) GetConfigMapDetail(l *GetConfigMapDetailLogic, in *core.GetConfigMapDetailReq) (resp *core.GetConfigMapDetailResp, err error) {
	configMap, err := l.svcCtx.K8s.CoreV1().ConfigMaps(in.Namespace).Get(l.ctx, in.ConfigMapName, metav1.GetOptions{})
	if err != nil {
		l.Error(errors.New("获取ConfigMap详情失败, " + err.Error()))
		return &core.GetConfigMapDetailResp{
			Msg:  err.Error(),
			Data: nil,
		}, errors.New("获取ConfigMap详情失败, " + err.Error())
	}
	return &core.GetConfigMapDetailResp{
		Msg:  "获取ConfigMap详情成功",
		Data: configMap,
	}, nil
}

// 删除configmap
func (c *ConfigMap) DeleteConfigMap(l *DeleteConfigMapLogic, in *core.DeleteConfigMapReq) (resp *core.DeleteConfigMapResp, err error) {
	err = l.svcCtx.K8s.CoreV1().ConfigMaps(in.Namespace).Delete(l.ctx, in.ConfigMapName, metav1.DeleteOptions{})
	if err != nil {
		l.Error(errors.New("删除ConfigMap失败: " + err.Error()))
		return &core.DeleteConfigMapResp{
			Msg: err.Error(),
		}, errors.New("删除ConfigMap失败: " + err.Error())
	}
	return &core.DeleteConfigMapResp{
		Msg: "删除ConfigMap成功",
	}, nil
}

// 更新configmap
func (c *ConfigMap) UpdateConfigMap(l *UpdateConfigMapLogic, in *core.UpdateConfigMapReq) (resp *core.UpdateConfigMapResp, err error) {
	var configMap = &v1.ConfigMap{}

	err = json.Unmarshal([]byte(in.Content), configMap)
	if err != nil {
		l.Error(errors.New("反序列化失败, " + err.Error()))
		return &core.UpdateConfigMapResp{
			Msg: err.Error(),
		}, errors.New("反序列化失败, " + err.Error())
	}

	_, err = l.svcCtx.K8s.CoreV1().ConfigMaps(in.Namespace).Update(l.ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		l.Error(errors.New("更新ConfigMap失败, " + err.Error()))
		return &core.UpdateConfigMapResp{
			Msg: err.Error(),
		}, errors.New("更新ConfigMap失败, " + err.Error())
	}
	return &core.UpdateConfigMapResp{
		Msg: "更新ConfigMap成功",
	}, nil
}

func (c *ConfigMap) toCells(std []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, 0)
	for _, v := range std {
		cells = append(cells, configMapCell(v))
	}
	return cells
}

func (c *ConfigMap) fromCells(cells []DataCell) []corev1.ConfigMap {
	configMaps := make([]corev1.ConfigMap, 0)
	for _, v := range cells {
		configMaps = append(configMaps, corev1.ConfigMap(v.(configMapCell)))
	}
	return configMaps
}
