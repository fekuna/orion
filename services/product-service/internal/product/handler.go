package product

import (
	"errors"
	"net/http"

	"github.com/fekuna/orion/pkg/logger"
	"github.com/fekuna/orion/pkg/response"
	"github.com/fekuna/orion/services/product-service/internal/httputil"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for the product domain.
// It holds no *zap.Logger — logging is done via the request-scoped logger
// injected into every context by the server's loggerMiddleware.
type Handler struct {
	uc UseCase
}

// NewHandler creates a new product HTTP handler.
func NewHandler(uc UseCase) *Handler {
	return &Handler{uc: uc}
}

// RegisterRoutes mounts all product routes onto the given Echo group.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.list)
	g.GET("/:id", h.getByID)
	g.POST("", h.create)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *Handler) list(c echo.Context) error {
	filter := Filter{
		Page:  httputil.ParseIntQuery(c, "page", 1),
		Limit: httputil.ParseIntQuery(c, "limit", 20),
		Name:  c.QueryParam("name"),
	}

	products, total, err := h.uc.GetProducts(c.Request().Context(), filter)
	if err != nil {
		logger.FromContext(c.Request().Context()).Error("list products failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, response.SuccessWithMeta(
		"products retrieved",
		toListData(products),
		response.NewMeta(total, filter.Page, filter.Limit),
	))
}

func (h *Handler) getByID(c echo.Context) error {
	id, err := httputil.ParseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	p, err := h.uc.GetProductByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "product")
		}
		logger.FromContext(c.Request().Context()).Error("get product by id failed",
			zap.String("product_id", id.String()),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, response.Success("product retrieved", toResponse(p)))
}

func (h *Handler) create(c echo.Context) error {
	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	p, err := h.uc.CreateProduct(c.Request().Context(), req)
	if err != nil {
		logger.FromContext(c.Request().Context()).Error("create product failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, response.Created(toResponse(p)))
}

func (h *Handler) update(c echo.Context) error {
	id, err := httputil.ParseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	var req UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	p, err := h.uc.UpdateProduct(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "product")
		}
		logger.FromContext(c.Request().Context()).Error("update product failed",
			zap.String("product_id", id.String()),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, response.Success("product updated", toResponse(p)))
}

func (h *Handler) delete(c echo.Context) error {
	id, err := httputil.ParseUUIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.uc.DeleteProduct(c.Request().Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "product")
		}
		logger.FromContext(c.Request().Context()).Error("delete product failed",
			zap.String("product_id", id.String()),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
