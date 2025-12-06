CREATE TABLE `us_base_user`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_no` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户编号[只做展示用途]',
    `mask_phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '加密手机号',
    `phone` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '脱敏手机号',
    `product` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '所属产品 0.辞源',
    `source` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '首次来源 100.提审标记 0.APP 1.豆伴匠 2.风的颜色 1001.听力熊 1002.微软',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 0.正常 1.冻结 2.注销',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uni_mask_phone_product` (`mask_phone`, `product`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT ='用户主表';

CREATE TABLE `us_user_auth`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `base_user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `identity_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '授权类型 1.手机号 2.微信授权 3.一键登录 4.设备标识',
    `identifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录账号 例：手机号、邮箱、用户名或第三方应用的唯一标识',
    `certificate` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码凭证 例：站内的密码，站外的token',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 0.正常 1.冻结 2.注销',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uni_identifier` (`identifier`, `identity_type`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT ='用户授权表';

CREATE TABLE `us_user`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `base_user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '主用户ID',
    `is_default` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否为默认账户 1.默认 2.非默认',
    `role_type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '角色类型 1.孩子[默认] 2.家长',
    `nickname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户昵称',
    `real_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '真实姓名',
    `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '用户头像',
    `gender` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '性别 0.未知 1.男 2.女',
    `grade` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '年级 10.学龄前 11.一年级 12.二年级 13.三年级 14.四年级 15.五年级 16.六年级 20.中小衔接 21.初一 22.初二 23.初三 31.高一 32.高二 33.高三 41.大学语文 42.中文专业 43.文学专家',
    `birthday` date DEFAULT NULL COMMENT '生日[年月日]',
    `province_id` int unsigned NOT NULL DEFAULT '0' COMMENT '省ID',
    `province_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '省名称',
    `city_id` int unsigned NOT NULL DEFAULT '0' COMMENT '市ID',
    `city_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '市名称',
    `area_id` int unsigned NOT NULL DEFAULT '0' COMMENT '区ID',
    `area_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '区名称',
    `is_info_finished` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否完善信息 1.未完善 2.已完善',
    `device_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '最近一次登录设备号',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '用户状态 1.正常 2.删除',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT ='用户信息表';

CREATE TABLE `us_user_login_log`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `base_user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '主用户ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '子用户ID',
    `product` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '所属产品 0.辞源',
    `identity_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '授权类型 1.手机号 2.微信 3.一键登录 4.设备标识',
    `command` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '操作类型 1.注册成功 2.注册失败 3.登录成功 4.登录失败 5.登出成功 6.登出失败',
    `client` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '客户端',
    `client_version` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '客户端版本号',
    `device_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '登录时设备号',
    `last_ip` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '登录ip',
    `os` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '手机系统',
    `os_version` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '手机系统版本',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_uid` (`base_user_id`),
    KEY `idx_uid_type_cmd` (`base_user_id`, `identity_type`, `command`),
    KEY `idx_ct` (`created_at`),
    KEY `idx_uid_type_prd_cmd` (`base_user_id`, `identity_type`, `product`, `command`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT ='用户登录日志表';

CREATE TABLE `us_white_list`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `phone` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '白名单手机号',
    `start_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '有效期起始时间',
    `end_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '有效期起始时间',
    `product` tinyint NOT NULL COMMENT '产品线',
    `source` tinyint NOT NULL COMMENT '来源',
    `status` tinyint NOT NULL COMMENT '状态 1.正常 2.失效',
    `remark` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uni_product_source_phone` (`product`,`source`,`phone`) USING BTREE,
    KEY `idx_start_end_time` (`start_time`,`end_time`),
    KEY `idx_status` (`status`),
    KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户白名单表';