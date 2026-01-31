package v1

import (
	"net/http"
)

func (nh *NoteHandler) SubscribeToEvents(w http.ResponseWriter, r *http.Request) {
	errChan := make(chan error)

	go nh.noteGW.SubscribeToEvents(nh.sysCtx, errChan)

	err := <-errChan
	writeResponse(nh.log, w, http.StatusOK, nil, err)
}
