package api

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type DiagnosticsRequest struct {
	Check string `json:"check"`
}

// DiagnosticsHandler runs a named probe from the ops catalog and returns its
// combined output. Used by on-call runbooks instead of manual SSH sessions.
func DiagnosticsHandler(c *gin.Context) {
	var req DiagnosticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid diagnostics request"})
		return
	}

	out, err := exec.Command("/bin/sh", "-c", req.Check).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: string(out)})
		return
	}

	c.JSON(http.StatusOK, ResponseJSON{Message: string(out)})
}
