package controllers

import (
    "log/slog"
    "net/http"
    "strconv"

    "github.com/go-chi/chi"
    "github.com/go-chi/render"

    "take-home-challenge/configuration"
    "take-home-challenge/coordinators"
    "take-home-challenge/middleware"
    "take-home-challenge/models"
    "take-home-challenge/services/api/controllers/payloads"
)

type MessagesController struct {
	coordinator coordinators.MessagesCoordinator
}

func NewMessagesController(conf configuration.Config) *MessagesController {
	return &MessagesController{
		coordinator: conf.MessagesCoordinator,
	}
}

func (m *MessagesController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Post("/", m.create)
	r.Post("/delete", m.delete)
	r.Get("/{messageID}", m.read)

	return r
}

func (m *MessagesController) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := ctx.Value(middleware.KeyLogger).(*slog.Logger)

	res := Response{}

	data := &models.Message{}
	if err := render.Bind(r, data); err != nil {
		res.BadRequest(w, r, "invalid request body")
		return
	}

	if err := m.coordinator.Create(data); err != nil {
		if coordinators.IsValidationError(err) {
			res.BadRequest(w, r, err.Error())
			return
		}

		l.ErrorContext(ctx, "Failed to create message",
			"userIDs", data.UserIDs,
			"error", err,
		)
		res.InternalServerError(w, r, "internal server error")
		return
	}

	res.JSON(w, r, data, http.StatusCreated)
}

func (m *MessagesController) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := ctx.Value(middleware.KeyLogger).(*slog.Logger)

	res := Response{}

	data := &payloads.MessagesMarkDeleted{}
	if err := render.Bind(r, data); err != nil {
		res.BadRequest(w, r, "invalid request body")
		return
	}

	if err := m.coordinator.MarkDeleted(data.DeletedWhen, data.IDs); err != nil {
		if coordinators.IsValidationError(err) {
			res.BadRequest(w, r, err.Error())
			return
		}
		
		l.ErrorContext(ctx, "Failed to mark messages as deleted",
			"messageIDs", data.IDs,
			"deletedWhen", data.DeletedWhen,
			"error", err,
		)
		res.InternalServerError(w, r, "internal server error")
		return
	}

	res.JSON(w, r, nil, http.StatusOK)
}

func (m *MessagesController) read(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    l := ctx.Value(middleware.KeyLogger).(*slog.Logger)

    res := Response{}

    messageIDParam := chi.URLParam(r, "messageID")
    if messageIDParam == "" {
        res.BadRequest(w, r, "message ID is required")
        return
    }

    messageID, err := strconv.Atoi(messageIDParam)
    if err != nil {
        res.BadRequest(w, r, "invalid message ID format")
        return
    }

    if messageID <= 0 {
        res.BadRequest(w, r, "message ID must be a positive integer")
        return
    }

    message, err := m.coordinator.Read(messageID)
    if err != nil {
        if coordinators.IsRecordNotFoundError(err) {
            res.NotFound(w, r, "message not found")
            return
        }
        
        l.ErrorContext(ctx, "Failed to read message",
            "messageID", messageID,
            "error", err,
        )
        res.InternalServerError(w, r, "internal server error")
        return
    }

    res.JSON(w, r, message, http.StatusOK)
}
