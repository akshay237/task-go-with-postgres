package apiserver

type ApiError struct {
	ErrCode string      `json:"errcode"`
	ErrData interface{} `json:"errdata"`
	ErrMsg  string      `json:"errmsg"`
}

func (o *ApiError) Error() string {
	return o.ErrMsg
}

func NewApiError(errcode string, errdata interface{}, errmsg string) *ApiError {
	return &ApiError{ErrCode: errcode, ErrData: errdata, ErrMsg: errmsg}
}
