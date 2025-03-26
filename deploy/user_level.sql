/*
Navicat MySQL Data Transfer

Source Server         : mysql57@lo
Source Server Version : 50729
Source Host           : localhost:3306
Source Database       : flexdrive

Target Server Type    : MYSQL
Target Server Version : 50729
File Encoding         : 65001

Date: 2024-03-28 09:50:20
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `user_level`
-- ----------------------------
DROP TABLE IF EXISTS `user_level`;
CREATE TABLE `user_level` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '会员级别id',
  `name` varchar(32) NOT NULL COMMENT '会员级别名称',
  `quotaSpace` bigint(20) NOT NULL DEFAULT '0' COMMENT '会员级别的会员空间 单位KB',
  `price` int(11) NOT NULL DEFAULT '0' COMMENT '会员级别的价格 单位分',
  `isDefault` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否新会员默认等级 1=是',
  `created` timestamp NOT NULL DEFAULT '1999-01-01 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=启用 99=删除',
  `primeCost` int(11) NOT NULL DEFAULT '0' COMMENT '原价 仅用于展示 单位分',
  `desc` varchar(256) NOT NULL DEFAULT '' COMMENT '介绍',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of user_level
-- ----------------------------
INSERT INTO `user_level` VALUES ('1', '青铜会员', '512000', '0', '1', '2024-03-28 09:40:46', '2024-03-28 09:49:15', '1', '100', '空间500M');
INSERT INTO `user_level` VALUES ('2', '白银会员', '5120000', '300', '0', '2024-03-28 09:41:41', '2024-03-28 09:49:15', '1', '500', '空间5G');
INSERT INTO `user_level` VALUES ('3', '黄金会员', '51200000', '1000', '0', '2024-03-28 09:42:47', '2024-03-28 09:49:15', '1', '5000', '空间50G');
INSERT INTO `user_level` VALUES ('4', '钻石会员', '512000000', '5000', '0', '2024-03-28 09:43:25', '2024-03-28 09:49:15', '1', '50000', '空间500G');
