CREATE TABLE `us_session` (
                              `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                              `session` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '会话',
                              `start_time` datetime NOT NULL COMMENT '开始时间',
                              `end_time` datetime NOT NULL COMMENT '结束时间',
                              `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                              PRIMARY KEY (`id`),
                              KEY `idx_session` (`session`) USING BTREE COMMENT '会话索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='会话表';


CREATE TABLE `us_session_user` (
                                   `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                                   `user_id` bigint NOT NULL COMMENT '用户ID',
                                   `session_id` bigint NOT NULL COMMENT '会话ID',
                                   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                   PRIMARY KEY (`id`),
                                   KEY `idx_user_id` (`user_id`) USING BTREE COMMENT 'user_id 索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户会话关系表';




CREATE TABLE `us_session_record` (
                                     `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                                     `session_id` bigint NOT NULL COMMENT '会话ID',
                                     `event_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '事件id',
                                     `event_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '事件类型',
                                     `event_time` int unsigned NOT NULL DEFAULT '0' COMMENT '事件时间戳',
                                     `page` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '事件页面',
                                     `appid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '应用ID',
                                     `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '手机号',
                                     `user_agent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '浏览器信息',
                                     `screen` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '分辨率',
                                     `language` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语言',
                                     `data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '携带数据',
                                     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                     PRIMARY KEY (`id`),
                                     KEY `idx_session_id` (`session_id`) USING BTREE COMMENT '会话id索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='会话记录表';