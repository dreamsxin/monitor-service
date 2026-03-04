-- 1. 站点信息表
CREATE TABLE IF NOT EXISTS monitored_sites (
    id VARCHAR PRIMARY KEY,                 -- UUID 字符串
    url VARCHAR NOT NULL,
    name VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- 2. 变化记录表（只记录差异，不存全文）
CREATE TABLE IF NOT EXISTS change_records (
    id VARCHAR PRIMARY KEY,
    site_id VARCHAR NOT NULL,                -- 关联 monitored_sites.id
    change_type VARCHAR NOT NULL,             -- 'js' 或 'content'
    file_path VARCHAR,                        -- 当 change_type='js' 时，记录 JS 文件路径
    change_diff TEXT,                          -- 差异内容
    snapshot_hash VARCHAR,                     -- 快照 MD5 哈希
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. 文件内容快照表（存储原始内容，用于调试）
CREATE TABLE IF NOT EXISTS file_contents (
    id VARCHAR PRIMARY KEY,
    site_id VARCHAR NOT NULL,
    file_path VARCHAR,                         -- 文件路径或页面 URL
    content_hash VARCHAR NOT NULL,              -- MD5 哈希
    content TEXT NOT NULL,                       -- 原始内容
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    change_record_id VARCHAR                     -- 可选，关联 change_records.id
);

-- 4. Hook 脚本库
CREATE TABLE IF NOT EXISTS hook_scripts (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    script_type VARCHAR NOT NULL,                 -- 'override', 'content_script', 'snippet'
    script_content TEXT NOT NULL,
    version VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. 站点 Hook 配置（多对多关联）
CREATE TABLE IF NOT EXISTS site_hooks (
    id VARCHAR PRIMARY KEY,
    site_id VARCHAR NOT NULL,
    script_id VARCHAR NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    config JSON,                                   -- 配置参数（JSON 格式）
    priority INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);