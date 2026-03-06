package sdmp

import (
	"net/http"
)

// BlobHandler serves message payloads over HTTP using message IDs as
// capability tokens. Knowing the message ID is sufficient to fetch
// the encrypted ciphertext — no authentication required.
//
// Supports Range requests for resumable downloads (critical for large
// file transfers). CDN-cacheable by design.
type BlobHandler struct {
	// Store will be the interface to transit storage once implemented.
}

// ServeHTTP handles GET /message/{id} requests.
func (h *BlobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	messageID := r.PathValue("id")
	if messageID == "" {
		http.Error(w, "message id required", http.StatusBadRequest)
		return
	}

	// TODO: look up message in transit storage by ID, serve the
	// encrypted blob with Content-Type: application/octet-stream.
	// Use http.ServeContent for automatic Range request support.
	http.Error(w, "not yet implemented", http.StatusNotImplemented)
}
