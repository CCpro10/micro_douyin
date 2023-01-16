package util

import "errors"

const (
	Success = iota

	WrongAuth
	NoAuth

	UserExisted
	UserNotExist
	WrongPassword

	VideoNotExist

	CommentNotExist
	CommentTooLong
	CommentIsEmpty

	SensitiveComment

	IsFollow
	NotFollow

	InternalServerError
	ParamError
)

var ErrCode2Msg = map[int]string{
	Success:             "请求成功",
	WrongAuth:           "用户登录已过期,请重新登陆",
	NoAuth:              "权限不足,请登录后查看哦",
	UserExisted:         "用户名已存在",
	UserNotExist:        "用户不存在",
	WrongPassword:       "密码错误",
	IsFollow:            "该关注记录已存在",
	NotFollow:           "该关注记录尚未存在",
	VideoNotExist:       "该视频不存在",
	CommentNotExist:     "该评论不存在",
	CommentTooLong:      "评论过长（内容需在500个字符以内）",
	CommentIsEmpty:      "评论不能为空",
	SensitiveComment:    "评论包含敏感词汇",
	InternalServerError: "服务器内部错误",
	ParamError:          "参数错误",
}

var (
	ErrWrongAuth        = errors.New(ErrCode2Msg[WrongAuth])
	ErrNoAuth           = errors.New(ErrCode2Msg[NoAuth])
	ErrUserExisted      = errors.New(ErrCode2Msg[UserExisted])
	ErrWrongPassword    = errors.New(ErrCode2Msg[WrongPassword])
	ErrIsFollow         = errors.New(ErrCode2Msg[IsFollow])
	ErrNotFollow        = errors.New(ErrCode2Msg[NotFollow])
	ErrSensitiveComment = errors.New(ErrCode2Msg[SensitiveComment])
	ErrParamError       = errors.New(ErrCode2Msg[ParamError])
)
