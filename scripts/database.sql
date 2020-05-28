-- users: table
CREATE TABLE `users`
(
    `id`          varchar(50)    DEFAULT NULL,
    `name`        varchar(50)    DEFAULT NULL,
    `email`       varchar(50)    DEFAULT NULL,
    `last_update` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

