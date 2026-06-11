package qr

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/adeelkhan/qr-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type generateRequest struct {
	InputText string `json:"input_text" binding:"required"`
	Title     string `json:"title"`
}

func (h *Handler) Generate(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	var req generateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code, err := h.svc.Generate(userID, req.InputText, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate QR code"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":           code.ID,
		"image_base64": base64.StdEncoding.EncodeToString(code.ImageData),
		"download_url": fmt.Sprintf("/api/v1/qr/%s/download", code.ID),
	})
}

func (h *Handler) List(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	codes, err := h.svc.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list QR codes"})
		return
	}
	type item struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		InputText string    `json:"input_text"`
		CreatedAt string    `json:"created_at"`
	}
	result := make([]item, len(codes))
	for i, code := range codes {
		result[i] = item{
			ID:        code.ID,
			Title:     code.Title,
			InputText: code.InputText,
			CreatedAt: code.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Get(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	codeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	code, err := h.svc.Get(userID, codeID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err == ErrForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           code.ID,
		"title":        code.Title,
		"input_text":   code.InputText,
		"image_base64": base64.StdEncoding.EncodeToString(code.ImageData),
		"created_at":   code.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *Handler) Delete(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	codeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(userID, codeID); err == ErrForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	} else if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Download(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	codeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	code, err := h.svc.Get(userID, codeID)
	if err == ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err == ErrForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="qr-%s.png"`, code.ID))
	c.Data(http.StatusOK, "image/png", code.ImageData)
}

// RegisterRoutes wires QR routes onto the given router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, jwtMiddleware gin.HandlerFunc) {
	protected := rg.Group("/qr")
	protected.Use(jwtMiddleware)
	protected.POST("/generate", h.Generate)
	protected.GET("", h.List)
	protected.GET("/:id", h.Get)
	protected.DELETE("/:id", h.Delete)
	protected.GET("/:id/download", h.Download)
}
