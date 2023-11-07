package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	"net/url"
	"url-shrotener/internal/api"
	"url-shrotener/internal/storage"
	"url-shrotener/tools"
)

type RequestAlias struct {
	Alias string `json:"alias" validate:"url,required"`
}

type ResponseURL struct {
	api.Response
	URL string `json:"url"`
}

func Get(log *slog.Logger, urlStorage storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.http.get.Get"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestAlias

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", tools.LogAttr("error", err.Error()))

			render.JSON(w, r, api.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", tools.LogAttr("error", err.Error()))

			render.JSON(w, r, api.ValidationError(validateErr))
			return
		}

		u, _ := url.Parse(req.Alias)
		alias := u.Path[1:]
		resURL, err := urlStorage.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", req.Alias)

			render.JSON(w, r, api.Error("not found"))
			return
		}

		if err != nil {
			log.Error("failed to get url", tools.LogAttr("error", err.Error()))

			render.JSON(w, r, api.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", resURL))

		render.JSON(w, r, ResponseURL{
			Response: api.OK(),
			URL:      resURL,
		})
	}
}
