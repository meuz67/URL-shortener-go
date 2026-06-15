package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"

	"url-shortener/internal/storage"
)

type URLHandler struct {
	baseURL string
	store   storage.URLStore
}

type createShortURLRequest struct {
	URL string `json:"url"`
}

type createShortURLResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func NewURLHandler(baseURL string, store storage.URLStore) *URLHandler {
	return &URLHandler{
		baseURL: strings.TrimRight(baseURL, "/"),
		store:   store,
	}
}

func (h *URLHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *URLHandler) CreateShortURL(c echo.Context) error {
	var request createShortURLRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	longURL := strings.TrimSpace(request.URL)
	if !isValidHTTPURL(longURL) {
		return echo.NewHTTPError(http.StatusBadRequest, "url must be a valid http or https URL")
	}

	item, err := h.store.Save(c.Request().Context(), longURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create short URL")
	}

	return c.JSON(http.StatusCreated, createShortURLResponse{
		Code:     item.Code,
		ShortURL: h.baseURL + "/" + item.Code,
		LongURL:  item.LongURL,
	})
}

func (h *URLHandler) Redirect(c echo.Context) error {
	code := strings.TrimSpace(c.Param("code"))
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "code is required")
	}

	item, err := h.store.FindByCode(c.Request().Context(), code)
	if errors.Is(err, storage.ErrURLNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "short URL not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to resolve short URL")
	}

	return c.Redirect(http.StatusFound, item.LongURL)
}

func isValidHTTPURL(value string) bool {
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}

	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
