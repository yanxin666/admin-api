CREATE TABLE `us_user_auth`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `base_user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `identity_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '授权类型 1.手机号 2.微信授权 3.一键登录 4.设备标识',
    `identifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录账号 例：手机号、邮箱、用户名或第三方应用的唯一标识',
    `certificate` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码凭证 例：站内的密码，站外的token',
    `product` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '所属产品 0.辞源',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '登录方式状态 0.正常 1.冻结 2.注销',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uni_type_identifier_product` (`identity_type`, `identifier`, `product`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户授权表';

CREATE TABLE `us_base_user`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_no` varchar(16) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户编号[只做展示用途]',
    `mask_phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '加密手机号',
    `phone` char(11) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '脱敏手机号',
    `product` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '所属产品 0.辞源',
    `source` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '首次来源 100.提审标记 0.APP 1.豆伴匠 2.风的颜色 1001.听力熊 1002.微软',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 0.正常 1.冻结 2.注销',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uni_mask_phone_product` (`mask_phone`, `product`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户信息表';

CREATE TABLE `us_user`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `base_user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `is_default` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否为默认账户 1.默认 2.非默认',
    `role_type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '角色类型 1.孩子[默认] 2.家长',
    `nickname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户昵称',
    `real_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '真实姓名',
    `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户头像',
    `gender` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '性别 0.未知 1.男 2.女',
    `grade` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '年级 10.学龄前 11.一年级 12.二年级 13.三年级 14.四年级 15.五年级 16.六年级 20.中小衔接 21.初一 22.初二 23.初三 31.高一 32.高二 33.高三 41.大学语文 42.中文专业 43.文学专家',
    `birthday` date DEFAULT NULL COMMENT '生日[年月日]',
    `province_id` int unsigned NOT NULL DEFAULT '0' COMMENT '省ID',
    `province_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '省名称',
    `city_id` int unsigned NOT NULL DEFAULT '0' COMMENT '市ID',
    `city_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '市名称',
    `area_id` int unsigned NOT NULL DEFAULT '0' COMMENT '区ID',
    `area_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '区名称',
    `is_info_finished` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否完善信息 1.未完善 2.已完善',
    `device_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '最近一次登录设备号',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '用户状态 1.正常 2.删除',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户信息表';

CREATE TABLE `us_user_opinion`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `urls` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'urls路径，“|” 分割，支持多个',
    `comment` varchar(500) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '意见',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '0.新意见 1.已回复',
    `reply` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '回复信息',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='意见反馈表';