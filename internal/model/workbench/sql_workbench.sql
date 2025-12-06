CREATE TABLE `wk_dept`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
    `parent_id` int unsigned NOT NULL COMMENT '父级id',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '部门简称',
    `full_name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '部门全称',
    `unique_key` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一值',
    `type` tinyint unsigned NOT NULL COMMENT '1.公司 2.子公司 3.部门',
    `status` tinyint unsigned NOT NULL COMMENT '0.禁用 1.开启',
    `order_num` int unsigned NOT NULL COMMENT '排序值',
    `remark` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_key` (`unique_key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='部门';

CREATE TABLE `wk_dictionary`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
    `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '0.配置集 !0.父级id',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
    `type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1.文本 2.数字 3.数组 4.单选 5.多选 6.下拉 7.日期 8.时间 9.单图 10.多图 11.单文件 12.多文件',
    `unique_key` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一值',
    `value` varchar(2048) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '配置值',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '0.禁用 1.开启',
    `order_num` int unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
    `remark` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_key` (`unique_key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='系统参数';

CREATE TABLE `wk_job`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '岗位名称',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '0.禁用 1.开启 ',
    `order_num` int unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '开启时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='工作岗位';

CREATE TABLE `wk_log`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
    `user_id` int unsigned NOT NULL COMMENT '操作账号',
    `ip` varchar(100) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'ip',
    `uri` varchar(200) COLLATE utf8mb4_general_ci NOT NULL COMMENT '请求路径',
    `type` tinyint unsigned NOT NULL COMMENT '1.登录日志 2.操作日志',
    `request` varchar(2048) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '请求数据',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '0.失败 1.成功',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='系统日志';

CREATE TABLE `wk_perm_menu`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
    `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
    `router` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '路由',
    `perms` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '权限',
    `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '0.目录 1.菜单 2.权限',
    `icon` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图标',
    `order_num` int unsigned DEFAULT '0' COMMENT '排序值',
    `view_path` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页面路径',
    `is_show` tinyint unsigned DEFAULT '1' COMMENT '0.隐藏 1.显示',
    `active_router` varchar(1024) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '当前激活的菜单',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='权限&菜单';

CREATE TABLE `wk_profession`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '职称',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '0.禁用 1.开启',
    `order_num` int unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='职称';

CREATE TABLE `wk_role`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
    `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
    `unique_key` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一标识',
    `remark` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
    `perm_menu_ids` json NOT NULL COMMENT '权限集',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '0.禁用 1.开启',
    `order_num` int unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_key` (`unique_key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='角色';

CREATE TABLE `wk_user`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
    `account` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账号',
    `password` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
    `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '姓名',
    `mobile` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
    `open_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '飞书open_id',
    `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '飞书user_id',
    `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮件',
    `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
    `gender` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '0.保密 1.女 2.男',
    `avatar` varchar(400) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
    `profession_id` int unsigned NOT NULL COMMENT '职称',
    `job_id` int unsigned NOT NULL COMMENT '岗位',
    `dept_id` int unsigned NOT NULL COMMENT '部门',
    `role_ids` json NOT NULL COMMENT '角色集',
    `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '0.禁用 1.开启',
    `order_num` int unsigned NOT NULL DEFAULT '0' COMMENT '排序值',
    `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `account` (`account`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户';