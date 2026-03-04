package handlers

import (
	"monitor-service/db"
	"monitor-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFileContent 获取文件内容快照
// @Summary      获取文件内容快照
// @Description  根据站点ID、文件路径或哈希值查询历史文件内容
// @Tags         内容快照
// @Produce      json
// @Param        site_id query string true "站点ID"
// @Param        file_path query string false "文件路径"
// @Param        hash query string false "内容哈希值(MD5)"
// @Success      200  {array}   models.FileContent
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /content [get]
func GetFileContent(c *gin.Context) {
	siteID := c.Query("site_id")
	filePath := c.Query("file_path")
	hash := c.Query("hash")

	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing site_id"})
		return
	}

	query := `SELECT id, site_id, file_path, content_hash, content, created_at, change_record_id 
              FROM file_contents WHERE site_id = ?`
	args := []interface{}{siteID}

	if filePath != "" {
		query += " AND file_path = ?"
		args = append(args, filePath)
	}
	if hash != "" {
		query += " AND content_hash = ?"
		args = append(args, hash)
	}
	query += " ORDER BY created_at DESC LIMIT 50"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var contents []models.FileContent
	for rows.Next() {
		var fc models.FileContent
		err := rows.Scan(&fc.ID, &fc.SiteID, &fc.FilePath, &fc.ContentHash, &fc.Content, &fc.CreatedAt, &fc.ChangeRecordID)
		if err != nil {
			continue
		}
		contents = append(contents, fc)
	}
	c.JSON(http.StatusOK, contents)
}
