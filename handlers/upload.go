package handlers

import (
	"crypto/md5"
	"fmt"
	"log"
	"monitor-service/db"
	"monitor-service/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadChange 上传变化记录
// @Summary      上传变化记录
// @Description  监控工具上传检测到的JS或内容变化
// @Tags         变化记录
// @Accept       json
// @Produce      json
// @Param        request body models.UploadChangeRequest true "上传数据"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /upload [post]
func UploadChange(c *gin.Context) {
	var req models.UploadChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找或创建站点
	var siteID string
	row := db.DB.QueryRow("SELECT id FROM monitored_sites WHERE url = ?", req.SiteURL)
	err := row.Scan(&siteID)
	if err != nil { // 站点不存在，创建新站点
		siteID = uuid.New().String()
		siteName := ""
		if req.SiteName != nil {
			siteName = *req.SiteName
		}
		_, err = db.DB.Exec(
			`INSERT INTO monitored_sites (id, url, name, created_at, is_active) 
             VALUES (?, ?, ?, ?, ?)`,
			siteID, req.SiteURL, siteName, time.Now(), 1,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create site: " + err.Error()})
			return
		}
	}

	// 插入变化记录
	recordID := uuid.New().String()
	_, err = db.DB.Exec(
		`INSERT INTO change_records 
         (id, site_id, change_type, file_path, change_diff, snapshot_hash, detected_at) 
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
		recordID, siteID, req.ChangeType, req.FilePath, req.ChangeDiff, req.SnapshotHash, time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save change record: " + err.Error()})
		return
	}

	// 如果有原始内容，保存到 file_contents 表
	if req.Content != nil && *req.Content != "" {
		contentHash := fmt.Sprintf("%x", md5.Sum([]byte(*req.Content)))
		contentID := uuid.New().String()
		_, err = db.DB.Exec(
			`INSERT INTO file_contents (id, site_id, file_path, content_hash, content, created_at, change_record_id) 
             VALUES (?, ?, ?, ?, ?, ?, ?)`,
			contentID, siteID, req.FilePath, contentHash, *req.Content, time.Now(), recordID,
		)
		if err != nil {
			log.Printf("Failed to save file content: %v", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "change recorded successfully",
		"site_id":   siteID,
		"record_id": recordID,
	})
}
