#######
table structure file
#######

CREATE DATABASE `dms` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE dms;
drop table if exists dms.user_whitelist;
create table user_whitelist(
    id int not null auto_increment,
    username varchar(20) NOT NULL,
    is_deleted tinyint not null default 1 comment '1 yes,0 no',
    created datetime not null default current_timestamp,
    updated datetime not null default current_timestamp on update current_timestamp,
    PRIMARY KEY(ID),
    unique key uidx_username(username)
)engine=innodb default charset=utf8;

insert into user_whitelist(username,is_deleted)values("admin",0);