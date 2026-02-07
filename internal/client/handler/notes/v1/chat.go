package v1

import "net/http"

func (nh *NoteHandler) Chat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	errChan := make(chan error, 1)

	go nh.noteGW.Chat(nh.sysCtx, errChan)

	select {
	case <-ctx.Done():
		return
	case err := <-errChan:
		writeResponse(nh.log, w, http.StatusOK, nil, err)
	}
}
