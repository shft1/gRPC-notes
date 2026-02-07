package v1

import "net/http"

func (nh *NoteHandler) UploadMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metricsTotal, err := nh.noteGW.UploadMetrics(ctx)

	if err != nil {
		writeResponse(nh.log, w, 0, nil, err)
		return
	}
	data := map[string]int64{"total": metricsTotal}
	writeResponse(nh.log, w, http.StatusOK, data, nil)
}
