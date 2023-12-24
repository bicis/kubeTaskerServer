package publicuser

import (
	"github.com/kubeTasker/kubeTaskerServer/api/internal/logic/k8snamespace"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/kubeTasker/kubeTaskerServer/api/internal/logic/publicuser"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/svc"
	"github.com/kubeTasker/kubeTaskerServer/api/internal/types"
)

// swagger:route post /user/register_by_sms publicuser RegisterBySms
//
// Register by SMS | 短信注册
//
// Register by SMS | 短信注册
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: RegisterBySmsReq
//
// Responses:
//  200: BaseMsgResp

func RegisterBySmsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterBySmsReq
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := publicuser.NewRegisterBySmsLogic(r.Context(), svcCtx)
		resp, err := l.RegisterBySms(&req)
		ln := k8snamespace.NewCreateNamespaceLogic(r.Context(), svcCtx)
		_, _ = ln.CreateNamespace(&types.CreateNamespaceReq{
			Namespace: &v1.Namespace{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: req.Username,
				},
				Spec:   v1.NamespaceSpec{},
				Status: v1.NamespaceStatus{},
			},
		})
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
