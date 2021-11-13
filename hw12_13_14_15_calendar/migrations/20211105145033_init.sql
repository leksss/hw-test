-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `event`
(
    `uuid`       VARCHAR(36)      NOT NULL DEFAULT '',
    `owner_id`   VARCHAR(255)     NOT NULL DEFAULT '',
    `title`      VARCHAR(255)     NOT NULL DEFAULT '',
    `started_at` TIMESTAMP        NULL     DEFAULT NULL,
    `ended_at`   TIMESTAMP        NULL     DEFAULT NULL,
    `text`       VARCHAR(255)     NOT NULL DEFAULT '',
    `notify_for` INT(11) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`uuid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `event`;
-- +goose StatementEnd
