CREATE TABLE `lv_user` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 Id',
    `stream_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '直播ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `user_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
    `im_user_id` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'IM用户ID',
    `role_type` tinyint NOT NULL DEFAULT '0' COMMENT '用户角色类型：1 主讲 2 学生',
    `is_mute` tinyint NOT NULL DEFAULT '2' COMMENT '是否禁言（1是，2否）',
    `mute_end_at` bigint NOT NULL DEFAULT '0' COMMENT '禁言结束时间戳',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_stream_user` (`stream_id`, `user_id`) COMMENT '直播ID和用户ID联合唯一索引'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '直播参与人员表';


CREATE TABLE `lv_chat_room`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 Id',
    `stream_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '直播ID',
    `room_id` varchar(50) NOT NULL COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '聊天房间号',
    `room_name` varchar(100) NOT NULL COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '聊天室名称',
    `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '类型：1.主聊天室 2.讲师聊天室',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '聊天室状态：1.正常 2.关闭',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_stream_id` (`stream_id`) COMMENT '直播ID索引'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '直播聊天室表';


CREATE TABLE `lv_signal_log`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 Id',
    `stream_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '直播ID',
    `chat_room_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '聊天室ID',
    `event` varchar(100) NOT NULL COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '事件名称',
    `event_type` varchar(100) NOT NULL COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '事件细分分类',
    `recv_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '信令接收方：1.全员 2.个人',
    `to_account` varchar(100) NOT NULL COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '信令接收方账号',
    `content` text NOT NULL COLLATE utf8mb4_general_ci COMMENT '信令内容（JSON 格式）',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_stream_id` (`stream_id`) COMMENT '直播ID索引',
    KEY `idx_chat_room_id` (`chat_room_id`) COMMENT '聊天室ID索引',
    KEY `idx_event` (`event`) COMMENT '事件名称索引'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '直播聊天室信令记录表';