# logging

A simple logging tool

#### 介绍

1.  在每天0点分割前天的日志为单独的文件，并以日期+.log为文件名后缀,然后删除两个月以前的日志文件。
2.  每月1日0点30分，对上个月的日志文件以tar.gz进行打包压缩。

#### 安装教程

1.  go get github.com/Jerryhax/logging


#### 使用说明

1.  依赖 github.com/cobfig/cron/v3


#### 参与贡献

  欢迎提交issues
1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

