框架源码
<pre>
├─cfg	--config.ini 读取, 写log
├─data	--程序运行数据
│  │  upinfo.tmp
│  │
│  └─userdata
│      ├─100
│      ├─110
│      └─50
├─helper	--一些方便的方法
│      export.go	--封装log文件
│      profilingtool.go
│      stack.go
│
├─misc
│  └─timer	--定时器
├─storage	--更新仓库文件档案（exe） 分析upinfo.tmp和更新数据库
│  │  config.ini	--配置文件
│  │  main.go	
│  │  make.sh	--linux编译文件
│  │  query.sh	--查询数据库脚本
│  │  run.sh	--linux运行文件
│  │  storage.sql	--建库文件
│  │  设计.md
│  │ 
│  └─core	--主要逻辑
└─util	
</pre>
