package handlers

import (
	"monitor-service/db"
	"monitor-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListSites 获取所有站点
// @Summary      获取所有站点
// @Description  返回所有被监控的站点列表
// @Tags         站点管理
// @Produce      json
// @Success      200  {array}   models.MonitoredSite
// @Failure      500  {object}  map[string]interface{}
// @Router       /sites [get]
func ListSites(c *gin.Context) {
	rows, err := db.DB.Query(`SELECT id, url, name, created_at, is_active FROM monitored_sites ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var sites []models.MonitoredSite
	for rows.Next() {
		var s models.MonitoredSite
		var isActive int
		err := rows.Scan(&s.ID, &s.URL, &s.Name, &s.CreatedAt, &isActive)
		if err != nil {
			continue
		}
		s.IsActive = isActive == 1
		sites = append(sites, s)
	}
	c.JSON(http.StatusOK, sites)
}

// GetChanges 查询变化记录
// @Summary      查询变化记录
// @Description  根据站点ID和变化类型查询历史变化
// @Tags         变化记录
// @Produce      json
// @Param        site_id query string false "站点ID"
// @Param        type query string false "变化类型 (js 或 content)"
// @Success      200  {array}   models.ChangeRecord
// @Failure      500  {object}  map[string]interface{}
// @Router       /changes [get]
func GetChanges(c *gin.Context) {
	siteID := c.Query("site_id")
	changeType := c.Query("type")

	query := `SELECT id, site_id, change_type, file_path, change_diff, snapshot_hash, detected_at 
              FROM change_records WHERE 1=1`
	args := []interface{}{}

	if siteID != "" {
		query += " AND site_id = ?"
		args = append(args, siteID)
	}
	if changeType != "" {
		query += " AND change_type = ?"
		args = append(args, changeType)
	}
	query += " ORDER BY detected_at DESC LIMIT 100"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var changes []models.ChangeRecord
	for rows.Next() {
		var ch models.ChangeRecord
		err := rows.Scan(&ch.ID, &ch.SiteID, &ch.ChangeType, &ch.FilePath, &ch.ChangeDiff, &ch.SnapshotHash, &ch.DetectedAt)
		if err != nil {
			continue
		}
		changes = append(changes, ch)
	}
	c.JSON(http.StatusOK, changes)
}
