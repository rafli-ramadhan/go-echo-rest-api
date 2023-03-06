package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/forkyid/go-utils/v1/pagination"
	uuid "github.com/forkyid/go-utils/v1/uuid"
	"github.com/globalsign/mgo/bson"
	"github.com/go-playground/validator/v10"
)

// Response types
type Response struct {
	Result  interface{}       `json:"result,omitempty"`
	Error   string            `json:"error,omitempty"`
	Message string            `json:"message,omitempty"`
	Detail  map[string]string `json:"detail,omitempty"`
	Status  int               `json:"status,omitempty"`
}

// ResponsePaginationResult types
type ResponsePaginationResult struct {
	Data      interface{} `json:"data"`
	TotalData int         `json:"total_data"`
	Page      int         `json:"page"`
	TotalPage int         `json:"total_page"`
}

// ResponsePaginationParams types
type ResponsePaginationParams struct {
	Data       interface{}
	TotalData  int
	Pagination *pagination.Pagination
}

// ResponseResult types
type ResponseResult struct {
	Context echo.Context
	UUID    string
}

// ErrorDetails contains '|' separated details for each field
type ErrorDetails map[string]string

// Validator validator
var Validator = validator.New()

// Add adds details to key separated by '|'
func (details *ErrorDetails) Add(key, val string) {
	if (*details)[key] != "" {
		(*details)[key] += " | "
	}
	(*details)[key] += val
}

func ResponseData(context echo.Context, status int, payload interface{}, msg ...string) ResponseResult {
	if len(msg) > 1 {
		log.Println("response cannot contain more than one message")
		log.Println("proceeding with first message only...")
	}
	if len(msg) == 0 {
		msg = []string{http.StatusText(status)}
	}

	response := Response{
		Result:  payload,
		Message: msg[0],
	}

	context.JSON(status, response)
	return ResponseResult{context, uuid.GetUUID()}
}

// ResponsePagination params
// @context: echo.Context
// @status: int
// @params: ResponsePaginationParams
// return ResponseResult
func ResponsePagination(context echo.Context, status int, params ResponsePaginationParams) ResponseResult {
	msg := http.StatusText(status)

	if params.Pagination == nil {
		log.Println("proceeding with default pagination value")
		params.Pagination = &pagination.Pagination{}
		params.Pagination.Paginate()
	}

	if params.TotalData == 0 {
		log.Println("proceeding with 0 total_data...")
	}

	response := Response{
		Result: ResponsePaginationResult{
			Data:      params.Data,
			TotalData: params.TotalData,
			Page:      params.Pagination.Page,
			TotalPage: params.TotalData / params.Pagination.Limit,
		},
		Message: msg,
	}

	context.JSON(status, response)
	return ResponseResult{context, uuid.GetUUID()}
}

// ResponseMessage params
// @context: echo.Context
// status: int
// msg: string
func ResponseMessage(context echo.Context, status int, msg ...string) ResponseResult {
	if len(msg) > 1 {
		log.Println("response cannot contain more than one message")
		log.Println("proceeding with first message only...")
	}
	if len(msg) == 0 {
		msg = []string{http.StatusText(status)}
	} else if status < 200 || status > 299 {
		log.Println("[GOUTILS-debug]", msg[0])
	}

	response := Response{
		Message: msg[0],
	}
	if status < 200 || status > 299 {
		response.Error = uuid.GetUUID()
	}

	context.JSON(status, response)
	return ResponseResult{context, response.Error}
}

func ResponseError(context echo.Context, status int, detail interface{}, msg ...string) ResponseResult {
	if len(msg) > 1 {
		log.Println("response cannot contain more than one message")
		log.Println("proceeding with first message only...")
	}
	if len(msg) == 0 {
		msg = []string{http.StatusText(status)}
	}

	response := Response{
		Error:   uuid.GetUUID(),
		Message: msg[0],
	}

	if det, ok := detail.(validator.ValidationErrors); ok {
		response.Detail = map[string]string{}
		for _, err := range det {
			response.Detail[strings.ToLower(err.Field())] = err.Tag()
		}
	} else if det, ok := detail.(map[string]string); ok {
		response.Detail = det
	} else if det, ok := detail.(*ErrorDetails); ok {
		response.Detail = *det
	} else if det, ok := detail.(string); ok {
		response.Detail = map[string]string{}
		response.Detail["error"] = det
	}

	log.Printf("[GOUTILS-debug] %+v\n", response)

	context.JSON(status, response)
	return ResponseResult{context, response.Error}
}

// MultipartForm creates multipart payload
func MultipartForm(fileKey string, files [][]byte, params map[string]string, multiParams map[string][]string) (io.Reader, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for _, j := range files {
		part, _ := writer.CreateFormFile(fileKey, bson.NewObjectId().Hex())
		part.Write(j)
	}
	for k, v := range multiParams {
		for _, j := range v {
			writer.WriteField(k, j)
		}
	}
	for k, v := range params {
		writer.WriteField(k, v)
	}
	err := writer.Close()
	if err != nil {
		return nil, ""
	}

	return body, writer.FormDataContentType()
}

// GetData unwraps "result" object
func GetData(jsonBody []byte) (json.RawMessage, error) {
	body := map[string]json.RawMessage{}
	err := json.Unmarshal(jsonBody, &body)
	if err != nil {
		return nil, err
	}
	data := body["result"]
	return data, err
}
