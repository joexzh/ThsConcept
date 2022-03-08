# ThsConcept

获取同花顺概念

需要`mongodb`, `mysql` (doc need update)

```bash
export MONGO_USER=example
export MONGO_PASSWORD=example
export MONGO_HOST_PORT=localhost:27017
export SERVER_PORT=8080
export mysql_user=example
export mysql_password=example
export mysql_host=localhost
export mysql_port=3306

# 从同花顺 q.10jqka.com.cn 获取概念 
ThsConcept -mode retrieve
# 或开启概念查询服务器
THsconcept -mode server
```