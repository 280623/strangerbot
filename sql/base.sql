CREATE TABLE `reports` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `reporter_id` bigint(20) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `report` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `reporter_id` (`reporter_id`),
  CONSTRAINT `reports_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `reports_ibfk_2` FOREIGN KEY (`reporter_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `users` (
 `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
 `chat_id` int(11) NOT NULL,
 `last_activity` datetime NOT NULL,
 `match_chat_id` int(11) DEFAULT NULL,
 `available` tinyint(1) NOT NULL DEFAULT '1',
 `register_date` datetime DEFAULT NULL,
 `previous_match` int(11) DEFAULT NULL,
 `allow_pictures` tinyint(1) NOT NULL,
 `banned_until` datetime DEFAULT NULL,
 `gender` tinyint(4) unsigned DEFAULT '0',
 `tags` varchar(128) DEFAULT '',
 `match_mode` tinyint(4) unsigned DEFAULT '0',
 PRIMARY KEY (`id`),
 UNIQUE KEY `chat_id` (`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `bot_question`
(
    `id`                   int(10) unsigned    NOT NULL AUTO_INCREMENT,
    `scene_type`           tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'scene type 1.profile 2.matching',
    `helper_title`         varchar(512)        NOT NULL DEFAULT '' COMMENT 'question helper title',
    `title`                varchar(128)        NOT NULL DEFAULT '' COMMENT 'question',
    `helper_text`          varchar(2048)       NOT NULL DEFAULT '' COMMENT 'helper text',
    `frontend_type`        tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT 'frontend type 1.select 2.multi select',
    `sort`                 int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'sort order by sort asc',
    `matching_mode`        tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT 'matching mode 1.question by question',
    `matching_question_id` int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'if scene is matching , it is origin question id, if 0 it will not support matching',
    `is_del`               tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'is del 0.false 1.true',
    `deleted_time`         int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'deleted time',
    `create_time`          int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'create time',
    `modify_time`          int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'update time',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_scene_type` (`scene_type`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = DYNAMIC COMMENT ='bot question table';

CREATE TABLE IF NOT EXISTS `bot_option`
(
    `id`              int(10) unsigned    NOT NULL AUTO_INCREMENT,
    `question_id`     int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'question id it form form_question table',
    `option_type`     tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT 'option type 1.value option',
    `label`           varchar(512)        NOT NULL DEFAULT '' COMMENT 'option label',
    `value`           varchar(512)        NOT NULL DEFAULT '' COMMENT 'option value',
    `is_matching_any` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'question is scene type matching,is matching any option? 0.false 1.true',
    `sort`            int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'sort order by sort asc',
    `is_del`          tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'is del 0.false 1.true',
    `deleted_time`    int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'deleted time',
    `create_time`     int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'create time',
    `modify_time`     int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'update time',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unq_question_id_option_value` (`question_id`, `value`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = DYNAMIC COMMENT ='bot select option table';

CREATE TABLE IF NOT EXISTS `bot_user_question_data`
(
    `id`           int(10) unsigned    NOT NULL AUTO_INCREMENT,
    `chat_id`      int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'filler unique id',
    `question_id`  int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'question id',
    `option_id`    int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'option id',
    `value`        varchar(512)        NOT NULL DEFAULT '' COMMENT 'option value',
    `is_del`       tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'is del 0.false 1.true',
    `deleted_time` int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'deleted time',
    `create_time`  int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'create time',
    `modify_time`  int(10) unsigned    NOT NULL DEFAULT '0' COMMENT 'update time',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `unq_chat_id_question_id_option_id` (`option_id`, `question_id`, `chat_id`) USING BTREE
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    ROW_FORMAT = DYNAMIC COMMENT ='bot user commit question data table';