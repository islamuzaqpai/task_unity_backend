package handler

import (
	"enactus/internal/apperrors"
	"enactus/internal/helpers"
	"enactus/internal/httpx"
	"enactus/internal/models/inputs"
	"enactus/internal/service"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AttendanceSessionHandler struct {
	AttendanceSessionS *service.AttendanceSessionService
}

func NewAttendanceSessionHandler(attendanceSessionS *service.AttendanceSessionService) *AttendanceSessionHandler {
	return &AttendanceSessionHandler{AttendanceSessionS: attendanceSessionS}
}

func (h *AttendanceSessionHandler) CreateSession(c *gin.Context) error {
	ctx := c.Request.Context()

	var req inputs.CreateAttendanceSessionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	v := helpers.NewValidator()
	if errs := helpers.Validate(req, v); errs != nil {
		return httpx.BadRequestValidation(errs)
	}

	currentUserId, role, err := getCurrentUserFromContext(c)
	if err != nil {
		return err
	}

	session, err := h.AttendanceSessionS.CreateSession(ctx, currentUserId, role, req)
	if err != nil {
		return mapAttendanceSessionError(err)
	}

	httpx.WriteJSON(c, 201, session)
	return nil
}

func (h *AttendanceSessionHandler) GetSession(c *gin.Context) error {
	ctx := c.Request.Context()

	sessionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid session id")
	}

	currentUserId, role, err := getCurrentUserFromContext(c)
	if err != nil {
		return err
	}

	session, err := h.AttendanceSessionS.GetSession(ctx, currentUserId, role, sessionId)
	if err != nil {
		return mapAttendanceSessionError(err)
	}

	httpx.WriteJSON(c, 200, session)
	return nil
}

func (h *AttendanceSessionHandler) ListSessions(c *gin.Context) error {
	ctx := c.Request.Context()

	currentUserId, role, err := getCurrentUserFromContext(c)
	if err != nil {
		return err
	}

	var departmentId *int
	if departmentIdStr := c.Query("department_id"); departmentIdStr != "" {
		parsed, parseErr := strconv.Atoi(departmentIdStr)
		if parseErr != nil {
			return httpx.BadRequest("invalid department_id")
		}
		departmentId = &parsed
	}

	var dateFrom *string
	if value := c.Query("date_from"); value != "" {
		dateFrom = &value
	}
	var dateTo *string
	if value := c.Query("date_to"); value != "" {
		dateTo = &value
	}

	sessions, err := h.AttendanceSessionS.ListSessions(ctx, currentUserId, role, departmentId, dateFrom, dateTo)
	if err != nil {
		return mapAttendanceSessionError(err)
	}

	httpx.WriteJSON(c, 200, sessions)
	return nil
}

func (h *AttendanceSessionHandler) BulkUpsertEntries(c *gin.Context) error {
	ctx := c.Request.Context()

	sessionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid session id")
	}

	var req inputs.BulkUpsertAttendanceEntriesInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}
	if len(req.Entries) == 0 {
		return httpx.BadRequest("entries are required")
	}

	v := helpers.NewValidator()
	if errs := helpers.Validate(req.Entries, v); errs != nil {
		return httpx.BadRequestValidation(errs)
	}

	currentUserId, role, err := getCurrentUserFromContext(c)
	if err != nil {
		return err
	}

	if err := h.AttendanceSessionS.BulkUpsertEntries(ctx, currentUserId, role, sessionId, req); err != nil {
		return mapAttendanceSessionError(err)
	}

	c.Status(200)
	return nil
}

func (h *AttendanceSessionHandler) PublishSession(c *gin.Context) error {
	ctx := c.Request.Context()

	sessionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid session id")
	}

	currentUserId, role, err := getCurrentUserFromContext(c)
	if err != nil {
		return err
	}

	if err := h.AttendanceSessionS.PublishSession(ctx, currentUserId, role, sessionId); err != nil {
		return mapAttendanceSessionError(err)
	}

	c.Status(200)
	return nil
}

func getCurrentUserFromContext(c *gin.Context) (int, string, error) {
	userIDValue, ok := c.Get("user_id")
	if !ok {
		return 0, "", httpx.BadRequest("invalid user id")
	}

	userId, ok := userIDValue.(int)
	if !ok {
		return 0, "", httpx.BadRequest("invalid user id")
	}

	claimsValue, ok := c.Get("claims")
	if !ok {
		return 0, "", httpx.Unauthorized("claims missing")
	}

	claims, ok := claimsValue.(map[string]interface{})
	if !ok {
		return 0, "", httpx.Unauthorized("claims missing")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", httpx.Unauthorized("role missing")
	}

	return userId, role, nil
}

func mapAttendanceSessionError(err error) error {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return httpx.NotFound("resource")
	case errors.Is(err, apperrors.ErrUnauthorized):
		return httpx.Unauthorized("insufficient permissions")
	default:
		if err != nil && (strings.Contains(err.Error(), "invalid date") || strings.Contains(err.Error(), "entries are required")) {
			return httpx.BadRequest(err.Error())
		}
		if err != nil && (strings.Contains(err.Error(), "does not belong") || strings.Contains(err.Error(), "already published") || strings.Contains(err.Error(), "department is not set")) {
			return httpx.BadRequest(err.Error())
		}
		return httpx.InternalError(err)
	}
}
