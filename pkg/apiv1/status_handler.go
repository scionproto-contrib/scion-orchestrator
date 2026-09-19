// Copyright 2026 OVGU Magdeburg
// SPDX-License-Identifier: Apache-2.0

package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netsys-lab/scion-orchestrator/pkg/metrics"
)

func AddStatusHandler(eng *gin.RouterGroup) {
	eng.GET("status", func(c *gin.Context) {
		c.JSON(http.StatusOK, metrics.Status)
	})
}
