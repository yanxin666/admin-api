CREATE TABLE `ns_learn_records`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课堂ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `completed_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '第一次完成时间',
    `completed_count` int unsigned NOT NULL DEFAULT '0' COMMENT '累计完成次数',
    `completed_sign` char(15) COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '完成标记',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_course_user_sign` (`course_id`,`user_id`,`completed_sign`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ai互动课堂学习记录';

CREATE TABLE `ns_homework_records`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `live_course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课堂ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '作业批改状态 1.批改中 2.批改完成 3.批改失败',
    `correct_comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '批改内容',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ai互动课堂作业记录表';

CREATE TABLE `ns_homework_detail`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `homework_records_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '作业记录id',
    `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'url',
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '识别内容',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态1.待识别 2.完成 3.失败',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ai互动课堂作业详情表';

-- 创建ai直播课用户预约表
CREATE TABLE `ns_live_appointment`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
    `live_series_id` bigint NOT NULL DEFAULT '0' COMMENT '系列id',
    `live_course_id` bigint NOT NULL DEFAULT '0' COMMENT '课堂ID',
    `appointment_time` datetime NOT NULL COMMENT '预约时间',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 1.未学习 2.学习中 3.已完成',
    `home_work_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课后作业状态 1.未完成 2.已完成',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ai直播课用户预约表';

-- 创建ai直播课程系列年级表
CREATE TABLE `ns_live_series_grade`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `live_series_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '系列id',
    `grade` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '适用年级',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ai直播课程系列年级表';

-- 创建ai直播课程系列表
CREATE TABLE `ns_live_series`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '系列名',
    `cover_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程系列封面',
    `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程系列介绍',
    `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '系列类型 1.试听系列 2.付费系列',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 1.正常 2.禁用',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ai直播课程系列表';

-- 创建ai直播课程系列课程关系表
CREATE TABLE `ns_live_series_course_relation`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `live_series_id` bigint NOT NULL COMMENT '系列 id',
    `live_course_id` bigint NOT NULL COMMENT '课程 id',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 创建ai直播老师表
CREATE TABLE `ns_live_teacher`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `teacher_code` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '数据源提供的老师ID',
    `teacher_name` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '老师名称',
    `teacher_nickname` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '老师昵称',
    `hand_audio` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课中举手发声音频',
    `tone_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '音色类型 1.minimax 2.火山',
    `teacher_tone` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '老师音色',
    `tone_model` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '音色模型',
    `teacher_introduce` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '老师简介',
    `teacher_gif` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '老师讲课动图(兜底图)',
    `teacher_png` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '老师讲课静图(兜底图)',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ai直播课程系列表';

-- 创建直播课程表
CREATE TABLE `ns_live_course`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `live_no` char(17) NOT NULL DEFAULT '' COMMENT '课程编号',
    `version` int NOT NULL DEFAULT '1' COMMENT '版本号',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程名称',
    `teacher_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '老师id',
    `label` json NOT NULL COMMENT '课程标签',
    `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程封面图',
    `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '课程介绍',
    `keyword` json NOT NULL COMMENT '关键知识点',
    `open_time` datetime NOT NULL COMMENT '开放时间',
    `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '类型 1.试听课 2.付费课',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '备注',
    `sequence` double(5, 2) unsigned NOT NULL DEFAULT '0.00' COMMENT '排序',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播课程表';

-- 创建直播课程详情表
CREATE TABLE `ns_live_course_detail`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `live_course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课堂ID',
    `teacher_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '老师ID',
    `teacher_gif` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '老师讲课动图',
    `teacher_png` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '老师讲课静图',
    `teacher_video` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '老师讲课视频URL',
    `duration` int NOT NULL DEFAULT '0' COMMENT '视频总时长，单位毫秒',
    `topic` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '课堂开始时的寒暄主题',
    `ppt_video` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'PPT视频URL',
    `full_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '全屏背景URL',
    `poster_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '视频开场背景URL',
    `note` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '课堂笔记',
    `homework` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '课后作业',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播课程详情表';

-- 创建直播课堂时间线表
CREATE TABLE `ns_live_course_timeline`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `live_course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课堂ID',
    `navigate` json NOT NULL COMMENT '导航进度条',
    `draft` json NOT NULL COMMENT '字幕',
    `draft_rag` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字幕纯文字',
    `simple_draft` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '提炼出的精简字幕',
    `block_time` json NOT NULL COMMENT '停顿气口',
    `front_event` json NOT NULL COMMENT '班主任欢迎环节',
    `small_talk_event` json NOT NULL COMMENT '寒暄环节',
    `small_talk_precast` json NOT NULL COMMENT '寒暄环节中的开场预制',
    `learn_event` json NOT NULL COMMENT '课中环节',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播课堂时间线表';

-- 创建直播课堂事件表
CREATE TABLE `ns_live_course_event`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `live_course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课堂ID',
    `start_id` float unsigned NOT NULL DEFAULT '0' COMMENT '数据源唯一ID',
    `event_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '事件类型 1.填空 2.选择',
    `appearance` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '题型风格 1.我说上句你来接 2.排序题 3.课前简答题 4.课后简答题 5.文本单选题 6.图片选择题 7.火眼金睛',
    `prompt_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'LLM Code',
    `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '标题',
    `task` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '题目描述',
    `talk` json NOT NULL COMMENT '交互过程中的会话',
    `without_answer_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '无答案图片URL',
    `within_answer_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '有答案图片URL',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播课堂事件表';

-- 创建直播课堂题目表
CREATE TABLE `ns_live_course_question`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `live_course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课堂ID',
    `live_course_event_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '互动事件ID',
    `question_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '题目类型 1.填空题 2.选择题',
    `question` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '问题',
    `answer_ask` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '答题要求',
    `tips` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '提示',
    `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '图片',
    `answer` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '题目答案',
    `analysis` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '解析',
    `sequence` double(5, 2) unsigned NOT NULL DEFAULT '0.00' COMMENT '排序',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播课堂题目表';

-- 创建直播课堂题目选项表
CREATE TABLE `ns_live_course_question_option`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `live_course_question_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '题目ID',
    `tag` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '选项名称',
    `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '选项内容',
    `tips` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '提示',
    `image` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '图片',
    `sequence` double(5,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '排序',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播课堂题目选项表';

-- 创建直播课堂学习记录表
CREATE TABLE `ns_live_record`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
    `appointment_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '预约id',
    `live_course_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课程id',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态1.进行中 2.结束',
    `completed_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '完成时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播课堂学习记录表';

-- 创建直播课堂参与用户表
CREATE TABLE `ns_live_user`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `live_record_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '直播记录ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `user_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '用户类型：1.真人 2.虚拟人',
    `role_id` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '预制数据源的角色ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='直播课堂参与用户表';

-- 创建直播课堂虚拟用户表
CREATE TABLE `ns_live_virtual_user`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '姓名',
    `gif` varchar(500) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '动态头像URL',
    `role_id` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '角色性格类型',
    `role_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '角色类型：1.学生 2.助教老师',
    `description` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '人物描述、性格、特长',
    `sex` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '性别：1.男 2.女',
    `status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '状态：1.有效 2.无效',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='直播课堂虚拟用户表';

-- 创建直播课堂学习记录表
CREATE TABLE `ns_live_answer`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `live_record_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '直播记录ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '直播用户ID',
    `code` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '业务类型',
    `game_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '互动ID',
    `content` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '内容',
    `role` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '角色：1 .模型 2.用户 3.虚拟用户',
    `index` int unsigned NOT NULL DEFAULT '0' COMMENT '视频播放时间戳',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='直播课堂学习记录表';