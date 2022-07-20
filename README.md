# dapr

#### 目标
框架目标是简化DDD开发难度，使开发人员可直接进行业务开发，不需关心技术细节与实现。实现技术与业务分离，提升开发效率与质量。


#### 介绍
本项目是对dapr二次开发，主要增加了DDD功能的支持。总体架构设计参考了Java Axon Framework框架。\

#### 部署
1.  下载本项目源代码到$GOPATH/src/github.com/dapr/dapr目录
2.  下载liuxd6825/components-contrib项目到$GOPATH/src/github.com/dapr/components-contrib目录
3.  按照INSTALL.md文档进行编译与安装。按文档可生成docker镜像，以及部署到k8s服务器上。

#### 功能
- 1.框架采用多租户模式设计，数据和方法中预留TenantId属性或参数。
- 2.框架采用采用接口、链式、函数式编程，可支持多种数据库扩展，目前仅支持MongoDB。
- 3.采用iris实现Http UserInterface层封装。
- 4.采用RSQL语言实现，前端复杂化查询。
- 5.对事件定义、事件注册、事件存储、事件发送、事件溯源等进行封装。
- 6.采用CQRS模式，可分别开发Command服务与Query服务。
- 7.对Repository进行了封装，采用接口与链式方式编程，可支持多种数据库，
- 8.优化了前端调用Command服务，异步交互的问题。降低开发复杂度。
- 9.在事件存储时增加聚合根关系维护功能和关系查询SDK,保证的事务，效率与易用性。
