

CREATE TABLE `user` (
  `id` integer PRIMARY KEY AUTOINCREMENT,
  `name` varchar(32) NOT NULL DEFAULT '' ,
  `email` varchar(32) NOT NULL ,
  `phone` varchar(12) NOT NULL DEFAULT '' ,
  `pwd` varchar(32) NOT NULL DEFAULT '' ,
  `levelId` integer NOT NULL DEFAULT '0' ,
  `totalCharge` integer NOT NULL DEFAULT '0' ,
  `quotaSpace` integer NOT NULL DEFAULT '0' , 
  `usedSpace` integer NOT NULL DEFAULT '0' , 
  `fileCount` integer NOT NULL DEFAULT '0' , 
  `lastLoginAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `lastLoginIp` varchar(16) NOT NULL DEFAULT '' , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `status` tinyint(4) NOT NULL DEFAULT '1' , 
  CONSTAINT `IDX_email` UNIQUE (`email`)
) ;

CREATE TABLE `user_level` (
  `id` integer PRIMARY KEY AUTOINCREMENT , 
  `name` varchar(32) NOT NULL , 
  `quotaSpace` integer NOT NULL DEFAULT '0' , 
  `price` integer NOT NULL DEFAULT '0' , 
  "isDefault" integer NOT NULL DEFAULT '0',
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `status` integer NOT NULL DEFAULT '1' 
);
 

CREATE TABLE `manager` (
  `id` integer PRIMARY KEY AUTOINCREMENT , 
  `name` varchar(32) NOT NULL DEFAULT '' , 
  `phone` varchar(12) NOT NULL DEFAULT '' , 
  `email` varchar(32) NOT NULL DEFAULT '' , 
  `pwd` varchar(32) NOT NULL DEFAULT '' , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `status` tinyint(4) NOT NULL DEFAULT '1' , 
  `roleId` integer NOT NULL DEFAULT '0' , 
  `isSuper` tinyint(4) NOT NULL DEFAULT '0' , 
  `lastLoginAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `lastLoginIp` varchar(16) NOT NULL DEFAULT '' ,
  CONSTRAINT `IDX_email` UNIQUE (email)
) ;

CREATE TABLE `role` (
  `id` integer PRIMARY KEY AUTOINCREMENT , 
  `name` varchar(32) NOT NULL DEFAULT '' , 
  `status` tinyint(4) NOT NULL DEFAULT '1' , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `permit` text NOT NULL
) ;

CREATE TABLE `oper_log` (
  `id` integer NOT NULL , 
  `managerId` integer NOT NULL , 
  `managerName` varchar(32) NOT NULL DEFAULT '' , 
  `operBiz` varchar(64) NOT NULL , 
  `operParams` text NOT NULL , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `status` tinyint(4) NOT NULL DEFAULT '1' 
) ;

CREATE TABLE `order` (
  `id` integer PRIMARY KEY AUTOINCREMENT , 
  `userId` integer NOT NULL , 
  `originLevelId` integer NOT NULL , 
  `awardLevelId` integer NOT NULL DEFAULT '0' , 
  `awardSpace` bigint(11) NOT NULL DEFAULT '0' , 
  `phone` varchar(12) NOT NULL DEFAULT '' , 
  `levelName` varchar(12) NOT NULL DEFAULT '' , 
  `totalAmount` integer NOT NULL DEFAULT '0' , 
  `payAmount` integer NOT NULL DEFAULT '0' , 
  `remark` text NOT NULL , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(4) NOT NULL DEFAULT '1' 
) ;

CREATE TABLE `node` (
  `id` integer PRIMARY KEY AUTOINCREMENT , 
  `nodeName` varchar(16) NOT NULL , 
  `nodeAddr` varchar(32) NOT NULL , 
  `totalSpace` integer NOT NULL DEFAULT '0' , 
  `usedSpace` integer NOT NULL DEFAULT '0' , 
  `fileCount` integer NOT NULL DEFAULT '0' , 
  `remark` text NOT NULL , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `status` tinyint(4) NOT NULL DEFAULT '1' , 
  `lastRegistered` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' 
) ;

--  todo 分库分表 by hash
CREATE TABLE `file_index` (
  `id` integer PRIMARY KEY AUTOINCREMENT , 
  `fileName` varchar(32) NOT NULL DEFAULT '' , 
  `fileHash` varchar(40) NOT NULL DEFAULT '' , 
  `nodeId` integer NOT NULL DEFAULT '0' , 
  `nodeId2` integer NOT NULL DEFAULT '0' , 
  `nodeId3` integer NOT NULL DEFAULT '0' , 
  `innerPath` varchar(256) NOT NULL DEFAULT '' , 
  `outerPath` varchar(256) NOT NULL DEFAULT '' , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `status` tinyint(4) NOT NULL DEFAULT '1' , 
  `referCount` integer NOT NULL DEFAULT '0' , 
  `size` integer NOT NULL DEFAULT '0' , 
  `space` integer NOT NULL DEFAULT '0' , 
  `desc` text NOT NULL 
) ;

-- todo 分库分表 by uid
CREATE TABLE `user_file` (
  `id` integer PRIMARY KEY AUTOINCREMENT , 
  `fileIndexId` integer NOT NULL DEFAULT '0' , 
  `userId` integer NOT NULL DEFAULT '0' , 
  `filePath` varchar(256) NOT NULL DEFAULT '' , 
  `fileName` varchar(256) NOT NULL DEFAULT '' , 
  `pathHash` varchar(40) NOT NULL DEFAULT '' , 
  `fileHash` varchar(40) NOT NULL DEFAULT '' , 
  `nodeId` integer NOT NULL DEFAULT '0' , 
  `isDir` tinyint(4) NOT NULL DEFAULT '0' , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `status` tinyint(4) NOT NULL DEFAULT '1' , 
  `size` integer NOT NULL DEFAULT '0' , 
  `space` integer NOT NULL DEFAULT '0' , 
  `desc` text NOT NULL 
) ;

CREATE TABLE `share` (
  `id` integer PRIMARY KEY AUTOINCREMENT , 
  `fileHash` varchar(40) NOT NULL DEFAULT '' , 
  `shareHash` varchar(40) NOT NULL DEFAULT '' , 
  `userId` integer NOT NULL DEFAULT '0' , 
  `userFileId` integer NOT NULL DEFAULT '0' , 
  `nodeId` integer NOT NULL DEFAULT '0' , 
  `fileName` varchar(32) NOT NULL DEFAULT '' , 
  `created` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' , 
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP , 
  `status` tinyint(4) NOT NULL DEFAULT '1' , 
  `expired` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' 
) ;

CREATE INDEX "IDX_shareHash"
ON "share" (
  "shareHash"
);

insert into manager (email,pwd,status,isSuper) values("admin@admin.com","e10adc3949ba59abbe56e057f20f883e",1,1);

