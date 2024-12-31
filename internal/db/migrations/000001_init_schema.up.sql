CREATE TABLE IF NOT EXISTS items (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(155) NOT NULL,
  `description` mediumtext,
  `is_done` bit(1) DEFAULT b'0',
  `created_at` datetime(3) NOT NULL,
  `modified_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_item_id` (`id`)
);

CREATE TABLE IF NOT EXISTS tags (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_tag_name_index` (`name`)
);

CREATE TABLE IF NOT EXISTS item_tags (
  `item_id` bigint unsigned NOT NULL,
  `tag_id` int unsigned NOT NULL,
  PRIMARY KEY (`item_id`,`tag_id`),
  KEY `idx_item_tag_tag_id` (`tag_id`),
  KEY `idx_item_tag_item_id` (`item_id`) USING BTREE,
  CONSTRAINT `fk_item_id` FOREIGN KEY (`item_id`) REFERENCES `items` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
);