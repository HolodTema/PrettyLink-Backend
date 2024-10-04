package save

import (
	"PrettyLinkBackend/internal/lib/api/response"
	"PrettyLinkBackend/internal/lib/logger/sl"
	"PrettyLinkBackend/internal/lib/random"
	"PrettyLinkBackend/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitepty"`
}

// omitempty can remove fields of "" from the struct
type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"` //alias is optional param
}

// TODO: maybe move aliasLength to config
const aliasLength = 6

type URLSaver interface {
	Save(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		//extend logger with const op
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		//unmarshall input request
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			//chi is Gin analog

			//return to finish request handling. Otherwise, render.JSON will
			//handle request endless...
			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		//validate struct Request:
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			//log the error just like err
			log.Error("invalid request", sl.Err(err))

			//log error just validated error
			render.JSON(w, r, response.ValidateError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.Save(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to save url"))

			return
		}

		//log that url saved successful with url id
		log.Info("url added", slog.Int64("id", id))

		render.JSON(
			w,
			r,
			Response{
				Response: response.OK(),
				Alias:    alias,
			})
	}
}
