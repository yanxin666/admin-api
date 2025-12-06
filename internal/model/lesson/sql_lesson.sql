CREATE TABLE `ls_user_lesson_record`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户ID',
    `user_schedule_id` bigint NOT NULL DEFAULT '0' COMMENT '用户规划ID',
    `schedule_detail_id` bigint NOT NULL DEFAULT '0' COMMENT '规划详情ID',
    `biz_type` tinyint NOT NULL DEFAULT '1' COMMENT '业务类型  1 阅读理解',
    `lesson_type` tinyint NOT NULL DEFAULT '1' COMMENT '规划业务类型，1.课节 2.月考',
    `lesson_id` bigint NOT NULL DEFAULT '0' COMMENT '课程ID',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1.进行中 2.结束 3.中断 4.跳过',
    `total_num` int NOT NULL DEFAULT '0' COMMENT '总题数',
    `right_num` int NOT NULL DEFAULT '0' COMMENT '正确题目数',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE COMMENT '用户ID索引',
    KEY `idx_schedule_detail_id` (`schedule_detail_id`) USING BTREE COMMENT '规划详情索引'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户学习记录表';

CREATE TABLE `ls_user_lesson_detail`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户ID',
    `user_lesson_record_id` bigint NOT NULL DEFAULT '0' COMMENT '用户学习记录ID',
    `question_id` bigint NOT NULL DEFAULT '0' COMMENT '问题ID',
    `question_usage_type` tinyint NOT NULL COMMENT '使用类型，1.例题 2.练习题 ',
    `lesson_type` tinyint NOT NULL DEFAULT '1' COMMENT '规划业务类型，1.课程 2.月考',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '答题 态.1.进行 2.完成 3.跳过',
    `result` tinyint NOT NULL DEFAULT '0' COMMENT '是否正确 0.未知，1.正确 2.错误',
    `start_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    `finish_time` datetime DEFAULT NULL COMMENT '完成时间',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE COMMENT '用户ID索引',
    KEY `idx_user_lesson_record_id` (`user_lesson_record_id`) USING BTREE COMMENT '用户学习记录ID索引',
    KEY `idx_question_id` (`question_id`) USING BTREE COMMENT '问题ID索引'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户学习详情表';

CREATE TABLE `ls_user_lesson_answer`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户ID',
    `user_lesson_detail_id` bigint NOT NULL DEFAULT '0' COMMENT '用户学习详情ID',
    `lesson_type` tinyint NOT NULL DEFAULT '1' COMMENT '规划业务类型 1.课程 2.月考（单条记录更新）',
    `medium_type` tinyint NOT NULL COMMENT '回答类型',
    `content` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容 类型为选择，则为选项 文本则为文本 语音则为解析文本',
    `url` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '资源地址，针对语音形式的回答',
    `duration` int NOT NULL DEFAULT '0' COMMENT '时长，单位秒',
    `analysis` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '解析，非选择题时，GPT解析记录',
    `conversation_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '大模型上下文标识',
    `is_good` tinyint NOT NULL DEFAULT '0' COMMENT '用户点赞 0.默认值 1.点赞 2.点踩',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE COMMENT '用户ID索引',
    KEY `idx_user_lesson_detail_id` (`user_lesson_detail_id`) USING BTREE COMMENT '用户学习详情ID索引'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户作答表';

CREATE TABLE `ls_course`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '课程名称',
    `grade` int unsigned NOT NULL DEFAULT '0' COMMENT '所属年级',
    `term` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '学期 1.上学期 2.下学期 3.学年',
    `biz_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课程所属业务类型 1.阅读理解 ',
    `plan_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '规划类型 1.长期规划 2.周规划',
    `intro` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '课程介绍',
    `learn_target` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '学习目标',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_grade_term_biztype` (`grade`, `term`, `biz_type`) USING BTREE COMMENT '唯一索引'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='课程表';

CREATE TABLE `ls_course_outline`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课程ID',
    `unit` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '单元名称',
    `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '章节名称',
    `lesson_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课程类型 1.主线课 2.月考 3.小灶课',
    `grade` int unsigned NOT NULL DEFAULT '0' COMMENT '所属年级',
    `lesson_group_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '课程组编号',
    `expect_month` int unsigned NOT NULL DEFAULT '0' COMMENT '预期月份',
    `expect_date` date DEFAULT NULL COMMENT '预期上课日期',
    `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
    `sub_title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '副标题',
    `learn_target` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '学习目标',
    `sequence` double(5, 2) unsigned NOT NULL DEFAULT '0.00' COMMENT '排序，正序，月份相同排序',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `remark` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '导入时使用的课节名称',
    `mark_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '标记 id, 主线课小语文 id，小灶课为原理讲解 id',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='课程大纲表';


CREATE TABLE `ls_user_schedule`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `plan_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '规划类型 1.长期规划 2.周规划',
    `biz_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课程所属业务类型 1.阅读理解 ',
    `is_dynamic` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否动态规划 1.否 固定课程 2.是 动态生成',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '规划状态 1 正常',
    `end_date` date DEFAULT NULL COMMENT '规划结束日期',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户规划表';

CREATE TABLE `ls_user_schedule_detail`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `user_schedule_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户规划ID',
    `lesson_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课程类型 1.课节 2.月考',
    `course_outline_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课程大纲ID',
    `course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课程表ID',
    `lesson_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课节ID',
    `planned_date` date DEFAULT NULL COMMENT '规划安排日期',
    `grade` int unsigned NOT NULL DEFAULT '0' COMMENT '所属年级',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '学习状态 1 未学习 2.进行中 3.已学习',
    `sequence` double(5, 2) unsigned NOT NULL DEFAULT '0.00' COMMENT '排序，正序，月份相同排序',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户规划详情表';
