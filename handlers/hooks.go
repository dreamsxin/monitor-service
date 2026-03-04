package handlers

import (
	"database/sql"
	"encoding/json"
	"monitor-service/db"
	"monitor-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSiteHooks 获取站点 Hook 脚本
// @Summary      获取站点 Hook 脚本
// @Description  返回指定站点所有启用的 Hook 脚本（用于客户端注入）
// @Tags         Hook脚本
// @Produce      json
// @Param        url query string true "站点完整URL"
// @Success      200  {array}   models.SiteHooksResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /site-hooks [get]
func GetSiteHooks(c *gin.Context) {
	siteURL := c.Query("url")
	if siteURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing url"})
		return
	}

	// 查询站点ID
	var siteID string
	row := db.DB.QueryRow("SELECT id FROM monitored_sites WHERE url = ?", siteURL)
	if err := row.Scan(&siteID); err != nil {
		// 站点不存在，返回空列表
		c.JSON(http.StatusOK, []models.SiteHooksResponse{})
		return
	}

	// 查询该站点的所有启用脚本，按优先级排序
	rows, err := db.DB.Query(`
        SELECT hs.script_type, hs.script_content, sh.config
        FROM site_hooks sh
        JOIN hook_scripts hs ON sh.script_id = hs.id
        WHERE sh.site_id = ? AND sh.enabled = 1
        ORDER BY sh.priority DESC
    `, siteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var hooks []models.SiteHooksResponse
	for rows.Next() {
		var scriptType, scriptContent string
		var configJSON sql.NullString
		rows.Scan(&scriptType, &scriptContent, &configJSON)

		hook := models.SiteHooksResponse{
			Type:    scriptType,
			Content: scriptContent,
		}
		if configJSON.Valid {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(configJSON.String), &config); err == nil {
				hook.Config = config
			}
		}
		hooks = append(hooks, hook)
	}
	c.JSON(http.StatusOK, hooks)
}
