
-- +migrate Up
CREATE TABLE IF NOT EXISTS rockets (
    id                VARCHAR(36) NOT NULL,
    class             VARCHAR(26) COLLATE utf8mb4_unicode_ci NOT NULL,
    launch_speed      VARCHAR(26) COLLATE utf8mb4_unicode_ci NOT NULL,
    mission           VARCHAR(26) COLLATE utf8mb4_unicode_ci NOT NULL,

    PRIMARY KEY (id)
    );

ALTER TABLE rockets ADD INDEX class_idx (class);
-- +migrate Down
