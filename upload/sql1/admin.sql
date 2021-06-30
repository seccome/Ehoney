/*
Navicat MySQL Data Transfer

Source Server Version : 50616

Target Server Type    : MYSQL
Target Server Version : 50616
File Encoding         : 65001

Date: 2021-06-18 12:11:50
*/

SET FOREIGN_KEY_CHECKS=0;

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
-- Records of e_admin
-- ----------------------------
INSERT INTO `e_admin` VALUES ('1', 'admin', 'admin', '2', null, null, '5');
