/*
Navicat MySQL Data Transfer

Source Server Version : 50616

Target Server Type    : MYSQL
Target Server Version : 50616
File Encoding         : 65001

Date: 2021-06-18 12:12:17
*/

SET FOREIGN_KEY_CHECKS=0;

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
-- Records of config
-- ----------------------------