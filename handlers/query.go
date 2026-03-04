package handlers

import (
	"monitor-service/db"
	"monitor-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListSites 获取所有站点（支持分页）
// @Summary      获取所有站点
// @Description  返回所有被监控的站点列表，支持分页
// @Tags         站点管理
// @Produce      json
// @Param        page query int false "页码，默认1" minimum(1)
// @Param        page_size query int false "每页数量，默认20，最大100" minimum(1) maximum(100)
// @Success      200  {object}  models.SitesPageResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /sites [get]
func ListSites(c *gin.Context) {
	// 解析分页参数
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val >= 1 {
			page = val
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
			return
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val >= 1 && val <= 100 {
			pageSize = val
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size parameter (must be 1-100)"})
			return
		}
	}

	offset := (page - 1) * pageSize

	// 查询总记录数
	var total int64
	err := db.DB.QueryRow("SELECT COUNT(*) FROM monitored_sites").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count sites: " + err.Error()})
		return
	}

	// 查询分页数据
	rows, err := db.DB.Query(
		`SELECT id, url, name, created_at, is_active 
         FROM monitored_sites 
         ORDER BY created_at DESC 
         LIMIT ? OFFSET ?`,
		pageSize, offset,
	)
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

	// 返回分页响应
	c.JSON(http.StatusOK, models.SitesPageResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Sites:    sites,
	})
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
