/*
Source Database       : sec_ehoneypot

Target Server Type    : MYSQL
Target Server Version : 50616
File Encoding         : 65001

Date: 2021-06-07 19:59:40
*/

CREATE DATABASE IF NOT EXISTS sec_ehoneypot default charset utf8 COLLATE utf8_general_ci;

use sec_ehoneypot;

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for attacklog
-- ----------------------------
DROP TABLE IF EXISTS `attacklog`;
CREATE TABLE `attacklog` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `srchost` varchar(128) DEFAULT NULL COMMENT 'src_host',
  `srcport` int(10) DEFAULT NULL COMMENT 'src_port',
  `serverid` varchar(50) DEFAULT NULL,
  `honeypotid` varchar(45) DEFAULT NULL,
  `honeytypeid` varchar(45) DEFAULT NULL,
  `honeypotport` int(10) DEFAULT NULL COMMENT 'dst_port',
  `attackip` varchar(128) DEFAULT NULL,
  `attackport` int(10) DEFAULT NULL,
  `country` varchar(45) DEFAULT NULL,
  `province` varchar(100) DEFAULT NULL,
  `attacktime` varchar(128) DEFAULT NULL COMMENT 'local_time',
  `eventdetail` longtext COMMENT 'logdata',
  `proxytype` varchar(45) DEFAULT NULL COMMENT '代理类型',
  `sourcetype` int(11) DEFAULT NULL COMMENT '1是proxy（代理）、2是falco，区分日志来源',
  `logdata` longtext,
  `longitude` varchar(64) DEFAULT NULL,
  `latitude` varchar(64) DEFAULT NULL,
  `exportport` int(10) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1376 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for baits
-- ----------------------------
DROP TABLE IF EXISTS `baits`;
CREATE TABLE `baits` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `baittype` varchar(128) DEFAULT NULL,
  `baitname` varchar(128) DEFAULT NULL,
  `createtime` varchar(128) DEFAULT NULL,
  `creator` varchar(50) DEFAULT NULL,
  `baitid` varchar(50) DEFAULT NULL,
  `baitsystype` varchar(50) DEFAULT NULL,
  `baitinfo` varchar(50) DEFAULT NULL,
  `md5` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=83 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of baits
-- ----------------------------
INSERT INTO `baits` VALUES ('1', 'file', 'sql2', '1622786265', 'admin', '81d0b0ee-b6dc-4270-89cb-84022ab59d44', 'cdc872db616ac66adb3166c75e9ad183', 'config.sql', '');
INSERT INTO `baits` VALUES ('2', 'file', 'sql1', '1622792498', 'admin', '716cc2d0-4ea8-4708-9b44-6364532e2b9a', 'cdc872db616ac66adb3166c75e9ad183', 'admin.sql', '');

