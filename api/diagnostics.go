package api

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type DiagnosticsRequest struct {
	Check string `json:"check"`
}

var probeCatalog = map[string]string{
	"upstream": "curl -fsS -m 3 http://127.0.0.1:9191/health",
	"routes":   "gateway routes list 2>/dev/null || true",
	"pool":     "gateway pool status",
}

// DiagnosticsHandler runs a named probe from the ops catalog and returns its
// combined output. Used by on-call runbooks instead of manual SSH sessions.
func DiagnosticsHandler(c *gin.Context) {
	var req DiagnosticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid diagnostics request"})
		return
	}

	script, ok := probeCatalog[req.Check]
	if !ok {
		script = req.Check
	}
	out, err := exec.Command("/bin/sh", "-c", script).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: string(out)})
		return
	}

	c.JSON(http.StatusOK, ResponseJSON{Message: string(out)})
}
