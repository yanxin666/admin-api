CREATE TABLE `supv_schedule`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `teacher_id` bigint NOT NULL COMMENT '讲师ID',
    `teacher_name` varchar(50) NOT NULL COMMENT '讲师名称',
    `course_id` bigint NOT NULL COMMENT '课程ID',
    `course_type` tinyint NOT NULL DEFAULT '0' COMMENT '类型 1-超级作文 2-超练阅读',
    `name` varchar(50) NOT NULL COMMENT '督学课名称',
    `stream_id` bigint NOT NULL COMMENT '直播ID',
    `effective_start` datetime NOT NULL COMMENT '活动生效时间',
    `effective_end` datetime NOT NULL COMMENT '活动结束时间',
    `appointment_start` datetime NOT NULL COMMENT '预约开始时间',
    `appointment_end` datetime NOT NULL COMMENT '预约结束时间',
    `start_time` datetime NOT NULL COMMENT '开课时间',
    `max_stock` int unsigned NOT NULL DEFAULT '0' COMMENT '预约人数上限',
    `surplus_stock` int unsigned NOT NULL DEFAULT '0' COMMENT '剩余预约人数',
    `schedule_type` tinyint NOT NULL DEFAULT '1' COMMENT '排课类型 1直播 2录播',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '排课状态 1待开放 2开放 3进行中 4结束',
    `bg_url` json NOT NULL COMMENT '伴学、督学课程背景图片',
    `content` json NOT NULL COMMENT '伴学、督学课程文案',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_teacher_id` (`teacher_id`),
    KEY `idx_course_id` (`course_id`),
    KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='督学排课表';

CREATE TABLE `lv_stream`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 Id',
    `teacher_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '讲师ID',
    `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '直播标题',
    `room_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间号',
    `playback_task_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '直播云端录制任务ID',
    `playback_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '直播回放地址',
    `is_mute` tinyint NOT NULL DEFAULT '2' COMMENT '是否禁言（1是，2否）',
    `start_time` datetime NOT NULL COMMENT '开播时间',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '直播状态：1.待开始 2.进行中 3.结束',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='直播表';



CREATE TABLE `ac_teacher`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `teacher_code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'code',
    `teacher_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名字',
    `teacher_nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
    `teacher_tone` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '音色',
    `tone_model` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '音色模型',
    `teacher_introduce` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '介绍',
    `teacher_gif` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '动图',
    `teacher_png` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '静图',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播老师表';

CREATE TABLE `supv_interaction`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '互动名字',
    `description` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述',
    `teaching_data` json DEFAULT NULL COMMENT '教师ppt',
    `data` json DEFAULT NULL COMMENT '互动内容',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='督学互动表';

CREATE TABLE `supv_schedule_interaction_link`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `schedule_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '督学排课ID',
    `interaction_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '督学互动ID',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 1.有效 2.无效',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT ='督学互动关联表';

