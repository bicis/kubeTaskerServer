package k8srolebinding

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kubeTasker/kubeTaskerServer/api/internal/logic/k8srolebinding"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/types"
)

// swagger:route post /k8srolebinding/create_role_binding k8srolebinding CreateRoleBinding
//
// createRoleBinding
//
// createRoleBinding
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: CreateRoleBindingReq
//
// Responses:
//  200: CreateRoleBindingResp

func CreateRoleBindingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateRoleBindingReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := k8srolebinding.NewCreateRoleBindingLogic(r.Context(), svcCtx)
		resp, err := l.CreateRoleBinding(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
