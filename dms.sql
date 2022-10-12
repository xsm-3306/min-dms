#######
table structure file
#######

CREATE DATABASE `dms` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE dms;
drop table if exists dms.user_whitelist;
CREATE TABLE `user_whitelist` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(20) NOT NULL,
  `is_deleted` tinyint(4) NOT NULL DEFAULT 1 COMMENT '1 yes,0 no',
  `created` datetime NOT NULL DEFAULT current_timestamp(),
  `updated` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

insert into user_whitelist(username,is_deleted)values("admin",0);

drop table if exists user_sqlexec_log;
CREATE TABLE `user_sqlexec_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL COMMENT 'user id reference tab user_whitelist primary key',
  `created` datetime NOT NULL DEFAULT current_timestamp(),
  `exec_result` varchar(50) DEFAULT NULL COMMENT '执行结果',
  `rows_inserted` int(11) DEFAULT NULL,
  `rows_updated` int(11) DEFAULT NULL,
  `rows_deleted` int(11) DEFAULT NULL,
  `recovery_id` varchar(50) DEFAULT NULL COMMENT '回退id，也作为回退文件名',
  PRIMARY KEY (`id`),
  KEY `idx_recovery_id` (`recovery_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;