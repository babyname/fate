/*
 Navicat Premium Data Transfer

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 50719
 Source Host           : localhost:3306
 Source Schema         : fate

 Target Server Type    : MySQL
 Target Server Version : 50719
 File Encoding         : 65001

 Date: 09/01/2018 17:21:15
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for five_phase
-- ----------------------------
DROP TABLE IF EXISTS `five_phase`;
CREATE TABLE `five_phase`  (
  `id` varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `created` datetime(0) NULL DEFAULT NULL,
  `updated` datetime(0) NULL DEFAULT NULL,
  `deleted` datetime(0) NULL DEFAULT NULL,
  `version` int(11) NULL DEFAULT 1,
  `first` varchar(2) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `second` varchar(2) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `third` varchar(2) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `fortune` varchar(8) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of five_phase
-- ----------------------------
INSERT INTO `five_phase` VALUES ('a325b888-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '木', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a3278d47-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '木', '火', '大吉');
INSERT INTO `five_phase` VALUES ('a32b5e23-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '木', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a32c966d-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '木', '金', '凶多吉少');
INSERT INTO `five_phase` VALUES ('a32d80c7-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '木', '水', '吉多于凶');
INSERT INTO `five_phase` VALUES ('a32ee059-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '火', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a32f7c7a-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '火', '火', '中吉');
INSERT INTO `five_phase` VALUES ('a330b52a-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '火', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a331c683-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '火', '金', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a33262ca-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '火', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a3339b53-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '土', '木', '大凶');
INSERT INTO `five_phase` VALUES ('a3343771-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '土', '火', '中吉');
INSERT INTO `five_phase` VALUES ('a3357001-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '土', '土', '吉');
INSERT INTO `five_phase` VALUES ('a336816f-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '土', '金', '吉多于凶');
INSERT INTO `five_phase` VALUES ('a3371db7-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '土', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a3382f28-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '金', '木', '大凶');
INSERT INTO `five_phase` VALUES ('a338f28b-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '金', '火', '大凶');
INSERT INTO `five_phase` VALUES ('a339dcf0-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '金', '土', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a33aee49-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '金', '金', '大凶');
INSERT INTO `five_phase` VALUES ('a33bb1a8-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '金', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a33cea31-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '水', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a33e49b8-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '水', '火', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a33fd03b-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '水', '土', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a340ba99-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '水', '金', '大吉');
INSERT INTO `five_phase` VALUES ('a34af4d2-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '木', '水', '水', '大吉');
INSERT INTO `five_phase` VALUES ('a34e0111-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '木', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a34faeb2-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '木', '火', '大吉');
INSERT INTO `five_phase` VALUES ('a351355f-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '木', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a351f8b4-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '木', '金', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a35294f4-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '木', '水', '中吉');
INSERT INTO `five_phase` VALUES ('a353cd6f-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '火', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a35505f1-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '火', '火', '中吉');
INSERT INTO `five_phase` VALUES ('a355c950-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '火', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a35728c9-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '火', '金', '大凶');
INSERT INTO `five_phase` VALUES ('a3586409-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '火', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a35999ca-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '土', '木', '吉多于凶');
INSERT INTO `five_phase` VALUES ('a35a85ce-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '土', '火', '大吉');
INSERT INTO `five_phase` VALUES ('a35be3c0-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '土', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a35cf53d-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '土', '金', '大吉');
INSERT INTO `five_phase` VALUES ('a35e2db7-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:03', '2018-01-05 10:25:03', NULL, 1, '火', '土', '水', '吉多于凶');
INSERT INTO `five_phase` VALUES ('a35ef105-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '金', '木', '大凶');
INSERT INTO `five_phase` VALUES ('a3605089-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '金', '火', '大凶');
INSERT INTO `five_phase` VALUES ('a361b017-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '金', '土', '吉凶参半');
INSERT INTO `five_phase` VALUES ('a362c1a5-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '金', '金', '大凶');
INSERT INTO `five_phase` VALUES ('a36384e1-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '金', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a3647497-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '水', '木', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3658500-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '水', '火', '大凶');
INSERT INTO `five_phase` VALUES ('a3666b00-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '水', '土', '大凶');
INSERT INTO `five_phase` VALUES ('a367f1a2-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '水', '金', '大凶');
INSERT INTO `five_phase` VALUES ('a369517e-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '火', '水', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a36b4d0b-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '木', '木', '中吉');
INSERT INTO `five_phase` VALUES ('a36cfac8-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '木', '火', '中吉');
INSERT INTO `five_phase` VALUES ('a36e815f-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '木', '土', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3702f19-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '木', '金', '大凶');
INSERT INTO `five_phase` VALUES ('a371b5ad-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '木', '水', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a372a023-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '火', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a374c323-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '火', '火', '大吉');
INSERT INTO `five_phase` VALUES ('a376e5d1-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '火', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a378455c-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '火', '金', '吉多于凶');
INSERT INTO `five_phase` VALUES ('a379cc05-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '火', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a37b54c6-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '土', '木', '中吉');
INSERT INTO `five_phase` VALUES ('a37cd93b-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '土', '火', '大吉');
INSERT INTO `five_phase` VALUES ('a37d9c94-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '土', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a37eadfe-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '土', '金', '大吉');
INSERT INTO `five_phase` VALUES ('a3800d97-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '土', '水', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3816d11-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '金', '木', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3825764-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '金', '火', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a383b71d-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '金', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a3853daa-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '金', '金', '大吉');
INSERT INTO `five_phase` VALUES ('a3869d46-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '金', '水', '大吉');
INSERT INTO `five_phase` VALUES ('a387d5bd-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '水', '木', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a388e756-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '水', '火', '大凶');
INSERT INTO `five_phase` VALUES ('a389d18d-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '水', '土', '大凶');
INSERT INTO `five_phase` VALUES ('a38ba665-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '水', '金', '吉凶参半');
INSERT INTO `five_phase` VALUES ('a38ce41f-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '土', '水', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a38eb37b-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '木', '木', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3903a24-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '木', '火', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a391c0ca-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '木', '土', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a393205c-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '木', '金', '大凶');
INSERT INTO `five_phase` VALUES ('a39458f0-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '木', '水', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a39654bf-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '火', '木', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3976614-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '火', '火', '吉凶参半');
INSERT INTO `five_phase` VALUES ('a398c5b8-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '火', '土', '吉凶参半');
INSERT INTO `five_phase` VALUES ('a399d729-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '火', '金', '大凶');
INSERT INTO `five_phase` VALUES ('a39a9aa8-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '火', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a39babec-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '土', '木', '中吉');
INSERT INTO `five_phase` VALUES ('a39cbd8f-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '土', '火', '大吉');
INSERT INTO `five_phase` VALUES ('a39d5986-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '土', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a39e6b02-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '土', '金', '大吉');
INSERT INTO `five_phase` VALUES ('a39f2e61-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '土', '水', '吉多于凶');
INSERT INTO `five_phase` VALUES ('a3a066c7-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '金', '木', '大凶');
INSERT INTO `five_phase` VALUES ('a3a1edef-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '金', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a3a34d09-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '金', '金', '中吉');
INSERT INTO `five_phase` VALUES ('a3a45e69-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '金', '水', '中吉');
INSERT INTO `five_phase` VALUES ('a3a4fac6-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '水', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a3a6334b-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '水', '火', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3a74500-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '水', '土', '吉');
INSERT INTO `five_phase` VALUES ('a3a82ef4-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '水', '金', '大吉');
INSERT INTO `five_phase` VALUES ('a3a96784-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '金', '水', '水', '中吉');
INSERT INTO `five_phase` VALUES ('a3aaee3b-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '木', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a3ac0407-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '木', '火', '大吉');
INSERT INTO `five_phase` VALUES ('a3acea00-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '木', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a3ae2277-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '木', '金', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3af0ccb-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '木', '水', '大吉');
INSERT INTO `five_phase` VALUES ('a3b01e3c-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '火', '木', '中吉');
INSERT INTO `five_phase` VALUES ('a3b1a4db-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '火', '火', '大凶');
INSERT INTO `five_phase` VALUES ('a3b32b72-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '火', '土', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3b46426-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '火', '金', '大凶');
INSERT INTO `five_phase` VALUES ('a3b5c398-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '火', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a3b72e91-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '土', '木', '大凶');
INSERT INTO `five_phase` VALUES ('a3b882bf-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '土', '火', '中吉');
INSERT INTO `five_phase` VALUES ('a3ba0953-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '土', '土', '中吉');
INSERT INTO `five_phase` VALUES ('a3bb68e2-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '土', '金', '中吉');
INSERT INTO `five_phase` VALUES ('a3bc5347-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '土', '水', '大凶');
INSERT INTO `five_phase` VALUES ('a3bdd9fb-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '金', '木', '凶多于吉');
INSERT INTO `five_phase` VALUES ('a3bfae97-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '金', '火', '凶多于');
INSERT INTO `five_phase` VALUES ('a3c10e2c-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '金', '土', '大吉');
INSERT INTO `five_phase` VALUES ('a3c1f88d-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '金', '金', '中吉');
INSERT INTO `five_phase` VALUES ('a3c37f41-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '金', '水', '大吉');
INSERT INTO `five_phase` VALUES ('a3c57b0d-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '水', '木', '大吉');
INSERT INTO `five_phase` VALUES ('a3c66564-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '水', '火', '大凶');
INSERT INTO `five_phase` VALUES ('a3c776cb-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '水', '土', '大凶');
INSERT INTO `five_phase` VALUES ('a3c8130c-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '水', '金', '大吉');
INSERT INTO `five_phase` VALUES ('a3c94b91-f1bf-11e7-bedc-ce15fda84d48', '2018-01-05 10:25:04', '2018-01-05 10:25:04', NULL, 1, '水', '水', '水', '中吉');

SET FOREIGN_KEY_CHECKS = 1;
