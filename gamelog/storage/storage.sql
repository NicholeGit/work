/*
Navicat MySQL Data Transfer

Source Server         : 10.100.0.181
Source Server Version : 50173
Source Host           : 10.100.0.181:3306
Source Database       : testStorage

Target Server Type    : MYSQL
Target Server Version : 50173
File Encoding         : 65001

Date: 2015-03-18 13:47:39
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `storage`
-- ----------------------------
DROP TABLE IF EXISTS `storage`;
CREATE TABLE `storage` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `account` varchar(255) NOT NULL COMMENT '账号\r\n',
  `objectId` int(11) NOT NULL COMMENT '宝物ID',
  `count` int(11) NOT NULL COMMENT '宝物使用数量',
  `type` int(11) NOT NULL COMMENT '宝物类型 1：普通仓库 2：vip仓库',
  `created` datetime NOT NULL COMMENT '修改日期',
  PRIMARY KEY (`id`),
  KEY `objectId` (`objectId`)
) ENGINE=MyISAM AUTO_INCREMENT=213017 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of storage
-- ----------------------------