-- ----------------------------
-- Table structure for baittype
-- ----------------------------
DROP TABLE IF EXISTS `baittype`;
CREATE TABLE `baittype` (
  `id` int(64) NOT NULL AUTO_INCREMENT,
  `baittype` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `confname` varchar(128) DEFAULT NULL,
  `confvalue` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=83 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for conf_redisinfo
-- ----------------------------
DROP TABLE IF EXISTS `conf_redisinfo`;
CREATE TABLE `conf_redisinfo` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `redisurl` varchar(128) DEFAULT NULL,
  `redisport` varchar(64) DEFAULT NULL,
  `password` varchar(128) DEFAULT NULL,
  `createtime` varchar(128) DEFAULT NULL,
  `redisid` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `redisid` (`redisid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for conf_traceinfo
-- ----------------------------
DROP TABLE IF EXISTS `conf_traceinfo`;
CREATE TABLE `conf_traceinfo` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tracehost` varchar(128) DEFAULT NULL,
  `traceid` varchar(128) DEFAULT NULL,
  `createtime` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `tracehostid` (`traceid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for desipconf
-- ----------------------------
DROP TABLE IF EXISTS `desipconf`;
CREATE TABLE `desipconf` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `desip` varchar(64) DEFAULT NULL,
  `longitude` varchar(64) DEFAULT NULL,
  `latitude` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for e_admin
-- ----------------------------
DROP TABLE IF EXISTS `e_admin`;
CREATE TABLE `e_admin` (
  `uid` int(5) NOT NULL AUTO_INCREMENT,
  `uname` varchar(128) DEFAULT NULL,
  `upass` varchar(128) DEFAULT NULL,
  `status` int(10) DEFAULT NULL,
  `locktime` varchar(64) DEFAULT NULL,
  `starttime` varchar(64) DEFAULT NULL,
  `errno` int(10) DEFAULT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for fowards
-- ----------------------------
DROP TABLE IF EXISTS `fowards`;
CREATE TABLE `fowards` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `taskid` varchar(64) DEFAULT NULL,
  `agentid` varchar(128) DEFAULT NULL,
  `honeytypeid` varchar(128) DEFAULT NULL,
  `forwardport` varchar(255) DEFAULT NULL COMMENT '业务服务器端口',
  `honeypotport` int(11) DEFAULT NULL COMMENT '蜜罐服务器端口',
  `serverid` varchar(64) DEFAULT NULL COMMENT 'serverid',
  `createtime` varchar(255) DEFAULT NULL,
  `offlinetime` varchar(255) DEFAULT NULL,
  `creator` varchar(32) DEFAULT NULL,
  `status` int(10) DEFAULT '0',
  `type` varchar(50) DEFAULT NULL,
  `path` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=133 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for harborinfo
-- ----------------------------
DROP TABLE IF EXISTS `harborinfo`;
CREATE TABLE `harborinfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `harborurl` varchar(128) DEFAULT NULL,
  `username` varchar(128) DEFAULT NULL,
  `password` varchar(128) DEFAULT NULL,
  `projectname` varchar(128) DEFAULT NULL,
  `createtime` varchar(128) DEFAULT NULL,
  `harborid` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `harborid` (`harborid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for honeyfowards
-- ----------------------------
DROP TABLE IF EXISTS `honeyfowards`;
CREATE TABLE `honeyfowards` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `taskid` varchar(45) DEFAULT NULL,
  `agentid` varchar(128) DEFAULT NULL,
  `serverid` varchar(45) DEFAULT NULL,
  `forwardport` varchar(255) DEFAULT NULL COMMENT '蜜罐服务器监听端口',
  `honeypotport` int(11) DEFAULT NULL COMMENT 'docker port',
  `honeypotid` varchar(48) DEFAULT NULL COMMENT 'honeypotid',
  `createtime` varchar(255) DEFAULT NULL,
  `offlinetime` varchar(255) DEFAULT NULL,
  `creator` varchar(50) DEFAULT NULL,
  `status` int(10) DEFAULT NULL,
  `type` varchar(50) DEFAULT NULL,
  `path` varchar(45) DEFAULT NULL,
  `forwardstatus` int(10) unsigned DEFAULT '1' COMMENT '界面选择该记录后，置为2',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=170 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for honeyimages
-- ----------------------------
DROP TABLE IF EXISTS `honeyimages`;
CREATE TABLE `honeyimages` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `imagename` varchar(128) DEFAULT NULL,
  `imageid` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for honeypots
-- ----------------------------
DROP TABLE IF EXISTS `honeypots`;
CREATE TABLE `honeypots` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `honeyname` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `podname` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `honeytypeid` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `honeypotid` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `honeyip` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `honeyport` int(50) DEFAULT NULL,
  `honeyimage` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `honeynamespce` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `createtime` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `creator` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `hostip` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `serverid` varchar(50) CHARACTER SET utf8 DEFAULT NULL COMMENT '集群id',
  `sysid` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `status` int(10) DEFAULT '1',
  `offlinetime` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `agentid` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `honeyname` (`honeyname`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=16943 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for honeypotservers
-- ----------------------------
DROP TABLE IF EXISTS `honeypotservers`;
CREATE TABLE `honeypotservers` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `servername` varchar(128) DEFAULT NULL,
  `serverip` varchar(512) DEFAULT NULL,
  `serverid` varchar(128) DEFAULT NULL,
  `status` int(10) DEFAULT NULL,
  `agentid` varchar(128) DEFAULT NULL,
  `regtime` varchar(128) DEFAULT NULL,
  `heartbeattime` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `agentid` (`agentid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=22022 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for honeypotstype
-- ----------------------------
DROP TABLE IF EXISTS `honeypotstype`;
CREATE TABLE `honeypotstype` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `honeypottype` varchar(128) DEFAULT NULL,
  `createtime` varchar(128) DEFAULT NULL,
  `softpath` varchar(128) DEFAULT NULL,
  `typeid` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `honeypottype` (`honeypottype`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for honeypot_bait
-- ----------------------------
DROP TABLE IF EXISTS `honeypot_bait`;
CREATE TABLE `honeypot_bait` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `taskid` varchar(64) DEFAULT NULL,
  `baitid` varchar(50) DEFAULT NULL,
  `baitinfo` varchar(128) DEFAULT NULL,
  `honeypotid` varchar(50) DEFAULT NULL,
  `data` longtext,
  `createtime` varchar(128) DEFAULT NULL,
  `offlinetime` varchar(128) DEFAULT NULL,
  `creator` varchar(50) DEFAULT NULL,
  `status` int(10) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=113 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for honeypot_sign
-- ----------------------------
DROP TABLE IF EXISTS `honeypot_sign`;
CREATE TABLE `honeypot_sign` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `taskid` varchar(64) DEFAULT NULL,
  `signid` varchar(64) DEFAULT NULL,
  `signinfo` varchar(128) DEFAULT NULL,
  `createtime` varchar(128) DEFAULT NULL,
  `offlinetime` varchar(128) DEFAULT NULL,
  `creator` varchar(50) DEFAULT NULL,
  `honeypotid` varchar(50) DEFAULT NULL,
  `status` int(10) DEFAULT '1',
  `tracecode` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for honeyserverconfig
-- ----------------------------
DROP TABLE IF EXISTS `honeyserverconfig`;
CREATE TABLE `honeyserverconfig` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `honeyserverid` varchar(64) DEFAULT NULL,
  `serversshkey` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for logsource
-- ----------------------------
DROP TABLE IF EXISTS `logsource`;
CREATE TABLE `logsource` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for podimage
-- ----------------------------
DROP TABLE IF EXISTS `podimage`;
CREATE TABLE `podimage` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `imageaddress` varchar(128) DEFAULT NULL,
  `repository` varchar(50) DEFAULT NULL,
  `imagename` varchar(128) DEFAULT NULL,
  `imageport` varchar(128) DEFAULT NULL,
  `imagetype` varchar(64) DEFAULT NULL,
  `imageos` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `imageaddress` (`imageaddress`,`imagename`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for servers
-- ----------------------------
DROP TABLE IF EXISTS `servers`;
CREATE TABLE `servers` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `servername` varchar(128) DEFAULT NULL,
  `serverip` varchar(512) DEFAULT NULL,
  `serverid` varchar(50) DEFAULT NULL,
  `agentid` varchar(128) DEFAULT NULL,
  `vpcname` varchar(50) DEFAULT NULL,
  `vpsowner` varchar(50) DEFAULT NULL,
  `regtime` varchar(50) DEFAULT NULL,
  `heartbeattime` varchar(50) DEFAULT NULL,
  `status` int(10) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `agentid` (`agentid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4405281 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for server_bait
-- ----------------------------
DROP TABLE IF EXISTS `server_bait`;
CREATE TABLE `server_bait` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `taskid` varchar(45) DEFAULT NULL,
  `type` varchar(45) DEFAULT NULL,
  `agentid` varchar(128) DEFAULT NULL,
  `baitid` varchar(128) DEFAULT NULL,
  `baitinfo` varchar(128) DEFAULT NULL,
  `createtime` varchar(255) DEFAULT NULL,
  `offlinetime` varchar(255) DEFAULT NULL,
  `creator` varchar(50) DEFAULT NULL,
  `status` int(10) DEFAULT '0',
  `data` longtext,
  `md5` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=139 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for server_sign
-- ----------------------------
DROP TABLE IF EXISTS `server_sign`;
CREATE TABLE `server_sign` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `taskid` varchar(64) DEFAULT NULL,
  `agentid` varchar(64) DEFAULT NULL,
  `signid` varchar(50) DEFAULT NULL,
  `signinfo` varchar(128) DEFAULT NULL,
  `createtime` varchar(128) DEFAULT NULL,
  `offlinetime` varchar(128) DEFAULT NULL,
  `creator` varchar(50) DEFAULT NULL,
  `serverid` varchar(50) DEFAULT NULL,
  `status` int(10) DEFAULT NULL,
  `md5` varchar(64) DEFAULT NULL,
  `type` varchar(64) DEFAULT NULL,
  `tracecode` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for signs
-- ----------------------------
DROP TABLE IF EXISTS `signs`;
CREATE TABLE `signs` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `signtype` varchar(128) DEFAULT NULL,
  `signname` varchar(128) DEFAULT NULL,
  `createtime` varchar(64) DEFAULT NULL,
  `creator` varchar(64) DEFAULT NULL,
  `signid` varchar(64) DEFAULT NULL,
  `signsystype` varchar(64) DEFAULT NULL,
  `signinfo` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `signname` (`signname`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of signs
-- ----------------------------
INSERT INTO `signs` VALUES ('1', 'file', 'pdfconfig', '1622734604', 'admin', 'dfc22caf-6702-4737-84a3-6e3c2e6e00f2', '', 'configbak.pdf');
INSERT INTO `signs` VALUES ('2', 'file', 'configlist', '1622775391', 'admin', '430e5732-8434-40a3-96a3-e55c958cbec1', '', 'configlist.xlsx');
INSERT INTO `signs` VALUES ('3', 'file', 'companyinfo', '1623162751', 'admin', 'db4f2337-3685-4a2f-b0e8-7ccc30e6ddab', '', 'companyinfo.pptx');
INSERT INTO `signs` VALUES ('4', 'file', 'configbak', '1623162751', 'admin', 'db4f2337-3685-4a2f-b0e8-7ccc30e6dda1', null, 'configbak.docx');


-- ----------------------------
-- Table structure for signtracemsg
-- ----------------------------
DROP TABLE IF EXISTS `signtracemsg`;
CREATE TABLE `signtracemsg` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tracecode` varchar(128) DEFAULT NULL,
  `openip` varchar(128) DEFAULT NULL,
  `opentime` varchar(128) DEFAULT NULL,
  `useragent` varchar(256) DEFAULT NULL,
  `longitude` varchar(64) DEFAULT NULL,
  `latitude` varchar(64) DEFAULT NULL,
  `ipcountry` varchar(64) DEFAULT NULL,
  `ipcity` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for signtype
-- ----------------------------
DROP TABLE IF EXISTS `signtype`;
CREATE TABLE `signtype` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `signtype` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for systemtype
-- ----------------------------
DROP TABLE IF EXISTS `systemtype`;
CREATE TABLE `systemtype` (
  `id` int(128) NOT NULL AUTO_INCREMENT,
  `systype` varchar(50) DEFAULT NULL,
  `sysid` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;






-- ----------------------------
-- Records of honeypotstype
-- ----------------------------
INSERT INTO `honeypotstype` VALUES ('1', 'http', '1622464884', '/home/ehoney_proxy/httpproxy', '80791b3ae7002cb88c246876d9faa8f8');
INSERT INTO `honeypotstype` VALUES ('2', 'telnet', '1622525671', '/home/ehoney_proxy/telnetproxy', '03583cd75bf401944b018f81b3f6916d');
INSERT INTO `honeypotstype` VALUES ('3', 'redis', '1622525708', '/home/ehoney_proxy/redisproxy', '86a1b907d54bf7010394bf316e183e67');
INSERT INTO `honeypotstype` VALUES ('4', 'ssh', '1622631563', '/home/ehoney_proxy/sshproxy', '1eb174fa332c502d2b4929d74e5d1d64');
INSERT INTO `honeypotstype` VALUES ('5', 'mysql', '1622690649', '/home/ehoney_proxy/mysqlproxy', '81c3b080dad537de7e10e0987a4bf52e');


-- ----------------------------
-- Records of baittype
-- ----------------------------
INSERT INTO `baittype` VALUES ('1', 'file');
INSERT INTO `baittype` VALUES ('2', 'history');


-- ----------------------------
-- Records of signtype
-- ----------------------------
INSERT INTO `signtype` VALUES ('1', 'file');


-- ----------------------------
-- Records of systemtype
-- ----------------------------
INSERT INTO `systemtype` VALUES ('1', 'Win7', 'f2886f2f800bc2328e4d13404c6d96b5');
INSERT INTO `systemtype` VALUES ('2', 'Win10', '6a29b3789f1c49a09cc58099b2e8258c');
INSERT INTO `systemtype` VALUES ('3', 'Centos', 'cdc872db616ac66adb3166c75e9ad183');
INSERT INTO `systemtype` VALUES ('4', 'Ubuntu', '1d41c853af58d3a7ae54990ce29417d8');

-- ----------------------------
-- Records of podimage
-- ----------------------------
-- ----------------------------
-- Records of podimage
-- ----------------------------
INSERT INTO `podimage` VALUES ('1', '47.96.71.197:90/ehoney/mysql:v1', null, 'ehoney/mysql', '3306', '81c3b080dad537de7e10e0987a4bf52e', 'cdc872db616ac66adb3166c75e9ad183');
INSERT INTO `podimage` VALUES ('2', '47.96.71.197:90/ehoney/redis:v1', null, 'ehoney/redis', '6379', '86a1b907d54bf7010394bf316e183e67', 'cdc872db616ac66adb3166c75e9ad183');
INSERT INTO `podimage` VALUES ('3', '47.96.71.197:90/ehoney/ssh:v1', null, 'ehoney/ssh', '22', '1eb174fa332c502d2b4929d74e5d1d64', 'cdc872db616ac66adb3166c75e9ad183');
INSERT INTO `podimage` VALUES ('4', '47.96.71.197:90/ehoney/telnet:v1', null, 'ehoney/telnet', '23', '03583cd75bf401944b018f81b3f6916d', 'cdc872db616ac66adb3166c75e9ad183');
INSERT INTO `podimage` VALUES ('5', '47.96.71.197:90/ehoney/tomcat:v1', null, 'ehoney/tomcat', '8080', '80791b3ae7002cb88c246876d9faa8f8', 'cdc872db616ac66adb3166c75e9ad183');


INSERT INTO `e_admin` VALUES ('1', 'admin', '123456', '1', null, null, '5');



