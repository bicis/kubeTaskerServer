package errorhandler

import (
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-common/msg/logmsg"

	"github.com/suyuan32/simple-admin-common/i18n"

	"github.com/kubeTasker/kubeTaskerServer/rpc/ent"
)

// UserEntError returns errors dealing with User functions.
func UserEntError(logger logx.Logger, err error, detail any) error {
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return errorx.NewInvalidArgumentError(i18n.TargetNotFound)
		case ent.IsConstraintError(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return errorx.NewInvalidArgumentError(i18n.UserRepeat)
		case ent.IsValidationError(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return errorx.NewInvalidArgumentError(i18n.ValidationError)
		case ent.IsNotSingular(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return errorx.NewInvalidArgumentError(i18n.NotSingularError)
		default:
			logger.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
			return errorx.NewInternalError(i18n.DatabaseError)
		}
	}
	return err
}
