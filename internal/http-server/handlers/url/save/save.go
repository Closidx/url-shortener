package save

import (
	"errors"
	resp "github/closidx/url-shortener/internal/lib/api/response"
	"github/closidx/url-shortener/internal/lib/logger/sl"
	"github/closidx/url-shortener/internal/lib/random"
	"github/closidx/url-shortener/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	CheckAliasExists(alias string) (bool, error)
}

func New(log *slog.Logger, urlSave URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "hanlders.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failet to decode request bosy", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias, err := getPreparedAlias(req.Alias, urlSave)
		if err != nil {
            log.Error("failed to get prepared alias", sl.Err(err))

            render.JSON(w, r, resp.Error("failed to get prepared alias"))

			return
        }

		id, err := urlSave.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exist", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exist"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		ResponseOK(w, r, alias)
	}
}

func getPreparedAlias(alias string, urlSave URLSaver) (string, error) {
	if alias != "" {
		return alias, nil
	}

	for {
		alias = random.NewRandomString(aliasLength)

		exist, err := urlSave.CheckAliasExists(alias)
		if err != nil {
            return "", err
        }

		if !exist {
			return alias, nil
		}
	}
}

func ResponseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
