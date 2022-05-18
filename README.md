# searcher
简单搜索引擎


### 安装
首先安装mysql数据库并建立一个空库utf8mb4格式
然后在config.json中修改配置文件
#### 编译
```shell
go build
```
#### 导入数据
将悟空数据集的csv文件放入resources文件夹下然后执行下面命令
```shell
./searcher --init=true
```

#### 运行
```shell
./searcher
```

### 功能说明

API接口：   
GET http://localhost:8080/api/v1/search?s=搜索内容  