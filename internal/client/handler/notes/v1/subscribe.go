package v1

import (
	"net/http"
)

func (nh *NoteHandler) SubscribeToEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errChan := make(chan error)

	go nh.noteGW.SubscribeToEvents(nh.sysCtx, errChan)

	select {
	case <-ctx.Done():
		return
	case err := <-errChan:
		writeResponse(nh.log, w, http.StatusOK, nil, err)
	}
}
