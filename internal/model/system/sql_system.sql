CREATE TABLE `sys_sms_record`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '加密手机号',
    `code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '调用sdk后的状态码',
    `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '状态码说明',
    `channel` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '渠道 1.腾讯云 2.阿里云',
    `template` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '短信模版id',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 1.成功 2.失败',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_temp_phone` (`template`, `phone`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='短信发送记录表';

CREATE TABLE `sys_dict_white` (
                                  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                                  `key` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典Key',
                                  `remark` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
                                  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
                                  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                                  PRIMARY KEY (`id`),
                                  KEY `idx_key` (`key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统字典表前端白名单';