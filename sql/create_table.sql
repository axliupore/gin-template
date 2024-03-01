create database gin_template;

use gin_template;

create table if not exists user
(
    id          bigint auto_increment primary key comment 'id',
    account     varchar(256)                          not null comment '账号',
    username    varchar(256)                          null comment '用户名',
    password    varchar(512)                          null comment '密码',
    avatar      varchar(1024)                         null comment '头像地址',
    email       varchar(128)                          null comment '邮箱',
    phone       varchar(128)                          null comment '电话',
    profile     varchar(2048)                         null comment '用户简介',
    gender      tinyint     default 0                 not null comment '男性:0, 女性:1',
    role        varchar(16) default 'user'            not null comment '用户角色:user, admin, ban...',
    create_time datetime    default CURRENT_TIMESTAMP not null comment '创建时间',
    update_time datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    is_delete   tinyint     default 0                 not null comment '是否删除',
    index idx_user_id (id),
    index idx_user_account (account),
    index idx_user_username (username)
) comment '用户表' collate = utf8mb4_unicode_ci
                   engine = InnoDB;