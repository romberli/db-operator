CREATE TABLE `t_mysql_installation_operation_info`
(
    `id`               int(11)       NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `addrs`            varchar(1000) NOT NULL COMMENT 'MySQL实例地址列表',
    `status`           tinyint(4)    NOT NULL DEFAULT '0' COMMENT '运行状态: 0-未运行, 1-运行中, 2-已完成, 3-已失败',
    `message`          mediumtext             DEFAULT NULL COMMENT '运行日志',
    `del_flag`         tinyint(4)    NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
    `create_time`      datetime(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `last_update_time` datetime(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '最后更新时间',
    PRIMARY KEY (`id`),
    KEY `idx01_create_time` (`create_time`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT = 'MySQL安装操作信息表';

CREATE TABLE `t_mysql_installation_operation_detail`
(
    `id`               int(11)       NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `operation_id`     varchar(100)  NOT NULL COMMENT '操作ID',
    `host_ip`          varchar(1000) NOT NULL COMMENT 'MySQL服务器IP',
    `port_num`         int(11)       NOT NULL COMMENT 'MySQL服务器端口',
    `status`           tinyint(4)    NOT NULL DEFAULT '0' COMMENT '运行状态: 0-未运行, 1-运行中, 2-已完成, 3-已失败',
    `message`          mediumtext             DEFAULT NULL COMMENT '运行日志',
    `del_flag`         tinyint(4)    NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
    `create_time`      datetime(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `last_update_time` datetime(6)   NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '最后更新时间',
    PRIMARY KEY (`id`),
    KEY `idx01_operation_id_status` (`operation_id`, `status`),
    KEY `idx02_host_ip_port_num` (`host_ip`, `port_num`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT = 'MySQL安装操作流水表';


