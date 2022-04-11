USE `d_meta`;

-- lock
DROP TABLE IF EXISTS `t_lock`;
CREATE TABLE IF NOT EXISTS `t_lock` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) DEFAULT '' COMMENT '锁名称，一般是服务名',
  `version` bigint(20) DEFAULT 1 COMMENT '版本号',
  `update_time` varchar(32) DEFAULT "" COMMENT '最近更新时间',
  `env` varchar(32) DEFAULT "" COMMENT '环境信息，如本机ip',
  PRIMARY KEY (`id`),
  KEY `name_index` (`name`),
  KEY `time_index` (`update_time`),
  KEY `env_index` (`env`),
  UNIQUE(name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1 COMMENT="保证集群条件下只有一个节点启动服务的乐观锁表";