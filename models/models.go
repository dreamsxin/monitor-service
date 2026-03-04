package models

import (
	"database/sql"
	"time"
)

type MonitoredSite struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// SitesPageResponse 站点列表分页响应
type SitesPageResponse struct {
	Total    int64           `json:"total"`     // 总记录数
	Page     int             `json:"page"`      // 当前页码
	PageSize int             `json:"page_size"` // 每页数量
	Sites    []MonitoredSite `json:"sites"`     // 站点列表
}

type ChangeRecord struct {
	ID           string    `json:"id"`
	SiteID       string    `json:"site_id"`
	ChangeType   string    `json:"change_type"`
	FilePath     *string   `json:"file_path,omitempty"`
	ChangeDiff   string    `json:"change_diff"`
	SnapshotHash string    `json:"snapshot_hash"`
	DetectedAt   time.Time `json:"detected_at"`
}

type FileContent struct {
	ID             string    `json:"id"`
	SiteID         string    `json:"site_id"`
	FilePath       *string   `json:"file_path,omitempty"`
	ContentHash    string    `json:"content_hash"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	ChangeRecordID *string   `json:"change_record_id,omitempty"`
}

type HookScript struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ScriptType    string    `json:"script_type"`
	ScriptContent string    `json:"script_content"`
	Version       string    `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
}

type SiteHook struct {
	ID        string         `json:"id"`
	SiteID    string         `json:"site_id"`
	ScriptID  string         `json:"script_id"`
	Enabled   bool           `json:"enabled"`
	Config    sql.NullString `json:"config,omitempty"`
	Priority  int            `json:"priority"`
	CreatedAt time.Time      `json:"created_at"`
}

type UploadChangeRequest struct {
	SiteURL      string  `json:"site_url" binding:"required"`
	SiteName     *string `json:"site_name"`
	ChangeType   string  `json:"change_type" binding:"required,oneof=js content"`
	FilePath     *string `json:"file_path"`
	ChangeDiff   string  `json:"change_diff" binding:"required"`
	SnapshotHash string  `json:"snapshot_hash" binding:"required"`
	Content      *string `json:"content"`
}

type SiteHooksResponse struct {
	Type    string                 `json:"type"`
	Content string                 `json:"content"`
	Config  map[string]interface{} `json:"config,omitempty"`
}
