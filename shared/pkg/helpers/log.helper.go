package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/natefinch/lumberjack"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	apiLogger    *zap.Logger
	systemLogger *zap.Logger
)

var LogAPIChannel = make(chan LogAPIParam, 100)
var LogSysChannel = make(chan LogSystemParam, 100)

type LogAPIParam struct {
	StatusCode    int                    `json:"status_code"`
	Message       string                 `json:"message"`
	Identifier    string                 `json:"identifier"`
	Timestamp     time.Time              `json:"timestamp"`
	HTTPMethod    string                 `json:"http_method"`
	RequestHeader interface{}            `json:"request_header"`
	QueryParams   interface{}            `json:"query_params"`
	RequestBody   interface{}            `json:"request_body"`
	ResponseCode  string                 `json:"response_code"`
	ResponseBody  map[string]interface{} `json:"response_body"`
	Endpoint      string                 `json:"endpoint"`
	OriginalURL   string                 `json:"original_url"`
	UserAgent     string                 `json:"user_agent"`
	ClientIP      string                 `json:"client_ip"`
	Username      string                 `json:"username"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
}

type LogSystemParam struct {
	Identifier string
	StatusCode int
	Location   string
	Message    string
	StartTime  time.Time
	EndTime    time.Time
	Duration   string
	Username   string
	Err        interface{}
}

type Log struct {
	StartTime time.Time
	EndTime   time.Time
	Location  string
	Message   string
	Err       interface{}
}

func InitLogger() {
	// Create logs directory if it does not exist
	if err := os.MkdirAll("storage/logs/api", os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("storage/logs/system", os.ModePerm); err != nil {
		panic(err)
	}

	// Encoder configuration
	cfg := config.AppConfig
	debugMode := cfg.Debug

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Dynamic log filename based on current month
	currentTime := time.Now()
	logAPIFilename := "storage/logs/api/api_" + currentTime.Format("2006-01") + ".log" // Format as YYYY-MM

	// API Logger setup with lumberjack for log rotation
	apiFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logAPIFilename, // Log file path
		MaxAge:     365,            // Keep log files for 365 days (1 year)
		MaxBackups: 12,             // Keep 12 backups (1 per month for a year)
		Compress:   true,           // Compress old log files
	})

	// Dynamic log filename based on current month
	logSysFilename := "storage/logs/system/system_" + currentTime.Format("2006-01") + ".log" // Format as YYYY-MM

	// System Logger setup with lumberjack for log rotation
	systemFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logSysFilename, // Log file path
		MaxAge:     365,            // Keep log files for 365 days (1 year)
		MaxBackups: 12,             // Keep 12 backups (1 per month for a year)
		Compress:   true,           // Compress old log files
	})

	// Console output configuration
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	var apiCores []zapcore.Core
	var systemCores []zapcore.Core

	if debugMode {
		apiCores = append(apiCores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
		systemCores = append(systemCores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
	}
	apiCores = append(apiCores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), apiFileWriter, zapcore.InfoLevel))
	systemCores = append(systemCores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), systemFileWriter, zapcore.InfoLevel))

	apiLogger = zap.New(zapcore.NewTee(apiCores...))
	systemLogger = zap.New(zapcore.NewTee(systemCores...))
}

func APILogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Read request body
		var requestBody interface{}
		contentType := c.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// Handle multipart form-data
			form, err := c.MultipartForm()
			if err != nil {
				requestBody = fmt.Sprintf("Error parsing form-data: %v", err)
			} else {
				formData := make(map[string]interface{})
				for key, values := range form.Value {
					if key == "password" || key == "re_password" || key == "old_password" {
						formData[key] = "[REDACTED]"
					} else {
						formData[key] = values
					}
				}
				requestBody = formData
			}
		} else {
			bodyBytes := c.Body()
			if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				// Handle form-urlencoded
				formData, err := url.ParseQuery(string(bodyBytes))
				if err != nil {
					requestBody = string(bodyBytes)
				} else {
					jsonFormData := make(map[string]interface{})
					for key, values := range formData {
						if key == "password" || key == "re_password" || key == "old_password" || key == "raw" {
							jsonFormData[key] = "[REDACTED]"
						} else {
							jsonFormData[key] = values
						}
					}
					requestBody = jsonFormData
				}
			} else {
				err := json.Unmarshal(bodyBytes, &requestBody)
				if err != nil {
					requestBody = string(bodyBytes)
				}
			}
		}

		// Let the request proceed
		err := c.Next()

		// Capture response
		statusCode := c.Response().StatusCode()
		responseBody := c.Response().Body()
		contentType = c.GetRespHeader("Content-Type") // Re-fetch content type for response
		endTime := time.Now()

		var jsonResponseBody map[string]interface{}
		if isBinaryContent(contentType) {
			// If the response is binary content (e.g. PDF, image), do not log the raw body, log metadata instead
			jsonResponseBody = map[string]interface{}{
				"content_type": contentType,
				"size":         len(responseBody),
				"message":      "Binary content returned, body is not logged.",
			}
		} else {
			// Process non-binary response body for logging
			if err := json.Unmarshal(responseBody, &jsonResponseBody); err != nil {
				jsonResponseBody = map[string]interface{}{
					"raw": string(responseBody),
				}
			} else {
				RedactFields(jsonResponseBody, []string{"key", "token", "password", "re_password", "old_password", "raw"})
			}
		}

		// Redact Authorization header
		headers := c.GetReqHeaders()
		if _, exists := headers["Authorization"]; exists {
			headers["Authorization"] = []string{"[REDACTED]"}
		}

		// Get other data
		identifier := c.GetRespHeader(fiber.HeaderXRequestID)
		userAgent := string(c.Request().Header.UserAgent())
		clientIP := c.IP()
		endpoint := c.Path()
		queryParams := c.Request().URI().QueryArgs()
		httpMethod := c.Method()
		statusCodeString := strconv.Itoa(statusCode)
		originalURL := c.OriginalURL()

		// Get username from session or set as empty string if nil
		var username string
		if sessionUsername := c.Locals("username"); sessionUsername != nil {
			username = sessionUsername.(string)
		} else {
			username = ""
		}

		// Log message
		message := "API LOG"
		if msg, ok := jsonResponseBody["message"]; ok {
			message = fmt.Sprintf("API Log | %s", msg.(string))
		}

		apiLogData := LogAPIParam{
			StatusCode:    statusCode,
			Message:       message,
			Identifier:    identifier,
			Timestamp:     time.Now(),
			HTTPMethod:    httpMethod,
			RequestHeader: headers,
			QueryParams:   queryParams,
			RequestBody:   requestBody,
			ResponseCode:  statusCodeString,
			ResponseBody:  jsonResponseBody,
			Endpoint:      endpoint,
			OriginalURL:   originalURL,
			UserAgent:     userAgent,
			ClientIP:      clientIP,
			Username:      username,
			StartTime:     startTime,
			EndTime:       endTime,
		}

		LogAPIChannel <- apiLogData

		return err
	}
}

// Helper function to determine if the content type is binary
func isBinaryContent(contentType string) bool {
	return !strings.HasPrefix(contentType, "application/json") && !strings.HasPrefix(contentType, "text/plain")
}

func GenerateLogAPI(apiLogData LogAPIParam) {

	// Log the request and response
	logFunc := apiLogger.Info
	if apiLogData.StatusCode >= 500 {
		logFunc = apiLogger.Error
	} else if apiLogData.StatusCode >= 400 {
		logFunc = apiLogger.Warn
	}

	duration := FormatDuration(apiLogData.StartTime, apiLogData.EndTime)

	logFunc(apiLogData.Message,
		zap.String("identifier", apiLogData.Identifier),
		zap.Time("timestamp", apiLogData.Timestamp),
		zap.String("http_method", apiLogData.HTTPMethod),
		zap.Any("request_header", apiLogData.RequestHeader),
		zap.Any("query_params", apiLogData.QueryParams),
		zap.Any("request_body", apiLogData.RequestBody),
		zap.String("response_code", apiLogData.ResponseCode),
		zap.Any("response_body", apiLogData.ResponseBody),
		zap.String("endpoint", apiLogData.Endpoint),
		zap.String("original_url", apiLogData.OriginalURL),
		zap.String("user_agent", apiLogData.UserAgent),
		zap.String("client_ip", apiLogData.ClientIP),
		zap.String("username", apiLogData.Username),
		zap.Time("start_time", apiLogData.StartTime),
		zap.Time("end_time", apiLogData.EndTime),
		zap.String("duration", duration),
	)
}

func GenerateLogSystem(logData LogSystemParam) {
	var (
		category         string
		humanTime        = logData.EndTime.Format(time.RFC1123)
		statusCodeString = strconv.Itoa(logData.StatusCode)
		duration         = FormatDuration(logData.StartTime, logData.EndTime)
	)

	switch {
	case logData.StatusCode >= 500:
		category = "FATAL"
	case logData.StatusCode >= 400:
		category = "ERROR"
	default:
		category = "INFO"
	}

	systemLogger.Info("System Log",
		zap.Time("timestamp", time.Now()),
		zap.String("category", category),
		zap.String("response_code", statusCodeString),
		zap.String("location", logData.Location),
		zap.String("message", logData.Message),
		zap.Time("start_time", logData.StartTime),
		zap.Time("end_time", logData.EndTime),
		zap.String("duration", duration),
		zap.String("identifier", logData.Identifier),
		zap.Any("username", logData.Username),
		zap.Any("errors", logData.Err),
		zap.String("human_time", humanTime),
	)
}

// Helper function to redact sensitive fields in a map
func RedactFields(data map[string]interface{}, fields []string) {
	for _, field := range fields {
		if _, ok := data[field]; ok {
			data[field] = "[REDACTED]"
		}
	}
	for _, value := range data {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			RedactFields(nestedMap, fields)
		} else if nestedArray, ok := value.([]interface{}); ok {
			for _, item := range nestedArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					RedactFields(itemMap, fields)
				}
			}
		}
	}
}

// GetAPILogger returns the initialized apiLogger instance
func GetAPILogger() *zap.Logger {
	if apiLogger == nil {
		InitLogger()
	}
	return apiLogger
}

func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func GetFunctionAndStructName(i interface{}) (string, string, string) {
	pc, file, _, ok := runtime.Caller(2) // 2 untuk mengambil caller dari fungsi pemanggil
	if !ok {
		return "", "", ""
	}
	fullFuncName := runtime.FuncForPC(pc).Name()
	funcName := fullFuncName[strings.LastIndex(fullFuncName, ".")+1:]

	contrName := reflect.TypeOf(i).Elem().Name()

	projectRoot := getProjectRoot()

	relativePath, err := filepath.Rel(projectRoot, file)
	if err != nil {
		relativePath = file
	}

	packagePath := strings.ReplaceAll(relativePath, string(filepath.Separator), "/")
	packagePath = strings.TrimSuffix(packagePath, filepath.Ext(packagePath))

	projectName := filepath.Base(projectRoot)
	if idx := strings.Index(packagePath, projectName); idx != -1 {
		packagePath = packagePath[idx+len(projectName)+1:]
	}

	return funcName, contrName, packagePath
}

func CreateLog(i interface{}) Log {
	funcName, contrName, packagePath := GetFunctionAndStructName(i)
	return Log{
		StartTime: time.Now(),
		Location:  fmt.Sprintf("%s/%s.%s", packagePath, contrName, funcName),
		Message:   "Passed",
	}
}

func InitialLogSystem() Log {
	return Log{
		StartTime: time.Now(),
		Message:   "Passed",
	}
}

// Function to create log
func CreateLogSystem(ctx context.Context, message string) *Log {
	// Extract log data from context
	identifier, _ := ctx.Value("identifier").(string)
	username, _ := ctx.Value("username").(string)
	funcName, contrName, packagePath := GetFunctionAndStructName(ctx.Value("function"))

	return &Log{
		StartTime: time.Now(),
		Location:  fmt.Sprintf("%s/%s.%s", packagePath, contrName, funcName),
		Message:   fmt.Sprintf("User: %s, ID: %s, Message: %s", username, identifier, message),
	}
}

func CreateLogSystem23(ctx context.Context, logData *Log) {
	// func CreateLogSystem23(ctx context.Context, location string, message string, startTime time.Time, err interface{}) {
	// Extract function and location
	// funcName, contrName, packagePath := GetFunctionAndStructName(i)

	// Create Log
	logSysData := LogSystemParam{
		StartTime: logData.StartTime,
		EndTime:   time.Now(),
		Location:  logData.Location,
		Message:   logData.Message,
		Err:       logData.Err,
	}

	// Get identifier and username from context
	if identifier, ok := ctx.Value("identifier").(string); ok {
		logSysData.Identifier = identifier
	}
	if username, ok := ctx.Value("username").(string); ok {
		logSysData.Username = username
	}

	// duration := FormatDuration(logSysData.StartTime, logSysData.EndTime)

	// Log debug waktu start dan end
	// log.Printf("Log for %s - StartTime: %s, EndTime: %s, Duration: %s\n",
	// 	logSysData.Location, logData.StartTime.Format(time.RFC3339), logSysData.EndTime.Format(time.RFC3339), duration)
	// log.Println("=============================================================")

	// LogSysChannel could be where your logs are processed
	LogSysChannel <- logSysData
}

// Extract identifier and username and insert into context
func ExtractIdentifierAndUsername(c *fiber.Ctx) context.Context {
	ctx := context.Background()

	identifier := c.GetRespHeader(fiber.HeaderXRequestID)
	username := ""
	if sessionUsername := c.Locals("username"); sessionUsername != nil {
		username = sessionUsername.(string)
	}
	var user_id float64
	if sessionUserID := c.Locals("user_id"); sessionUserID != nil {
		user_id = sessionUserID.(float64)
	}
	var is_admin bool
	if sessionUserAdmin := c.Locals("is_admin"); sessionUserAdmin != nil {
		is_admin = sessionUserAdmin.(bool)
	}

	ctx = context.WithValue(ctx, CtxKeyIdentifier, identifier)
	ctx = context.WithValue(ctx, CtxKeyUsername, username)
	ctx = context.WithValue(ctx, CtxKeyUserID, user_id)
	ctx = context.WithValue(ctx, CtxKeyIsAdmin, is_admin)

	return ctx
}

func InitialLogExtractIdentifierAndUsername(c *fiber.Ctx, i interface{}) (context.Context, Log) {
	ctx := context.Background()

	// extract
	identifier := c.GetRespHeader(fiber.HeaderXRequestID)
	username := ""
	if sessionUsername := c.Locals("name"); sessionUsername != nil {
		username = sessionUsername.(string)
	}

	// passing to context
	ctx = context.WithValue(ctx, CtxKeyIdentifier, identifier)
	ctx = context.WithValue(ctx, CtxKeyUsername, username)

	log := CreateLog(i)

	return ctx, log
}

func FormatDuration(startTime, endTime time.Time) string {
	duration := endTime.Sub(startTime)
	durationInMilliseconds := duration.Seconds() * 1000 // Menghitung durasi dalam milidetik
	return fmt.Sprintf("%.4fms", durationInMilliseconds)
}

func LogSystemWithDefer(ctx context.Context, logData *Log) func() {
	// Log masuk
	// logData.Message = "Entering"
	CreateLogSystem23(ctx, logData)
	log.Println(logData)

	// `defer` yang akan dieksekusi saat fungsi selesai
	return func() {
		if logData.Err != nil {
			logData.Message = "Exiting with error: " + logData.Err.(string)
		} else {
			logData.Message = "Exiting"
		}

		CreateLogSystem23(ctx, logData)
	}
}

func LogBaseResponse(logData *Log, response BaseResponse) BaseResponse {
	logData.Message = response.Message
	logData.Err = response.Errors
	return response
}
