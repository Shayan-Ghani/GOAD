package model

import (
	"time"
	// "fmt"

)


type Item struct {
	ID uint
	Name string
	Description string
	CreatedAt time.Time
	ModifiedAt time.Time
	DeletedAt time.Time
	Tags []Tag
}

type Tag struct {
	Items []Item
	ID uint
	Name string
}

type TagItems struct {
	ItemID uint
	TagID uint
}

// -- Generate Tabels
// CREATE TABLE `items` (
//   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
//   `created_at` datetime(3) DEFAULT NULL,
//   `updated_at` datetime(3) DEFAULT NULL,
//   `deleted_at` datetime(3) DEFAULT NULL,
//   `name` varchar(191) NOT NULL,
//   `description` longtext,
//   PRIMARY KEY (`id`),
//   KEY `idx_items_deleted_at` (`deleted_at`)
// );

// CREATE TABLE `tags` (
//   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
//   `created_at` datetime(3) DEFAULT NULL,
//   `updated_at` datetime(3) DEFAULT NULL,
//   `deleted_at` datetime(3) DEFAULT NULL,
//   `name` varchar(191) DEFAULT NULL,
//   PRIMARY KEY (`id`),
//   UNIQUE KEY `idx_tags_name` (`name`),
//   KEY `idx_tags_deleted_at` (`deleted_at`)
// );

// CREATE TABLE `item_tags` (
//   `item_id` bigint unsigned NOT NULL,
//   `tag_id` bigint unsigned NOT NULL,
//   PRIMARY KEY (`item_id`,`tag_id`),
//   KEY `fk_item_tags_tag` (`tag_id`),
//   CONSTRAINT `fk_item_tags_item` FOREIGN KEY (`item_id`) REFERENCES `items` (`id`),
//   CONSTRAINT `fk_item_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
// );