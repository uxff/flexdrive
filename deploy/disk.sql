

CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '会员id',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '会员姓名',
  `email` varchar(32) NOT NULL COMMENT '邮箱',
  `phone` varchar(12) NOT NULL DEFAULT '' COMMENT '手机号 ',
  `pwd` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `levelId` int(11) NOT NULL DEFAULT '0' COMMENT '级别id',
  `totalCharge` int(11) NOT NULL DEFAULT '0' COMMENT '累计充值 单位分',
  `quotaSpace` bigint(20) NOT NULL DEFAULT '0' COMMENT '当前拥有的空间 单位KB',
  `usedSpace` bigint(20) NOT NULL DEFAULT '0' COMMENT '当前已用空间 单位KB',
  `fileCount` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件数量',
  `lastLoginAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后登录时间',
  `lastLoginIp` varchar(16) NOT NULL DEFAULT '' COMMENT '最后登录ip',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=正常启用 99=账户冻结 ',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_level` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '会员级别id',
  `name` varchar(32) NOT NULL COMMENT '会员级别名称',
  `quotaSpace` bigint(20) NOT NULL DEFAULT '0' COMMENT '会员级别的会员空间 单位KB',
  `price` int(11) NOT NULL DEFAULT '0' COMMENT '会员级别的价格 单位分',
  `isDefault` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否新会员默认等级 1=是',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=启用 99=删除',
  `primeCost` int(11) NOT NULL DEFAULT '0' COMMENT '原价 仅用于展示 单位分',
  `desc` varchar(256) NOT NULL DEFAULT '' COMMENT '介绍',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `manager` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '管理员id',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '管理员名称',
  `phone` varchar(12) NOT NULL DEFAULT '' COMMENT '管理员手机号',
  `email` varchar(32) NOT NULL DEFAULT '' COMMENT '管理员email',
  `pwd` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=正常 99=删除',
  `roleId` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
  `isSuper` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否是超管 1=超管',
  `lastLoginAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后登录时间',
  `lastLoginIp` varchar(16) NOT NULL DEFAULT '' COMMENT '最后登录ip',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色id',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '角色名称',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=正常启用 99=删除',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `permit` text NOT NULL COMMENT '授权内容 json',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `oper_log` (
  `id` int(11) NOT NULL COMMENT '操作记录id',
  `managerId` int(11) NOT NULL COMMENT '操作员id 对应员工id',
  `managerName` varchar(32) NOT NULL DEFAULT '' COMMENT '操作人名称',
  `operBiz` varchar(64) NOT NULL COMMENT '操作业务 枚举',
  `operParams` text NOT NULL COMMENT '操作内容参数',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '操作时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=正常',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `order` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单id',
  `userId` int(11) NOT NULL COMMENT '会员id',
  `originLevelId` int(11) NOT NULL COMMENT '会员原等级id',
  `awardLevelId` int(11) NOT NULL DEFAULT '0' COMMENT '会员购买的等级id',
  `awardSpace` bigint(11) NOT NULL DEFAULT '0' COMMENT '本次购买的容量空间 单位KB',
  `phone` varchar(12) NOT NULL DEFAULT '' COMMENT '会员手机号',
  `levelName` varchar(12) NOT NULL DEFAULT '' COMMENT '等级名',
  `totalAmount` int(11) NOT NULL DEFAULT '0' COMMENT '订单价格 单位分',
  `payAmount` int(11) NOT NULL DEFAULT '0' COMMENT '实付款金额 单位分',
  `outOrderNo` varchar(40) NOT NULL DEFAULT '' COMMENT '第三方支付通道订单号',
  `remark` text NOT NULL COMMENT '订单备注',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=待付款 2=未付款关闭 3=已付款 4=退款中 5=已退款',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `node` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `nodeName` varchar(40) NOT NULL COMMENT '节点名',
  `nodeAddr` varchar(32) NOT NULL COMMENT '节点地址 集群中服务地址',
  `totalSpace` bigint(20) NOT NULL DEFAULT '0' COMMENT '全部空间 单位KB',
  `usedSpace` bigint(20) NOT NULL DEFAULT '0' COMMENT '使用的空间 单位KB',
  `unusedSpace` bigint(20) NOT NULL DEFAULT '0' COMMENT '未使用空间',
  `fileCount` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件数量',
  `remark` text NOT NULL COMMENT '房间备注',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '添加时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=正常启用 2=注册超时 99=删除 ',
  `lastRegistered` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后注册时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- todo 分库分表 by hash
CREATE TABLE `file_index` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '文件索引id',
  `fileName` varchar(64) NOT NULL DEFAULT '' COMMENT '文件名 无用',
  `fileHash` varchar(40) NOT NULL DEFAULT '' COMMENT '文件内容哈希',
  `nodeId` int(11) NOT NULL DEFAULT '0' COMMENT '所在节点名 第一副本所在节点',
  `nodeId2` int(11) NOT NULL DEFAULT '0' COMMENT '所在节点名 第二副本所在节点',
  `nodeId3` int(11) NOT NULL DEFAULT '0' COMMENT '所在节点名 第三副本所在节点',
  `innerPath` varchar(256) NOT NULL DEFAULT '' COMMENT '文件在服务器路径',
  `outerPath` varchar(256) NOT NULL DEFAULT '' COMMENT '文件外部访问路径',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 0=未就绪 1=就绪 98=上传失败 99=删除',
  `referCount` int(11) NOT NULL DEFAULT '0' COMMENT '被引用数量',
  `size` bigint(20) NOT NULL DEFAULT '0' COMMENT '大小 单位Byte',
  `space` bigint(20) NOT NULL DEFAULT '0' COMMENT '占用空间单位 单位KB',
  `desc` text NOT NULL COMMENT '描述信息',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- todo 分库分表 by uid
CREATE TABLE `user_file` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '文件id',
  `fileIndexId` int(10) NOT NULL DEFAULT '0' COMMENT '文件索引id',
  `userId` int(11) NOT NULL DEFAULT '0' COMMENT '会员id',
  `filePath` varchar(256) NOT NULL DEFAULT '' COMMENT '文件父路径 /结尾',
  `fileName` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `pathHash` varchar(40) NOT NULL DEFAULT '' COMMENT '路径哈希，hash(filePath)',
  `fileHash` varchar(40) NOT NULL DEFAULT '' COMMENT '文件内容哈希',
  `nodeId` int(11) NOT NULL DEFAULT '0' COMMENT '所在节点名 第一副本所在节点',
  `isDir` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否是目录',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=正常 2=隐藏 99=下架',
  `size` bigint(20) NOT NULL DEFAULT '0' COMMENT '大小 单位Byte 目录则记录0',
  `space` bigint(20) NOT NULL DEFAULT '0' COMMENT '占用空间单位 单位KB 目录则记录0',
  `desc` text NOT NULL COMMENT '描述信息',
  PRIMARY KEY (`id`),
  INDEX IDX_User_PathHash (`userId`, `pathHash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `share` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '文件id',
  `fileHash` varchar(40) NOT NULL DEFAULT '' COMMENT '文件内容哈希',
  `shareHash`  varchar(40) NOT NULL DEFAULT '' COMMENT 'share哈希 用于路由',
  `userId` int(11) NOT NULL DEFAULT '0' COMMENT '分享者会员id',
  `userFileId` int(11) NOT NULL DEFAULT '0' COMMENT '分享者会员文件索引id',
  `nodeId` int(11) NOT NULL DEFAULT '0' COMMENT '所在节点名',
  `fileName` varchar(64) NOT NULL DEFAULT '' COMMENT '文件名',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=正常 2=隐藏 99=已删除',
  `expired` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '分享有效期',
  PRIMARY KEY (`id`),
  INDEX IDX_ShareHash (`shareHash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `offline_task` (
  `id` int(11) NOT NULL COMMENT '任务id',
  `userId` int(11) NOT NULL COMMENT '会员id',
  `dataurl` varchar(256) NOT NULL DEFAULT '' COMMENT '资源地址',
  `fileName` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1=下载中 2=下载完成 3=下载失败 4=已保存',
  `parentUserFileId` int(11) NOT NULL DEFAULT '0' COMMENT '父目录id 对应userFile.Id',
  `userFileId` int(11) NOT NULL DEFAULT '0' COMMENT '会员文件id',
  `fileHash` varchar(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小',
  `remark` varchar(256) NOT NULL DEFAULT '' COMMENT '备注 比如失败原因',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='离线任务表';

