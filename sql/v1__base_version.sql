CREATE TABLE `t_sys_token_info`
(
    `id`               int(11)      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `token`            varchar(100) NOT NULL COMMENT 'token',
    `app_id`           int(11)               DEFAULT NULL COMMENT '应用ID',
    `remark`           varchar(200)          DEFAULT NULL COMMENT '备注',
    `del_flag`         tinyint(4)   NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
    `create_time`      datetime(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `last_update_time` datetime(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '最后更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx01_token` (`token`),
    KEY `idx02_app_id` (`app_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='token表';

CREATE TABLE `t_mysql_operation_info`
(
    `id`               int(11)      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `operation_type`   tinyint(4)   NOT NULL COMMENT '操作类型: 1-新增, 2-删除',
    `addrs`            varchar(200) NOT NULL COMMENT 'MySQL实例地址列表',
    `status`           tinyint(4)   NOT NULL DEFAULT '0' COMMENT '运行状态: 0-未运行, 1-运行中, 2-已完成, 3-已失败',
    `message`          mediumtext            DEFAULT NULL COMMENT '运行日志',
    `del_flag`         tinyint(4)   NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
    `create_time`      datetime(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `last_update_time` datetime(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '最后更新时间',
    PRIMARY KEY (`id`),
    KEY `idx01_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT = 'MySQL操作信息表';

CREATE TABLE `t_mysql_operation_detail`
(
    `id`               int(11)      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `operation_id`     int(11)      NOT NULL COMMENT '操作ID',
    `host_ip`          varchar(100) NOT NULL COMMENT 'MySQL服务器IP',
    `port_num`         int(11)      NOT NULL COMMENT 'MySQL服务器端口',
    `status`           tinyint(4)   NOT NULL DEFAULT '0' COMMENT '运行状态: 0-未运行, 1-运行中, 2-已完成, 3-已失败',
    `message`          mediumtext            DEFAULT NULL COMMENT '运行日志',
    `del_flag`         tinyint(4)   NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
    `create_time`      datetime(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `last_update_time` datetime(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '最后更新时间',
    PRIMARY KEY (`id`),
    KEY `idx01_operation_id_status` (`operation_id`, `status`),
    KEY `idx02_host_ip_port_num` (`host_ip`, `port_num`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT = 'MySQL操作流水表';

CREATE TABLE `t_mysql_operation_lock`
(
    `id`               int(11)      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `addr`             varchar(100) NOT NULL COMMENT 'MySQL实例监听地址',
    `del_flag`         tinyint(4)   NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
    `create_time`      datetime(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `last_update_time` datetime(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '最后更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx01_addr` (`addr`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT = 'MySQL操作锁表';

