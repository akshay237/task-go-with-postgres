package config

import (
	"os"
	"strings"
	"task-go-with-postgres/database"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
)

type ConfigError struct {
	Op   string
	Desc string
	Err  error
}

func (e *ConfigError) Error() string {
	return "#ConfigError|Op:" + e.Op + "|Desc: " + e.Desc + "|Err:" + e.Err.Error()
}

type ConfigAPIServer struct {
	Port *int `json:"port" validate:"required,numeric"`
}

type ServerConfig struct {
	ScratchDir      *string                  `json:"scratchdir" validate:"required"`
	PostgresConfig  *database.ConfigPostgres `json:"postgres" validate:"required"`
	APIServerConfig *ConfigAPIServer         `json:"apiserver" validate:"required"`
}

func readFileFullContent(filepath string) ([]byte, error) {
	rawjsonbytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, &ConfigError{
			Op:   "FILE_READ",
			Desc: "File opening has failed",
			Err:  err,
		}
	}
	return rawjsonbytes, nil
}

func validateJsonStruct(structObj any) error {
	validate := validator.New()
	err := validate.Struct(structObj)
	if err != nil {
		msg := "Validation failed for fields: "
		fields := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			fields = append(fields, strings.ToLower(err.Field()))
		}
		msg = msg + strings.Join(fields, ", ")
		return &ConfigError{
			Op:   "JSON_VALIDIATION",
			Desc: msg,
			Err:  err,
		}
	} else {
		return nil
	}
}

func ParseConfigFromFile(configfilepath string) (*ServerConfig, error) {
	configfilebytes, read_err := readFileFullContent(configfilepath)
	if read_err != nil {
		return nil, read_err
	}

	serverConfig := &ServerConfig{}
	unmarshal_err := sonic.Unmarshal(configfilebytes, &serverConfig)
	if unmarshal_err != nil {
		return nil, &ConfigError{
			Op:   "JSON_UNMARSHAL",
			Desc: "Malformed JSON",
			Err:  unmarshal_err,
		}
	}

	validate_err := validateJsonStruct(serverConfig)
	if validate_err != nil {
		return nil, validate_err
	}

	return serverConfig, nil
}
