package errcode

import (
	errors "errors"
	strings "github.com/hopeio/cherry/utils/strings"
	io "io"
)

func (x ErrCode) String() string {

	switch x {
	case Success:
		return "成功"
	case Canceled:
		return "操作取消"
	case Unknown:
		return "未知错误"
	case InvalidArgument:
		return "无效的参数"
	case DeadlineExceeded:
		return "响应超时"
	case NotFound:
		return "未发现"
	case AlreadyExists:
		return "已经存在"
	case PermissionDenied:
		return "操作无权限"
	case ResourceExhausted:
		return "资源不足"
	case FailedPrecondition:
		return "操作被拒绝"
	case Aborted:
		return "操作终止"
	case OutOfRange:
		return "超出范围"
	case Unimplemented:
		return "未实现"
	case Internal:
		return "内部错误"
	case Unavailable:
		return "服务不可用"
	case DataLoss:
		return "数据丢失"
	case Unauthenticated:
		return "身份未验证"
	case SysError:
		return "系统错误"
	case DBError:
		return "数据库错误"
	case RowExists:
		return "记录已存在"
	case RedisErr:
		return "Redis错误"
	case IOError:
		return "io错误"
	case UploadFail:
		return "上传失败"
	case UploadCheckFail:
		return "检查文件失败"
	case UploadCheckFormat:
		return "文件格式或大小有问题"
	case TimeTooMuch:
		return "尝试次数过多"
	case ParamInvalid:
		return "参数错误"
	}
	return ""
}

func (x ErrCode) MarshalGQL(w io.Writer) {
	w.Write(strings.QuoteToBytes(x.String()))
}

func (x *ErrCode) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = ErrCode(i)
		return nil
	}
	return errors.New("enum need integer type")
}
