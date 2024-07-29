# ************************************************************
# Sequel Pro SQL dump
# Version 5425
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.3.0)
# Database: glot
# Generation Time: 2024-07-29 09:10:15 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table tb_menu
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tb_menu`;

CREATE TABLE `tb_menu` (
                           `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                           `parent_id` bigint NOT NULL DEFAULT '0' COMMENT 'pid',
                           `mode` tinyint NOT NULL DEFAULT '1' COMMENT '类型 1-目录 2-菜单 3-按扭',
                           `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
                           `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '编码',
                           `route_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '路由名称',
                           `route_path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '路由路径',
                           `component` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页面组件',
                           `meta` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '资源元数据',
                           `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
                           `sort` int NOT NULL DEFAULT '1' COMMENT '排序',
                           `created_by` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '创建人',
                           `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_by` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '更新人',
                           `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `uq_code` (`code`),
                           KEY `idx_pid` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单表';

LOCK TABLES `tb_menu` WRITE;
/*!40000 ALTER TABLE `tb_menu` DISABLE KEYS */;

INSERT INTO `tb_menu` (`id`, `parent_id`, `mode`, `name`, `code`, `route_name`, `route_path`, `component`, `meta`, `status`, `sort`, `created_by`, `created_at`, `updated_by`, `updated_at`)
VALUES
    (119,0,2,'首页','home','home','/home','layout.base$view.home','{\"i18nKey\":\"route.home\",\"iconType\":1,\"icon\":\"mdi:monitor-dashboard\",\"layout\":\"\",\"page\":\"\",\"order\":1,\"pathParam\":\"\"}',1,1,'107','2024-07-24 12:38:59','','2024-07-24 12:38:59'),
    (120,0,1,'系统管理','manage','manage','/manage','layout.base','{\"i18nKey\":\"route.manage\",\"iconType\":1,\"icon\":\"carbon:cloud-service-management\",\"order\":2}',1,2,'107','2024-07-24 12:46:51','107','2024-07-29 14:41:52'),
    (121,120,2,'用户管理','manage_user','manage_user','/manage/user','view.manage_user','{\"i18nKey\":\"route.manage_user\",\"iconType\":1,\"icon\":\"ic:round-manage-accounts\",\"order\":3}',1,3,'107','2024-07-24 12:47:57','107','2024-07-29 14:52:51'),
    (122,120,2,'角色管理','manage_role','manage_role','/manage/role','view.manage_role','{\"i18nKey\":\"route.manage_role\",\"iconType\":1,\"icon\":\"carbon:user-role\",\"layout\":\"\",\"page\":\"\",\"order\":4,\"pathParam\":\"\"}',1,4,'107','2024-07-24 12:48:58','','2024-07-24 12:48:58'),
    (123,120,2,'菜单管理','manage_menu','manage_menu','/manage/menu','view.manage_menu','{\"i18nKey\":\"route.manage_menu\",\"iconType\":1,\"icon\":\"material-symbols:route\",\"layout\":\"\",\"page\":\"\",\"order\":5,\"pathParam\":\"\"}',1,5,'107','2024-07-24 12:49:49','','2024-07-24 12:49:49'),
    (127,0,2,'租户管理','tenant','tenant','/tenant','layout.base$view.tenant','{\"i18nKey\":\"route.tenant\",\"iconType\":1,\"icon\":\"ic:twotone-business\",\"order\":6}',1,6,'107','2024-07-25 19:46:28','107','2024-07-26 17:13:46'),
    (130,121,3,'用户租户','user-tenant','','','','',1,0,'107','2024-07-29 14:52:51','','2024-07-29 14:52:51');

/*!40000 ALTER TABLE `tb_menu` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tb_role
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tb_role`;

CREATE TABLE `tb_role` (
                           `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                           `name` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色名',
                           `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色编码',
                           `desc` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '描述',
                           `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
                           `created_by` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '创建人',
                           `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_by` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '更新人',
                           `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `uq_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色表';

LOCK TABLES `tb_role` WRITE;
/*!40000 ALTER TABLE `tb_role` DISABLE KEYS */;

INSERT INTO `tb_role` (`id`, `name`, `code`, `desc`, `status`, `created_by`, `created_at`, `updated_by`, `updated_at`)
VALUES
    (101,'超级管理员','R_SUPER','超级管理员',1,'admin','2024-07-18 15:19:04','','2024-07-29 11:36:56'),
    (102,'管理员','R_ADMIN','管理员',1,'admin','2024-07-18 15:20:03','','2024-07-29 11:37:01'),
    (103,'普通用户','R_USER','普通用户',1,'admin','2024-07-18 15:20:20','','2024-07-29 11:37:04'),
    (105,'测试用户','R_TEST','测试用户',1,'101','2024-07-19 18:02:12','107','2024-07-24 12:57:38');

