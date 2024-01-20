package apiserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type APIHandler interface {
	RegisterRoutes(r chi.Router)
}

var (
	ErrMissingAPIArguments    = errors.New("MISSING_API_ARGUMENTS")
	ErrInsufficientPermission = errors.New("INSUFFICIENT_PERMISSION")
	ErrUnknown                = errors.New("UNKNOWN_ERROR")
)

var (
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	JSON_PARSE_ERROR      = "JSON_PARSE_ERROR"
)

// const errdefaultresponse = "{\"err\":{\"errcode\":\"UNKNOWN_ERROR\"},\"data\":null,\"msg\":\"Unknown error\"}"

type APIResponseErrJson struct {
	Errcode string      `json:"errcode"`
	ErrData interface{} `json:"errdata"`
}

type APIResponse struct {
	Err  interface{} `json:"err"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func APIResponseOK(w http.ResponseWriter, r *http.Request, data interface{}, msg string) error {
	responseobj := &APIResponse{
		Err:  nil,
		Data: data,
		Msg:  msg,
	}
	jsbarr, _ := json.Marshal(responseobj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeerr := w.Write(jsbarr)
	return writeerr
}

func apiResponseErr(w http.ResponseWriter, r *http.Request, status int, errcode string, errdata interface{}, msg string) error {
	responseobj := &APIResponse{
		Err: &APIResponseErrJson{
			Errcode: errcode,
			ErrData: errdata,
		},
		Data: nil,
		Msg:  msg,
	}
	jsbarr, _ := json.Marshal(responseobj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, writeerr := w.Write(jsbarr)
	return writeerr
}

func APIResponseBadRequest(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string) error {
	return apiResponseErr(w, r, http.StatusBadRequest, errcode, errdata, msg)
}

func APIResponseUnauthorized(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string) error {
	return apiResponseErr(w, r, http.StatusUnauthorized, errcode, errdata, msg)
}

func APIResponseForbidden(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string) error {
	return apiResponseErr(w, r, http.StatusForbidden, errcode, errdata, msg)
}

func APIResponseConflict(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string) error {
	return apiResponseErr(w, r, http.StatusConflict, errcode, errdata, msg)
}

func APIResponseGone(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string) error {
	return apiResponseErr(w, r, http.StatusGone, errcode, errdata, msg)
}

func APIResponseUnprocessableEntity(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string) error {
	return apiResponseErr(w, r, http.StatusUnprocessableEntity, errcode, errdata, msg)
}

func APIResponseNotAcceptable(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string) error {
	return apiResponseErr(w, r, http.StatusNotAcceptable, errcode, errdata, msg)
}

func APIResponseInternalServerError(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string) error {
	return apiResponseErr(w, r, http.StatusInternalServerError, errcode, errdata, msg)
}

func APIFailedInternalAPICall(w http.ResponseWriter, r *http.Request, errcode string, errdata interface{}, msg string, statusCode int) error {
	return apiResponseErr(w, r, statusCode, errcode, errdata, msg)
}

func JsonBodyParser(r io.Reader, reqVar interface{}) (interface{}, error) {

	body, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &reqVar)
	if err != nil {
		return "", err
	}

	return ValidateStruct(reqVar)
}

type ErrStructure struct {
	Field string `json:"field"`
	Param string `json:"param"`
	Tag   string `json:"err"`
	Info  string `json:"info"`
}

func ValidateStruct(structVar interface{}) (interface{}, error) {
	validate := validator.New()
	if err := validate.RegisterValidation("regex", validateRegex); err != nil {
		log.Err(err).Msg("register failed")
		return "", err
	}

	err := validate.Struct(structVar)
	if err != nil {
		msg := "Validation failed for field(s): "
		fields := []string{}
		errStructs := []ErrStructure{}

		for _, err := range err.(validator.ValidationErrors) {

			errStruct, regexErrMessgae := getRegexErrorMessage(err.Field(), err.Param(), err.Tag())

			fields = append(fields, strings.ToLower(err.Field())+regexErrMessgae)

			errStructs = append(errStructs, errStruct)
		}
		msg = msg + strings.Join(fields, "; ")

		return errStructs, errors.New(msg)
	}

	return "", nil

}

func validateRegex(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	regex := fl.Param()
	match, _ := regexp.MatchString(regex, value)
	return match
}

func getRegexErrorMessage(field, param, tag string) (ErrStructure, string) {

	var errMessage string

	errMessage = `field ` + field + ` doesn't present or contains unwanted charachters`

	switch tag {
	case "max":
		errMessage = errMessage + ` field ` + field + ` can't have more than ` + param + ` charachters`
	case "min":
		errMessage = errMessage + ` field ` + field + ` should have mminimum ` + param + ` charachters`
	}

	return ErrStructure{
		Field: field,
		Param: param,
		Tag:   tag,
		Info:  errMessage,
	}, errMessage
}
