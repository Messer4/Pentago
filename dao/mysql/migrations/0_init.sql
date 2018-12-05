CREATE TABLE `users` (
  `usr_id`                   INT(11)             NOT NULL AUTO_INCREMENT,
  `usr_email`                VARCHAR(255)        NOT NULL,
  `usr_password`        VARCHAR(255)                NOT NULL,
  PRIMARY KEY (`usr_id`),
  UNIQUE KEY `users_usr_email_uindex` (`usr_email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `game` (
  `gm_id`                   INT(11)             NOT NULL AUTO_INCREMENT,
  `usr_id_1`               INT(11)       NOT NULL,
  `usr_id_2`        INT(11)              NOT NULL,
  `result`          VARCHAR(10),
    PRIMARY KEY (`gm_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;