/*!40000 ALTER TABLE `tb_role` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tb_role_menu
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tb_role_menu`;

CREATE TABLE `tb_role_menu` (
                                `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                                `role_id` bigint NOT NULL COMMENT '角色id',
                                `menu_id` bigint NOT NULL COMMENT '菜单按钮id',
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色菜单表';

LOCK TABLES `tb_role_menu` WRITE;
/*!40000 ALTER TABLE `tb_role_menu` DISABLE KEYS */;

INSERT INTO `tb_role_menu` (`id`, `role_id`, `menu_id`)
VALUES
    (84,101,120),
    (85,101,121),
    (86,101,122),
    (87,101,123),
    (88,101,127),
    (89,101,119),
    (93,103,119),
    (105,101,130),
    (106,102,119),
    (107,105,119),
    (108,105,120),
    (109,105,121),
    (110,105,127),
    (111,105,122);

/*!40000 ALTER TABLE `tb_role_menu` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tb_tenant
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tb_tenant`;

CREATE TABLE `tb_tenant` (
                             `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                             `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '租户编码',
                             `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '租户名称',
                             `contacts` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '联系人',
                             `phone` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
                             `address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '公司地址',
                             `size` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '公司规模',
                             `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1-正常 2-冻结',
                             `created_by` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '创建人',
                             `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             `updated_by` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '更新人',
                             `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `uq_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='租户表';

LOCK TABLES `tb_tenant` WRITE;
/*!40000 ALTER TABLE `tb_tenant` DISABLE KEYS */;

INSERT INTO `tb_tenant` (`id`, `code`, `name`, `contacts`, `phone`, `address`, `size`, `status`, `created_by`, `created_at`, `updated_by`, `updated_at`)
VALUES
    (1,'T1001','系统默认','admin','13800135225','','10',1,'admin','2024-07-26 18:59:49','107','2024-07-29 11:19:25'),
    (2,'T1002','测试业务','小张','15837382938','','',1,'107','2024-07-29 11:18:00','107','2024-07-29 13:29:06');

/*!40000 ALTER TABLE `tb_tenant` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tb_tenant_role
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tb_tenant_role`;

CREATE TABLE `tb_tenant_role` (
                                  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                                  `tenant_id` bigint NOT NULL COMMENT '租户id',
                                  `role_id` bigint NOT NULL COMMENT '角色id',
                                  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='租户角色表';

LOCK TABLES `tb_tenant_role` WRITE;
/*!40000 ALTER TABLE `tb_tenant_role` DISABLE KEYS */;

INSERT INTO `tb_tenant_role` (`id`, `tenant_id`, `role_id`)
VALUES
    (14,1,102),
    (15,1,101),
    (16,1,105),
    (17,1,103),
    (22,2,102),
    (23,2,103);

/*!40000 ALTER TABLE `tb_tenant_role` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tb_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tb_user`;

CREATE TABLE `tb_user` (
                           `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                           `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
                           `password` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
                           `phone` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
                           `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',
                           `nickname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
                           `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
                           `tenant_id` bigint NOT NULL COMMENT '租户id',
                           `created_by` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '创建人',
                           `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_by` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '更新人',
                           `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `uq_username` (`username`),
                           KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

LOCK TABLES `tb_user` WRITE;
/*!40000 ALTER TABLE `tb_user` DISABLE KEYS */;

INSERT INTO `tb_user` (`id`, `username`, `password`, `phone`, `email`, `nickname`, `status`, `tenant_id`, `created_by`, `created_at`, `updated_by`, `updated_at`)
VALUES
    (107,'super','123456','15033292333','lsk@tt.com','lsk',1,1,'101','2024-07-23 18:48:02','107','2024-07-29 11:20:07'),
    (108,'admin','123456','13598750284','lxq@tt.com','lxq',1,1,'101','2024-07-23 18:48:49','','2024-07-29 11:20:11'),
    (109,'user','123456','15895857041','lpp@tt.com','lpp',1,1,'101','2024-07-23 18:49:37','107','2024-07-29 15:08:00'),
    (111,'test','123456','13800136778','test@tt.com','aaaa',1,1,'107','2024-07-23 22:13:11','107','2024-07-29 15:08:35'),
    (113,'t1002','123456','13522333223','aa@tt.com','aabb',1,2,'107','2024-07-29 16:02:45','107','2024-07-29 16:04:57');

/*!40000 ALTER TABLE `tb_user` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table tb_user_role
# ------------------------------------------------------------

DROP TABLE IF EXISTS `tb_user_role`;

CREATE TABLE `tb_user_role` (
                                `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
                                `user_id` bigint NOT NULL COMMENT '用户id',
                                `role_id` bigint NOT NULL COMMENT '角色id',
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户角色表';

LOCK TABLES `tb_user_role` WRITE;
/*!40000 ALTER TABLE `tb_user_role` DISABLE KEYS */;

INSERT INTO `tb_user_role` (`id`, `user_id`, `role_id`)
VALUES
    (26,108,102),
    (31,107,101),
    (38,109,103),
    (39,111,105),
    (41,113,103);

/*!40000 ALTER TABLE `tb_user_role` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
