package helpers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Log     *Log        `json:"log,omitempty"`
}

type SuccessResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Meta    *Meta        `json:"meta,omitempty"`
}

type Meta struct {
	RequestID  string      `json:"request_id"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	CurrentPage int     `json:"current_page"`
	TotalItems  int     `json:"total_items"`
	TotalPages  int     `json:"total_pages"`
	ItemPerPage int     `json:"item_per_page"`
	FromRow     int     `json:"from_row"`
	ToRow       int     `json:"to_row"`
	Self        string  `json:"self"`
	Next        *string `json:"next"`
	Prev        *string `json:"prev"`
}

type ErrorResponse struct {
	Status  int          `json:"status"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Errors  *interface{} `json:"errors,omitempty"`
}

func ResponseFormatter(c *fiber.Ctx, res BaseResponse) error {
	// Insert log
	var username string
	if sessionUsername := c.Locals("username"); sessionUsername != nil {
		username = sessionUsername.(string)
	} else {
		username = ""
	}
	// log.Println(res.Log)

	// Default value if res.Log is nil
	var location string
	var startTime time.Time

	if res.Log != nil {
		location = res.Log.Location
		startTime = res.Log.StartTime
	}

	// log.Printf("Unhandled error in Response Formatter: %v", res.Errors)

	logSysData := LogSystemParam{
		Identifier: c.GetRespHeader(fiber.HeaderXRequestID),
		StatusCode: res.Status,
		Location:   location,
		Message:    res.Message,
		StartTime:  startTime,
		EndTime:    time.Now(),
		Err:        res.Errors,
		Username:   username,
	}

	LogSysChannel <- logSysData

	res.Log = nil
	if res.Status != fiber.StatusBadRequest ||
		(res.Message != "Invalid or malformed request query" && res.Message != "Invalid or malformed request body") {
		res.Errors = nil
	}

	return c.Status(res.Status).JSON(res)
}
