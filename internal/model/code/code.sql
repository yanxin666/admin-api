CREATE TABLE `rc_redeem_code` (
                                  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '兑换码',
                                  `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 1.未兑换 2.已兑换',
                                  `benefits_group_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '权益组表ID',
                                  `batch` int NOT NULL COMMENT '批次 格式20240101 当天多个批次就以此类推 202401011',
                                  `valid_date` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '兑换有效期',
                                  `source` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '来源 1.外部导入 2.内容生成',
                                  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
                                  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
                                  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                  PRIMARY KEY (`id`) USING BTREE,
                                  UNIQUE KEY `idx_code` (`code`) USING BTREE COMMENT '兑换码索引'
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='兑换码表';