package k8snamespace

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kubeTasker/kubeTaskerServer/api/internal/logic/k8snamespace"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/types"
)

// swagger:route post /namespace/create_namespace k8snamespace CreateNamespace
//
// createNamespace
//
// createNamespace
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: CreateNamespaceReq
//
// Responses:
//  200: CreateNamespaceResp

func CreateNamespaceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateNamespaceReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := k8snamespace.NewCreateNamespaceLogic(r.Context(), svcCtx)
		resp, err := l.CreateNamespace(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
