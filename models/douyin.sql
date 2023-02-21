USE douyin;
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`          bigint(20)   NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(20)   NOT NULL COMMENT '用户ID',
    `username`    varchar(32)  NOT NULL COLLATE utf8mb4_general_ci COMMENT '姓名',
    `password`    varchar(255) NOT NULL COLLATE utf8mb4_general_ci COMMENT '密码',
    `description` varchar(512)          DEFAULT '' COLLATE utf8mb4_general_ci COMMENT '简介',
    `avatar`      varchar(255) NOT NULL DEFAULT '' COLLATE utf8mb4_general_ci COMMENT '头像',
    `bg_image`    varchar(255) NOT NULL DEFAULT '' COLLATE utf8mb4_general_ci COMMENT '背景图片',
    `create_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY `pk_users_id` (`id`),
    UNIQUE KEY `uk_users_user_id` (`user_id`),
    UNIQUE KEY `uk_users_username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '用户表';

DROP TABLE IF EXISTS `user_favor_videos`;
CREATE TABLE `user_favor_videos`
(
    `id`        bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id`   bigint(20) NOT NULL COMMENT '点赞用户ID',
    `video_id`  bigint(20) NOT NULL COMMENT '被点赞视频ID',
    `is_favor` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0表示没有点赞, 1表示已经点赞',
    PRIMARY KEY `pk_favor_id` (`id`),
    UNIQUE KEY `idx_favor_user_id` (`user_id`, `video_id`),
    INDEX `idx_favor_user_follow_id` (`user_id`, `is_favor`),
    INDEX `idx_favor_to_user_follow_id` (`video_id`, `is_favor`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '点赞表';

DROP TABLE IF EXISTS `relations`;
CREATE TABLE `relations`
(
    `id`         bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id`    bigint(20) NOT NULL COMMENT '用户id',
    `to_user_id` bigint(20) NOT NULL COMMENT '关系用户id',
    `is_follow`  tinyint(4) NOT NULL DEFAULT 1 COMMENT '0表示未关注,1表示已经关注',
    PRIMARY KEY `pk_relations_id` (`id`),
    UNIQUE KEY `idx_relations_user_id` (`user_id`, `to_user_id`),
    INDEX `idx_relations_user_follow_id` (`user_id`, `is_follow`),
    INDEX `idx_relations_to_user_follow_id` (`to_user_id`, `is_follow`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '关系表';

DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`
(
    `id`        bigint(20)   NOT NULL AUTO_INCREMENT,
    `video_id`  bigint(20)   NOT NULL COMMENT '视频id',
    `user_id`   bigint(20)   NOT NULL COMMENT '用户id',
    `title`     varchar(255) NOT NULL COLLATE utf8mb4_general_ci COMMENT '视频标题',
    `play_url`  varchar(255) NOT NULL COLLATE utf8mb4_general_ci COMMENT '视频地址',
    `cover_url` varchar(255) NOT NULL COLLATE utf8mb4_general_ci COMMENT '视频封面地址',
    `is_delete` tinyint(4)   NOT NULL DEFAULT 0 COMMENT '0表示未删除,1表示已经删除',
    `create_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
    `update_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY `pk_videos_id` (`id`),
    INDEX `idx_videos_create_at` (`create_at`, `is_delete`),
    INDEX `idx_videos_user_delete_id` (`user_id`, `is_delete`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '视频表';

DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `id`        bigint(20)   NOT NULL AUTO_INCREMENT,
    `user_id`   bigint(20)   NOT NULL COMMENT '用户id',
    `video_id`  bigint(20)   NOT NULL COMMENT '视频id',
    `content`   varchar(512) NOT NULL COLLATE utf8mb4_general_ci COMMENT '评论内容',
    `is_delete` tinyint(4)   NOT NULL DEFAULT 0 COMMENT '0表示未删除,1表示已经删除',
    `create_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY `pk_comments_id` (`id`),
    INDEX `idx_comments_video_id` (`video_id`, `is_delete`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '评论表';
