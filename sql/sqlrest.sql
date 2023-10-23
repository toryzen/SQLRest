#
# SQL Export


CREATE DATABASE IF NOT EXISTS `sqlrest` DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;
USE `sqlrest`;




SET @PREVIOUS_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS = 0;


DROP TABLE IF EXISTS `project`;
DROP TABLE IF EXISTS `db_source`;
DROP TABLE IF EXISTS `data_api`;
DROP TABLE IF EXISTS `auth_key`;


CREATE TABLE `auth_key` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `authkey` varchar(128) DEFAULT NULL,
  `project_id` varchar(32) DEFAULT '0',
  `api_ids` varchar(1024) NOT NULL DEFAULT '',
  `memo` varchar(128) DEFAULT '',
  `created_stime` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified_stime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_user` varchar(128) DEFAULT NULL,
  `modified_user` varchar(128) DEFAULT NULL,
  `is_del` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `authkey` (`authkey`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;


CREATE TABLE `data_api` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `api_id` varchar(8) DEFAULT '',
  `apiname` varchar(128) NOT NULL DEFAULT '',
  `project_id` varchar(8) DEFAULT '0',
  `db_id` varchar(128) NOT NULL DEFAULT '',
  `joint` varchar(128) DEFAULT '',
  `sourcesql` text NOT NULL,
  `memo` varchar(128) DEFAULT '',
  `created_stime` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified_stime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_user` varchar(128) DEFAULT NULL,
  `modified_user` varchar(128) DEFAULT NULL,
  `is_del` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq` (`api_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=89 DEFAULT CHARSET=utf8;


