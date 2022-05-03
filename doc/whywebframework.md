# Why you need a web framework

`net/http` 提供了基础Web功能：
- 监听端口
- 映射静态路由
- 解析HTTP报文     
-   
然而一些Web开发中的简单的需求并不支持，需要手工实现：
- 动态路由：例如`hello/:name` , `hello/*` 等规则
- 鉴权：没有分组/统一鉴权的能力，需要在每个路由映射的handler中实现
- 模板：没有统一简化的HTML机制
- ...
  