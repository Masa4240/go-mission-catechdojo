CREATE DATABASE IF NOT EXISTS ca_mission;
USE ca_mission;


CREATE TABLE IF NOT EXISTS rankratio
(
  `ranklevel`         varchar(255),
  `weight`     int
);

INSERT INTO rankratio (ranklevel, weight) VALUES ("SR", 5);
INSERT INTO rankratio (ranklevel, weight) VALUES ("R", 15);
INSERT INTO rankratio (ranklevel, weight) VALUES ("N", 80);

CREATE TABLE IF NOT EXISTS formal_character_list
(
    `id`         int NOT NULL AUTO_INCREMENT,
    `created_at`         datetime(3) DEFAULT NULL,
    `updated_at`         datetime(3) DEFAULT NULL,
    `deleted_at`         datetime(3) DEFAULT NULL,
    `name`     varchar(255),
    `rank`     varchar(255),
    `weight`     int,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS user_characterlist
(
    -- gorm.Model
    `id`         int NOT NULL AUTO_INCREMENT,
    `created_at`         datetime(3) DEFAULT NULL,
    `updated_at`         datetime(3) DEFAULT NULL,
    `deleted_at`         datetime(3) DEFAULT NULL,
    `character_id`         int,
    `user_id`     int,
    PRIMARY KEY (`id`)
);
INSERT INTO user_characterlist (`id`, `character_id`, `user_id`) VALUES (1, 1, 1);

CREATE TABLE IF NOT EXISTS user_list
(
    `id`         int NOT NULL AUTO_INCREMENT,
    `created_at`         datetime(3) DEFAULT NULL,
    `updated_at`         datetime(3) DEFAULT NULL,
    `deleted_at`         datetime(3) DEFAULT NULL,
    `name`         varchar(255),
    PRIMARY KEY (`id`)
);

INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (1, "SR", "Pika-Chu", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (2, "SR", "Mew", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (3, "SR", "Eevee", 1);

INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (4, "R", "Ditto", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (5, "R", "Lapras", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (6, "R", "Pidgeot", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (7, "R", "Starmie", 1);

INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (8, "N", "Zubat", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (9, "N", "Paras", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (10, "N", "Caterpie", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (11, "N", "Weedle", 1);
INSERT INTO formal_character_list (`id`, `rank`, `name`, `weight`) VALUES (12, "N", "Rattata", 2);