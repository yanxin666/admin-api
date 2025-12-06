CREATE TABLE `ac_course`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `course_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程编号',
    `use_cases` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '应用场景：1.练习 2.直播互动',
    `course_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课程类型 1.超练作文慢练 2.超练作文快练 3.超练阅读',
    `version` int NOT NULL DEFAULT '1' COMMENT '版本号',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程名称',
    `lesson_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课节名称',
    `subject` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '学科 1.语文',
    `unit` varchar(64) NOT NULL DEFAULT '' COMMENT '单元',
    `teacher_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '老师名称',
    `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程封面图',
    `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程介绍',
    `open_time` datetime NOT NULL COMMENT '开放时间',
    `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '类型 1.试听课 2.付费课',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `level` json DEFAULT NULL COMMENT '难度',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '备注',
    `extra` json DEFAULT NULL COMMENT '拓展字段，用于存json',
    `sequence` double(5,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '排序',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AiClass课程表';

CREATE TABLE `ac_course_chapter`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `chapter_no` varchar(64) NOT NULL DEFAULT '' COMMENT '章节编号',
    `type` tinyint NOT NULL DEFAULT '1' COMMENT '1 任务 2 文章',
    `course_id` bigint NOT NULL DEFAULT '0' COMMENT '课程ID',
    `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题，后续不再维护',
    `title` varchar(64) NOT NULL COMMENT '精简版标题',
    `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程封面图',
    `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程介绍',
    `index` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '圈数',
    `guide_video` varchar(255) NOT NULL DEFAULT '' COMMENT '引导视频',
    `sequence` double NOT NULL DEFAULT '0' COMMENT '系列序号',
    `teacher_id` bigint NOT NULL DEFAULT '0' COMMENT '老师id',
    `status` tinyint unsigned NOT NULL DEFAULT '3' COMMENT '状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `can_learn` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否能学 1.可以 2.不可以',
    `is_new` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否上新 1.否 2.是',
    `extra` json DEFAULT NULL COMMENT '拓展字段，用于存json',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_chapter` (`chapter_no`),
    KEY `idx_course` (`course_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='课程章节表';

CREATE TABLE `ac_chapter_task`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `task_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务编号',
    `title` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
    `subtitle` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '子标题',
    `description` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述',
    `course_id` bigint NOT NULL DEFAULT '0' COMMENT '课堂id',
    `course_chapter_id` bigint NOT NULL DEFAULT '0' COMMENT '章节id',
    `progress_prefix` varchar(32) NOT NULL DEFAULT '' COMMENT '进度前缀、时空稳定值、时空通道已建立等',
    `extra` json NOT NULL COMMENT '拓展字段，用于存json',
    `type` int NOT NULL COMMENT '类型：1.文章段落 2.身份卡 3.文章总结 4.时空长廊 200.主线任务 201.副本任务',
    `sequence` double(5,2) NOT NULL DEFAULT '0.00' COMMENT '排序',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否删除 1.未删除 2.已删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_task` (`task_no`),
    KEY `idx_course_chapter` (`course_chapter_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='课程章节任务表';

CREATE TABLE `ac_sub_task`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mode` tinyint NOT NULL DEFAULT '1' COMMENT '子任务模式 1.普通模式 2. 困难模式',
    `sub_task_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '子任务编号',
    `chapter_task_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '章节任务ID',
    `title` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '标题',
    `article` text COMMENT '原文资料',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '描述',
    `image_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '图片地址',
    `keywords` json NOT NULL COMMENT '关键词列表(json格式)',
    `extra` json NOT NULL COMMENT '拓展字段(json格式)',
    `type` int NOT NULL COMMENT '类型 1.普通任务 2.身份卡 3.聊天',
    `sequence` double(5,2) NOT NULL DEFAULT '0.00' COMMENT '排序',
    `deleted` tinyint NOT NULL DEFAULT '1' COMMENT '是否删除 1.否 2.是',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_sub_task` (`sub_task_no`),
    KEY `idx_task` (`chapter_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='子任务表';

CREATE TABLE `ac_medal`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name` varchar(32) NOT NULL DEFAULT '' COMMENT '名称',
    `description` varchar(256) NOT NULL DEFAULT '' COMMENT '描述',
    `type` tinyint NOT NULL COMMENT '类型 1.金钥匙 2.银钥匙 3.铜钥匙',
    `prompt_knowledge` varchar(64) NOT NULL DEFAULT '' COMMENT 'p工程类关联知识点',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='任务奖牌';

CREATE TABLE `ac_sub_task_medal`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `sub_task_id` bigint NOT NULL COMMENT '子任务ID',
    `medal_id` bigint NOT NULL COMMENT '任务奖牌ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_sub_task_id` (`sub_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='子任务奖牌关联表';