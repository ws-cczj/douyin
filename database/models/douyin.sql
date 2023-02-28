USE
douyin;
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`           bigint(20) NOT NULL COMMENT '用户ID',
    `username`          varchar(32)  NOT NULL COMMENT '用户名',
    `password`          varchar(255) NOT NULL COMMENT '用户密码',
    `signature`         varchar(512) NOT NULL DEFAULT '' COLLATE utf8mb4_general_ci COMMENT '个性签名',
    `avatar`            varchar(255) NOT NULL DEFAULT '' COLLATE utf8mb4_general_ci COMMENT '头像',
    `background_image`  varchar(255) NOT NULL DEFAULT '' COLLATE utf8mb4_general_ci COMMENT '背景图片',
    `work_count`        bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户视频数',
    `favor_count`       bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户点赞数',
    `total_favor_count` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户获赞数',
    `follow_count`      bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户关注数',
    `follower_count`    bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户粉丝数',
    `create_at`         datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at`         datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY `pk_users_id` (`id`),
    UNIQUE KEY `uk_users_user_id` (`user_id`),
    UNIQUE KEY `uk_users_username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE utf8mb4_general_ci COMMENT '用户表';

DROP TABLE IF EXISTS `user_favor_videos`;
CREATE TABLE `user_favor_videos`
(
    `id`       bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`  bigint(20) NOT NULL COMMENT '点赞用户ID',
    `video_id` bigint(20) NOT NULL COMMENT '被点赞视频ID',
    `is_favor` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0表示没有点赞, 1表示已经点赞',
    PRIMARY KEY `pk_favor_id` (`id`),
    UNIQUE KEY `uk_favor_user_id` (`user_id`, `video_id`),
    KEY        `idx_favor_user_id` (`user_id`, `is_favor`),
    KEY        `idx_favor_video_id` (`video_id`, `is_favor`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE utf8mb4_general_ci COMMENT '点赞表';

DROP TABLE IF EXISTS `user_relations`;
CREATE TABLE `user_relations`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`    bigint(20) NOT NULL COMMENT '用户id',
    `to_user_id` bigint(20) NOT NULL COMMENT '关系用户id',
    `is_friend`  tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否互相关注,0表示未互相关注,1表示互相关注',
    `is_follow`  tinyint(4) NOT NULL DEFAULT 1 COMMENT '0表示未关注,1表示已经关注',
    `create_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY `pk_relations_id` (`id`),
    UNIQUE KEY `uk_relations_user_id` (`user_id`, `to_user_id`),
    KEY          `idx_relations_to_user_id` (`to_user_id`, `user_id`, `is_follow`),
    KEY          `idx_relations_user_follow_id` (`user_id`, `is_follow`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '关系表';

DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `video_id`      bigint(20) NOT NULL COMMENT '视频id',
    `user_id`       bigint(20) NOT NULL COMMENT '用户id',
    `title`         varchar(255) NOT NULL COMMENT '视频标题',
    `play_url`      varchar(255) NOT NULL COMMENT '视频地址',
    `cover_url`     varchar(255) NOT NULL COMMENT '视频封面地址',
    `favored_count` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '视频获赞数',
    `comment_count` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '视频评论数',
    `is_delete`     tinyint(4) NOT NULL DEFAULT 0 COMMENT '0表示未删除,1表示已经删除',
    `create_at`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
    `update_at`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY `pk_videos_id` (`id`),
    KEY             `idx_videos_create_at` (`create_at`, `is_delete`),
    KEY             `idx_videos_video_delete_id` (`video_id`, `is_delete`),
    KEY             `idx_videos_user_delete_id` (`user_id`, `is_delete`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE utf8mb4_general_ci COMMENT '视频表';

DROP TABLE IF EXISTS `video_comments`;
CREATE TABLE `video_comments`
(
    `id`        bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`   bigint(20) NOT NULL COMMENT '用户id',
    `video_id`  bigint(20) NOT NULL COMMENT '视频id',
    `content`   varchar(512) NOT NULL COMMENT '评论内容',
    `is_delete` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0表示未删除,1表示已经删除',
    `create_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY `pk_comments_id` (`id`),
    KEY         `idx_comments_video_id` (`video_id`, `is_delete`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE utf8mb4_general_ci COMMENT '评论表';
