CREATE TABLE IF NOT EXISTS `auth` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    -- 用户名，需要唯一
    `username` TEXT NOT NULL UNIQUE,
    -- 密码
    `password` TEXT NOT NULL,
    -- 允许用户使用的端口，默认为 0
    `port` INTEGER DEFAULT 0 NOT NULL,
    -- 用户创建时间
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- 最后更新时间
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- 用户是否活跃
    `active` BOOLEAN DEFAULT 1,
    -- 备注
    `remark` TEXT
);


CREATE TRIGGER IF NOT EXISTS `trig_auth_set_updated_at` 
AFTER UPDATE
ON `auth`
FOR EACH ROW
BEGIN
    UPDATE `auth` SET `updated_at` = CURRENT_TIMESTAMP 
    WHERE  `id` = old.id;
END;
