package middleware

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/utils"
)

func LogMiddleware(inp map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		trace := true

		fields := map[string]interface{}{
			consts.ContextRequestURI:         utils.ConstructURL(c.Request),
			consts.ContextRequestMethod:      c.Request.Method,
			consts.ContextRequestIP:          c.ClientIP(),
			consts.ContextRequestID:          utils.GetRequestIDFromRequest(c.Request),
			consts.ContextRequestURITemplate: utils.GetRequestRoute(c),
			consts.ContextService:            logger.GetService(),
		}

		if logger.GetRequestDumpStatus() {
			if req, err := utils.GetRequestDump(c.Request); err == nil {
				fields[consts.ContextRequestDump] = *req
			}
		}
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, consts.LogData, fields)
		c.Request = c.Request.WithContext(ctx)

		defer func() {
			if trace {
				traceMsg := fmt.Sprintf("Stacktrace: %v", string(debug.Stack()))
				logger.
					Log().
					WithContext(ctx).
					Panic(traceMsg)
			}
		}()
		ww := utils.NewResponseWriterWrapper(c.Writer)
		c.Writer = ww

		start := time.Now()
		logger.Log().WithContext(ctx).Info("started handling request")

		c.Next()

		fields = map[string]interface{}{
			consts.ContextRequestStatus:    c.Writer.Status(),
			consts.ContextRequestTimetaken: time.Since(start).String(), // Convert duration to a string
		}
		if logger.GetResponseDumpStatus() {
			fields[consts.ContextResponseDump] = ww.GetResponseData()
		}
		logger.
			Log().
			WithContext(ctx).
			WithFields(fields).
			Info("completed handling request")
		trace = false
	}
}
