CREATE TABLE `c_conversation` (
      `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
      `user_id` bigint NOT NULL COMMENT '用户ID',
      `conversation_id` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'dify会话id',
      `prompt_id` bigint DEFAULT '0' COMMENT 'P工程ID',
      `title` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '新会话' COMMENT '会话标题',
      `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
      PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI会话表';

CREATE TABLE `c_conversation_record` (
      `id` bigint NOT NULL AUTO_INCREMENT,
      `conversation_id` bigint NOT NULL COMMENT 'AI会话ID',
      `content` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '消息内容',
      `source` tinyint NOT NULL COMMENT '消息来源 1 系统 2 GPT 3 用户',
      `type` tinyint NOT NULL COMMENT '消息类型  1 文字 2 语音',
      `audio_url` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '语音文件地址',
      `message_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'dify返回的消息id',
      `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
      PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI聊天记录表';