# ThsConcept

获取同花顺概念

需要`mongodb`

```bash
export MONGO_USER=example
export MONGO_PASSWORD=example
export MONGO_HOST_PORT=192.168.0.1:27017
# 从同花顺 q.10jqka.com.cn 获取概念 
ThsConcept -mode retrieve
# 或开启概念查询服务器
THsconcept -mode server
```