CREATE TABLE `hospital` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL COMMENT 'The hospital name',
  `display_name` varchar(200) NOT NULL DEFAULT '' COMMENT 'The display name',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `employee` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'The primary key',
  `hospital_id` bigint NOT NULL,
  `username` varchar(50) NOT NULL,
  `first_name` varchar(100) NOT NULL DEFAULT '',
  `last_name` varchar(100) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_name` (`username`),
  KEY `idx_hid` (`hospital_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `task` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'The primary key',
  `hospital_id` bigint NOT NULL,
  `owner_id` bigint NOT NULL,
  `title` varchar(100) NOT NULL COMMENT 'The task title',
  `description` varchar(500) NOT NULL COMMENT 'The task description',
  `priority` varchar(50) NOT NULL COMMENT 'The task priority. Could be one of urgent, hight, low',
  `status` varchar(50) NOT NULL COMMENT 'The task status. Could be one of open, failed, completed',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_hid` (`hospital_id`),
  KEY `idx_oid` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
