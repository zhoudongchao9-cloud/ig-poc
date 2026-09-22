package api

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type DiagnosticsRequest struct {
	Probe string `json:"probe"`
}

// DiagnosticsHandler runs a named probe from the ops catalog and returns its
// combined output. Used by on-call runbooks instead of manual SSH sessions.
func DiagnosticsHandler(c *gin.Context) {
	var input DiagnosticsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid diagnostics request"})
		return
	}

	out, err := exec.Command("/bin/sh", "-c", input.Probe).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: string(out)})
		return
	}

	c.JSON(http.StatusOK, ResponseJSON{Message: string(out)})
}
