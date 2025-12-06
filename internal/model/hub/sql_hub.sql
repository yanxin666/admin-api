CREATE TABLE `bk_lesson_snapshot`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `lesson_no` int unsigned NOT NULL DEFAULT '0' COMMENT '课节来源编号',
    `node_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '节点类型 1.小灶课 2.小语文 3.大语文',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程名称',
    `point_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '知识点名称',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `grade` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '年级 10.学龄前 11.一年级 12.二年级 13.三年级 14.四年级 15.五年级 16.六年级 20.中小衔接 21.初一 22.初二 23.初三 31.高一 32.高二 33.高三 41.大学语文 42.中文专业 43.文学专家',
    `version` int unsigned NOT NULL DEFAULT '1' COMMENT '版本号',
    `operate_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '操作状态 0.未推送 1.已推送 2.已下架',
    `operate_id` bigint NOT NULL DEFAULT '0' COMMENT '操作人ID',
    `data` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '原始数据',
    `app_version` varchar(64) NOT NULL DEFAULT '' COMMENT 'APP版本',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='课程快照表';

CREATE TABLE `bk_write_ppt_snapshot`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `lesson_no` int unsigned NOT NULL DEFAULT '0' COMMENT '来源编号',
    `unit` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '单元',
    `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '标题',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `lesson_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课节类型 1.技巧 2.赏析',
    `lesson_category` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课节类型 1.技巧 2.赏析',
    `version` int unsigned NOT NULL DEFAULT '1' COMMENT '版本号',
    `operate_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '操作状态 0.未推送 1.已推送 2.已下架',
    `operate_id` bigint NOT NULL DEFAULT '0' COMMENT '操作人ID',
    `data` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '原始数据',
    `app_version` varchar(64) NOT NULL DEFAULT '' COMMENT 'APP版本',
    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='写作PPT快照表';

CREATE TABLE `bk_live_snapshot`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `lesson_no` char(17) NOT NULL DEFAULT '' COMMENT '来源编号',
    `version_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '版本号编号',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程名称',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `version` int unsigned NOT NULL DEFAULT '1' COMMENT '版本号',
    `operate_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '操作状态 0.未推送 1.已推送 2.已下架',
    `operate_id` bigint NOT NULL DEFAULT '0' COMMENT '操作人ID',
    `data` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '原始数据',
    `app_version` varchar(64) NOT NULL DEFAULT '' COMMENT 'APP版本',
    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='AI直播快照表';

CREATE TABLE `bk_super_train_snapshot`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '编号',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程名称',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `version` int unsigned NOT NULL DEFAULT '1' COMMENT '版本号',
    `operate_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '操作状态 0.未推送 1.已推送 2.已下架',
    `operate_id` bigint NOT NULL DEFAULT '0' COMMENT '操作人ID',
    `data` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '原始数据',
    `app_version` varchar(64) NOT NULL DEFAULT '' COMMENT 'APP版本',
    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='超能训练快照表';