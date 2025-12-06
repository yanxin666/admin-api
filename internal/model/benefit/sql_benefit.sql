CREATE TABLE `bt_benefits_resource`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '权益名称',
    `type` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '权益类型 1.快解题 2.关卡 3.风的颜色 4.学习规划',
    `is_limited_free` tinyint unsigned NOT NULL DEFAULT 2 COMMENT '是否限免，1.不是 2.是',
    `description` TEXT DEFAULT NULL COMMENT '权益描述',
    `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '权益状态 1.启用 2.禁用',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='权益资源表';


-- 关卡权益详情表
CREATE TABLE `bt_benefits_detail`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `benefits_resource_id` bigint unsigned NOT NULL COMMENT '权益资源表ID',
    `resource_type` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '资源类型 11.快解知识点 12.快解题单元 13.快解全部 21.阅读理解关卡 22.阅读理解全部 31.风的颜色主题 32.风的颜色全部',
    `resource_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '资源ID,仅对 resource_type = 1 有效',
    `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '权益状态 1.启用 2.禁用',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='关卡权益详情表';

-- 用户权益变更记录表
CREATE TABLE `bt_benefits_user_record`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '权益变更类型 1: 发放 2: 回收',
    `order_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '订单ID，已废弃',
    `order_no` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '订单号',
    `benefits_resource_id` bigint NOT NULL DEFAULT '0' COMMENT '权益资源表ID',
    `user_benefits_id` bigint unsigned NOT NULL COMMENT '用户权益表',
    `source` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '来源 1：内部员工；2：渠道三部 3：中台；4：蓝V企业店； 5：产品运营',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE COMMENT '用户权益变更记录表userId 索引'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户权益变更记录表';

-- 权益通知记录表
CREATE TABLE `bt_benefits_notify_log`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `order_no` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '订单编号',
    `source` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '来源 1.豆伴匠',
    `data` json DEFAULT NULL COMMENT '通知原始数据',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='权益通知记录表';

-- 权益组表
CREATE TABLE `bt_benefits_group`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '权益名称',
    `code` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '编码',
    `user_type` tinyint unsigned NOT NULL DEFAULT 2 COMMENT '用户类型 1.单个子账号 2.主账号',
    `type` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '1.单品权益 2.组合权益',
    `intro` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '权益介绍',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='权益组表';

-- 权益组详情表
CREATE TABLE `bt_benefits_group_detail`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `benefits_group_id` bigint unsigned NOT NULL COMMENT '权益组表ID',
    `benefits_resource_id` bigint unsigned NOT NULL COMMENT '权益资源表ID',
    `kind` tinyint unsigned NOT NULL COMMENT '权益类别，1.永久 2.时效的 3.时间段',
    `from_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    `days` int DEFAULT '-1' COMMENT '生效天数，-1 代表永久',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='权益组详情表';

-- 渠道订单表
CREATE TABLE `bt_channel_order`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `order_no` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '外部订单号',
    `benefits_group_id` bigint unsigned NOT NULL COMMENT '权益类别表ID',
    `mask_phone` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '脱敏手机号',
    `user_type` tinyint unsigned NOT NULL DEFAULT 2 COMMENT '用户类型 1.单个子账号 2.主账号',
    `order_type` tinyint unsigned NOT NULL DEFAULT 2 COMMENT '订单类型 1.发放 2.回收',
    `source` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '来源 1.内部员工 2.渠道三部 3.中台 4.蓝V企业店 5.产品运营',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='渠道订单表';

-- 渠道子订单表
CREATE TABLE `bt_channel_sub_order`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `channel_order_id` bigint unsigned NOT NULL COMMENT '渠道订单表ID',
    `base_user_id` bigint unsigned NOT NULL COMMENT '主用户ID',
    `user_id` bigint unsigned NOT NULL COMMENT '子用户ID',
    `benefits_resource_id` bigint unsigned NOT NULL COMMENT '权益资源表ID',
    `benefits_group_detail_id` bigint unsigned NOT NULL COMMENT '权益组详情表ID',
    `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '订单状态.1.发放中 2.已发放 3.回收中 4.已回收',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='渠道子订单表';