package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	"url-shrotener/internal/api"
	"url-shrotener/internal/storage"
	"url-shrotener/tools"
)

type RequestURL struct {
	URL string `json:"url" validate:"required,url"`
}

type ResponseAlias struct {
	api.Response
	Alias string `json:"alias,omitempty"`
}

func Save(log *slog.Logger, urlStorage storage.URLStorage, aliasLen int, aliasAlpha string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.http.save.Save"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestURL

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

		var alias string
		var success bool
		// TODO: add async
		for !success {
			alias = tools.NewRandomString(aliasLen, aliasAlpha)

			err = urlStorage.SaveURL(req.URL, alias)

			if errors.Is(err, storage.ErrURLExists) {
				alias, err = urlStorage.GetAlias(req.URL)
				if err != nil {
					log.Error("failed to get existing url alias", tools.LogAttr("error", err.Error()))

					render.JSON(w, r, api.Error(err.Error()))
					return
				}
				break
			}

			if errors.Is(err, storage.AliasExists) {
				continue
			}

			success = true
		}

		if err != nil {
			log.Error("failed to add url", tools.LogAttr("error", err.Error()))

			render.JSON(w, r, api.Error("failed to add url"))
			return
		}

		log.Info("url added")

		if r.TLS != nil {
			render.JSON(w, r, ResponseAlias{
				Response: api.OK(),
				Alias:    "https://" + r.Host + "/" + alias,
			})
		} else {
			render.JSON(w, r, ResponseAlias{
				Response: api.OK(),
				Alias:    "http://" + r.Host + "/" + alias,
			})
		}
	}
}
