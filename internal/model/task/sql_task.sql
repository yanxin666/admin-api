CREATE TABLE `bk_sync_task`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `operate_id` int NOT NULL DEFAULT '0' COMMENT '操作人id',
    `type` tinyint NOT NULL DEFAULT '1' COMMENT '1:导入',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '0:未开始 1:进行中 2:已完成 3:失败',
    `param` json DEFAULT NULL COMMENT '排序',
    `file_id` varchar(64) NOT NULL DEFAULT '' COMMENT '文件id',
    `file_name` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名称',
    `file_sheet` int unsigned NOT NULL DEFAULT '1' COMMENT '文件表对象',
    `file_sheet_name` varchar(255) NOT NULL COMMENT '文件表对象名称',
    `start_time` varchar(64) NOT NULL DEFAULT '' COMMENT '开始时间',
    `end_time` varchar(64) NOT NULL DEFAULT '' COMMENT '结束时间',
    `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '错误信息',
    `total` int NOT NULL DEFAULT '0' COMMENT '总条数',
    `filter_num` int NOT NULL DEFAULT '0' COMMENT '过滤条数',
    `pre_num` int NOT NULL DEFAULT '0' COMMENT '预处理条数',
    `success_num` int NOT NULL DEFAULT '0' COMMENT '成功条数',
    `fail_num` int NOT NULL DEFAULT '0' COMMENT '失败条数',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='管理后台异步任务';

CREATE TABLE `bk_sync_task_log`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `task_id` bigint NOT NULL COMMENT '任务id',
    `index` bigint NOT NULL COMMENT '当前任务的数据下标',
    `data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '原始数据',
    `md5` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'data字段进行md5,去重使用',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 0.未开始 1.进行中 2.已完成 3.失败',
    `errors_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '错误原因',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_md5` (`md5`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='导入数据快照表';

CREATE TABLE `bk_schedule_log`
(
    `id` bigint NOT NULL AUTO_INCREMENT,
    `data` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '原始数据',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 0:未推送 1.已推送',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='题库规划导入快照表';