CREATE TABLE `db_source` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `db_id` varchar(8) NOT NULL DEFAULT '',
  `dbname` varchar(128) NOT NULL DEFAULT '',
  `ip` varchar(128) NOT NULL DEFAULT '',
  `port` int(11) NOT NULL DEFAULT '0',
  `user` varchar(128) NOT NULL DEFAULT '',
  `pwd` varchar(128) NOT NULL DEFAULT '',
  `memo` varchar(128) DEFAULT NULL,
  `created_stime` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified_stime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_user` varchar(128) DEFAULT NULL,
  `modified_user` varchar(128) DEFAULT NULL,
  `is_del` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq` (`db_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;


CREATE TABLE `project` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `project_id` varchar(8) DEFAULT '',
  `project_name` varchar(128) DEFAULT '',
  `memo` varchar(128) DEFAULT '',
  `created_stime` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified_stime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_user` varchar(128) DEFAULT NULL,
  `modified_user` varchar(128) DEFAULT NULL,
  `is_del` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq` (`project_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=54 DEFAULT CHARSET=utf8;




SET FOREIGN_KEY_CHECKS = @PREVIOUS_FOREIGN_KEY_CHECKS;


SET @PREVIOUS_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS = 0;


LOCK TABLES `auth_key` WRITE;
ALTER TABLE `auth_key` DISABLE KEYS;
INSERT INTO `auth_key` (`id`, `authkey`, `project_id`, `api_ids`, `memo`, `created_stime`, `modified_stime`, `created_user`, `modified_user`, `is_del`) VALUES 
	(1,'101e1c66e7b2e821','5f51665a','-1','管理密钥','2023-05-30 17:22:33','2023-10-23 14:31:15','toryzen','toryzen',0);
ALTER TABLE `auth_key` ENABLE KEYS;
UNLOCK TABLES;


LOCK TABLES `data_api` WRITE;
ALTER TABLE `data_api` DISABLE KEYS;
INSERT INTO `data_api` (`id`, `api_id`, `apiname`, `project_id`, `db_id`, `joint`, `sourcesql`, `memo`, `created_stime`, `modified_stime`, `created_user`, `modified_user`, `is_del`) VALUES 
	(1,'c616e72a','获取项目列表','5f51665a','616e7b2e','','select * from project a  left join  \n(select count(1) as allcnt from project where is_del = 0 ) b \non 1=1 where is_del = 0 {limit}','','2023-05-09 15:36:06','2023-10-23 14:11:39','toryzen','toryzen',0),
	(2,'e8e8181r','获取DB列表','5f51665a','616e7b2e','','select * from db_source a  left join  \n(select count(1) as allcnt from db_source where is_del = 0 ) b \non 1=1 where is_del = 0  {limit}','','2023-05-09 15:36:06','2023-05-30 17:22:13','toryzen','toryzen',0),
	(3,'e8181ee2','获取AuthKey列表','5f51665a','616e7b2e','','select id,authkey,(select project_name from project where project_id = a.project_id) project_name,project_id,api_ids,memo,created_stime,created_user,modified_stime,modified_user,allcnt from auth_key a  left join  \n(select count(1) as allcnt from auth_key where is_del = 0 ) b \non 1=1 where is_del = 0 {limit}','','2023-05-09 15:36:06','2023-05-30 17:22:14','toryzen','toryzen',0),
	(4,'8181e1ca','获取API列表','5f51665a','616e7b2e','','select id,api_id,apiname,(select project_name from project where project_id = a.project_id) project_name,project_id,\n(select dbname from db_source where db_id = a.db_id) dbname,db_id,joint,sourcesql,memo,created_stime,created_user,modified_stime,modified_user,allcnt\n  from data_api a  left join  \n(select count(1) as allcnt from data_api where is_del = 0 ) b \non 1=1 where is_del = 0 {limit}','','2023-05-09 15:36:06','2023-05-30 17:22:15','toryzen','toryzen',0),
	(5,'611c6164','新增项目','5f51665a','616e7b2e','','insert into project (project_id, project_name,memo,created_user,modified_user) values (LEFT(MD5(UUID()), 16),\'{project_name}\',\'{memo}\',\'{username}-{employeeid}\',\'{username}-{employeeid}\')','','2023-05-09 15:42:07','2023-05-30 17:22:15','toryzen','toryzen',0),
	(6,'e7b2e813','新增DB源','5f51665a','616e7b2e','','insert into db_source (db_id,dbname,ip,port,user,pwd,memo,created_user,modified_user) values (LEFT(MD5(UUID()), 16),\'{dbname}\',\'{ip}\',\'{port}\',\'{user}\',\'{pwd}\',\'{memo}\',\'{username}-{employeeid}\',\'{username}-{employeeid}\')','','2023-05-09 15:46:27','2023-05-30 17:22:16','toryzen','toryzen',0),
	(7,'81e616e3','新增AuthKey','5f51665a','616e7b2e','','insert into auth_key (authkey,project_id,api_ids,memo,created_user,modified_user) values (LEFT(MD5(UUID()), 16),\'{project_id}\',\'{api_ids}\',\'{memo}\',\'{username}-{employeeid}\',\'{username}-{employeeid}\')','','2023-05-09 15:50:54','2023-05-30 17:22:17','toryzen','toryzen',0),
	(8,'7b2e8161','新增API','5f51665a','616e7b2e','','insert into data_api (api_id,apiname,project_id,db_id,joint,sourcesql,memo,created_user,modified_user) values (LEFT(MD5(UUID()), 16),\'{apiname}\',\'{project_id}\',\'{db_id}\',\'{joint}\',\'{sourcesql}\',\'{memo}\',\'{username}-{employeeid}\',\'{username}-{employeeid}\')','','2023-05-09 15:53:04','2023-05-30 17:22:17','toryzen','toryzen',0),
	(9,'e7b2e812','更新项目信息','5f51665a','616e7b2e','','update project set project_name=\'{project_name}\' , memo = \'{memo}\' , modified_user = \'{username}-{employeeid}\' where id = {id}','','2023-05-09 20:09:02','2023-05-30 17:22:18','toryzen','toryzen',0),
	(10,'1e1c616e','更新DB源信息','5f51665a','616e7b2e','','update db_source set dbname=\'{dbname}\' , ip = \'{ip}\' , port=\'{port}\' , user=\'{user}\' , pwd=\'{pwd}\' , memo = \'{memo}\' , modified_user = \'{username}-{employeeid}\' where id = \'{id}\'','','2023-05-09 20:11:00','2023-05-30 17:22:19','toryzen','toryzen',0),
	(11,'716e71c6','更新AuthKey信息','5f51665a','616e7b2e','','update auth_key set api_ids=\'{api_ids}\' , memo = \'{memo}\'  , modified_user = \'{username}-{employeeid}\' where authkey = \'{authkey}\'','','2023-05-09 20:12:54','2023-05-30 17:22:19','toryzen','toryzen',0),
	(12,'1616e7b2','更新API信息','5f51665a','616e7b2e','','update data_api set apiname=\'{apiname}\' ,joint=\'{joint}\', sourcesql = \'{sourcesql}\' , memo = \'{memo}\'  , modified_user = \'{username}-{employeeid}\' where id = \'{id}\'','','2023-05-09 20:14:48','2023-05-30 17:22:20','toryzen','toryzen',0),
	(13,'e8181ee7','删除项目','5f51665a','616e7b2e','','update project set is_del=\'1\' where id = \'{project_id}\'','','2023-05-09 20:17:56','2023-05-30 17:22:21','toryzen','toryzen',0),
	(14,'b2e81aca','删除DB源','5f51665a','616e7b2e','','update db_source set is_del=1 where id = \'{db_id}\'','','2023-05-09 20:18:35','2023-05-30 17:22:21','toryzen','toryzen',0),
	(15,'81acavb2','删除AuthKey','5f51665a','616e7b2e','','update auth_key set is_del=\'1\' where authkey = \'{auth_key}\'','','2023-05-09 20:19:10','2023-05-30 17:22:22','toryzen','toryzen',0),
	(16,'ab2e81aa','删除API','5f51665a','616e7b2e','','update data_api set is_del=\'1\' where id = \'{api_id}\'','','2023-05-09 20:19:22','2023-05-30 17:22:23','toryzen','toryzen',0);
ALTER TABLE `data_api` ENABLE KEYS;
UNLOCK TABLES;


LOCK TABLES `db_source` WRITE;
ALTER TABLE `db_source` DISABLE KEYS;
INSERT INTO `db_source` (`id`, `db_id`, `dbname`, `ip`, `port`, `user`, `pwd`, `memo`, `created_stime`, `modified_stime`, `created_user`, `modified_user`, `is_del`) VALUES 
	(1,'616e7b2e','sqlrest','127.0.0.1',3306,'root','root','','2023-05-09 15:46:51','2023-05-30 17:22:46','toryzen','toryzen',0);
ALTER TABLE `db_source` ENABLE KEYS;
UNLOCK TABLES;


LOCK TABLES `project` WRITE;
ALTER TABLE `project` DISABLE KEYS;
INSERT INTO `project` (`id`, `project_id`, `project_name`, `memo`, `created_stime`, `modified_stime`, `created_user`, `modified_user`, `is_del`) VALUES 
	(1,'5f51665a','SQLRest','通用数据网关','2023-05-09 15:35:08','2023-05-30 17:22:57','toryzen','toryzen',0);
ALTER TABLE `project` ENABLE KEYS;
UNLOCK TABLES;




SET FOREIGN_KEY_CHECKS = @PREVIOUS_FOREIGN_KEY_CHECKS;


