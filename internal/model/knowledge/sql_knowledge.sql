CREATE TABLE `kn_question`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `question_no` varchar(32) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '题目编号',
    `biz_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '业务类型 1.深度探测 2.爬天梯',
    `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '题目类型，1：选择；2：简答；',
    `review_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `level` int unsigned NOT NULL DEFAULT '0' COMMENT '等级',
    `high_ladder_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '天梯ID',
    `primary_point_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '知识点ID',
    `material` text COLLATE utf8mb4_general_ci COMMENT '题目素材，比如选取的阅读素材片段',
    `ask` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '问题',
    `answer` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '答案，针对简答题，用于答案判别',
    `analysis` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '题目解析，用于答案判别',
    `version` int unsigned NOT NULL DEFAULT '0' COMMENT '版本',
    `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建者，程序导入为0，其他为操作者ID',
    `updated_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新者，程序导入为0，其他为操作者ID',
    `effective_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '有效期起始时间',
    `expiry_time` datetime NOT NULL DEFAULT '9999-12-31 23:59:59' COMMENT '有效期截止时间',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='题目表';

CREATE TABLE `kn_question_point`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `question_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '问题ID',
    `point_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '知识点ID',
    `is_primary` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否为主知识点 1.是 2.否',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='问题知识点表';

CREATE TABLE `kn_material`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '素材名称',
    `author` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '作者',
    `source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '素材来源',
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '素材内容，比如选取的阅读素材片段',
    `background` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '写作背景',
    `author_intro` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '作者介绍',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='题目素材表';

CREATE TABLE `kn_example`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '例题内容，JSON格式 旧版',
    `explain` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '例题内容，JSON格式 新版',
    `plain_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'JSON中的纯文本内容',
    `plain_text_article` text COLLATE utf8mb4_general_ci NOT NULL COMMENT 'JSON中的纯文本文章',
    `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
    `sub_title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '副标题',
    `lesson_notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '课堂笔记',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='例题表';

CREATE TABLE `kn_lesson_question`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `question_no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '题目编号',
    `node_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '节点类型 1.小灶课 2.小语文 3.大语文',
    `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '试题类型 1：单选，2：多选，3：填空，4：判断，5：简答，6：阅读题 7：作文',
    `usage_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '使用类型，1：例题；2：练习题；3：候补题',
    `grade_phase` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '年级学段，1：小学；2：初中；3：高中',
    `review_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `level` int unsigned NOT NULL DEFAULT '0' COMMENT '等级',
    `material_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '题目素材Id',
    `example_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '例题id',
    `ask` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '问题',
    `answer` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '答案，针对简答题，用于答案判别',
    `analysis` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '题目解析，用于答案判别',
    `start_tts` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '开场白TTS',
    `duration` bigint unsigned NOT NULL DEFAULT '0' COMMENT '预计用时',
    `source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源',
    `version` int unsigned NOT NULL DEFAULT '0' COMMENT '版本',
    `created_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建者，程序导入为0，其他为操作者ID',
    `updated_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '更新者，程序导入为0，其他为操作者ID',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='学练题目表';

CREATE TABLE `kn_lesson_question_option`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `lesson_question_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '题目ID',
    `sequence` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '选项展示顺序[正序]',
    `option_label` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '选项标签[例.A、B、C、D等]',
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '选项内容,描述性内容或词汇ID',
    `is_answer` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否为正确答案 1.是 2.否',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='学练题目选项表';

CREATE TABLE `kn_lesson`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `lesson_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '课节来源编号',
    `lesson_group_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '课程组编号',
    `node_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '导入数据时使用，节点类型 1.小灶课 2.小语文 3.大语文',
    `parent_id` bigint unsigned NOT NULL COMMENT '导入数据时使用，只有大语文才会用到，用作关联父子级',
    `level` int unsigned NOT NULL DEFAULT '0' COMMENT '难度等级 只有小灶课才会使用，1-4个等级',
    `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '课节名称',
    `lesson_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课程类型 1.主线课 2.月考  3.小灶课',
    `review_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='课节表';

CREATE TABLE `kn_lesson_resource`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `lesson_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课节ID',
    `question_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '题目ID',
    `sequence` double(5, 2) unsigned NOT NULL DEFAULT '0.00' COMMENT '排序，正序，从1开始，默认跨度1',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='课节资源表';

CREATE TABLE `kn_lesson_point`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `lesson_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '课节ID',
    `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '知识点名称',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='课节知识点表';

CREATE TABLE `kn_lesson_question_tts`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `question_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '题目ID',
    `content` varchar(500) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'TTS文本',
    `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'TTS音频地址',
    `apply_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '用途类型：1.开场白',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_question_id` (`question_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='题目TTS表';