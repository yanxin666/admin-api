CREATE TABLE `ai_course`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `course_no` bigint NOT NULL COMMENT '来源编号，导数据时使用',
    `name` varchar(30) NOT NULL COMMENT '课程名称',
    `description` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '课程简介',
    `teacher_name` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '主讲老师',
    `pre_step` json NOT NULL COMMENT '预制课程步骤',
    `homework` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '作业要求',
    `evaluate` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '评分标准',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态(0.待上架 1.正常 2.已下架)',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AI课程信息表';

CREATE TABLE `ai_live_role`
(
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `role_no` bigint unsigned NOT NULL COMMENT '数据源角色编号，导数据时使用',
    `role_type` varchar(30) NOT NULL COMMENT '角色类型 "teacher":教师 | "assistant":助教 | "student":学生',
    `name` varchar(30) NOT NULL COMMENT '名称',
    `feature` varchar(1000) NOT NULL COMMENT '特点',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='直播虚拟角色';

CREATE TABLE `ai_live_beta_records`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `phone` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户手机号',
    `user_name` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
    `company` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '公司名',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '审核状态 1.审核中 2.审核通过 3.审核不通过',
    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注信息',
    `operate_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '操作人id',
    `operate_at` int unsigned NOT NULL DEFAULT '0' COMMENT '操作时间戳',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;