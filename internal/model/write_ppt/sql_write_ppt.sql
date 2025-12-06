CREATE TABLE `aw_grade_course_series`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `grade` tinyint NOT NULL COMMENT '年级',
    `course_series` char(16) COLLATE utf8mb4_general_ci NOT NULL COMMENT '课程系列',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='私教课堂年级大纲对应表';

CREATE TABLE `aw_course_outline`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '单元名',
    `unit` tinyint NOT NULL DEFAULT '0' COMMENT '单元',
    `category` tinyint NOT NULL DEFAULT '0' COMMENT '分类 1.细节 2.表达 3.布局',
    `topic` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '写作要求',
    `series` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '课程系列',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '数据是否已处理 1.未处理 2.已处理',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='私教课堂写作大纲表';

CREATE TABLE `aw_lesson`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `course_outline_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '私教写作大纲',
    `lesson_no` bigint unsigned NOT NULL DEFAULT '0' COMMENT '原数据来源ID',
    `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '课程名',
    `sub_title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '子标题',
    `lesson_number` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课节序号',
    `lesson_type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '课节类型 1.技巧 2.赏析',
    `lesson_content` json NOT NULL COMMENT '讲稿',
    `lesson_content_text` json NOT NULL COMMENT '讲稿中的纯文本内容',
    `notes` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '笔记',
    `ppt_files` json NOT NULL COMMENT 'ppt 文件 url 列表',
    `mind_slices` json NOT NULL COMMENT '脑图步骤 json',
    `mind_full` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '完整脑图 json',
    `audio_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '讲稿音频地址',
    `review_status` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_course_id` (`course_outline_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='私教课堂写作课节表';
