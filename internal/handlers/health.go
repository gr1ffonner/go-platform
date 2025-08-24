package handlers

import (
	httputils "go-platform/pkg/utils/http-utils"
	"net/http"
)

// Health godoc
//
//	@Summary	Health check
//	@Tags		Health
//	@Accept		json
//	@Produce	json
//	@Success	200	{string}	httputils.Status		"ok"
//	@Failure	500	{string}	httputils.ErrorResponse	"internal server error"
//	@Router		/live [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httputils.WriteResponse(w, http.StatusOK, "ok", nil, nil)
}
