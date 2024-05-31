package delete_url

import (
	"errors"
	"github/closidx/url-shortener/internal/http-server/handlers/url/save"
	resp "github/closidx/url-shortener/internal/lib/api/response"
	"github/closidx/url-shortener/internal/lib/logger/sl"
	"github/closidx/url-shortener/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UrlDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDelete UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "hanlders.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDelete.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", alias)

				render.JSON(w, r, resp.Error("url not found"))

				return
			}

			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("URL deleted", slog.String("alias", alias))
		save.ResponseOK(w, r, alias)
	}
}
