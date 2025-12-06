CREATE TABLE `us_doushen_vip`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `vip_type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'vip类型 1.豆神会员（基础版） 2.豆神会员（典藏版） 3.豆神超级会员 4.豆神会员（特别版）',
    `level` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '等级',
    `from_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '权益开始时间',
    `end_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '权益结束时间',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态 1.有效 2.失效',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='豆神VIP表';