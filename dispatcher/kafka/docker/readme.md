# kafka manage ui
http://localhost:8080

```bash
docker-compose up -d
docker-compose logs -f
```

```bash
# 进入 kafka 容器
docker exec -it <kafka-container-id> bash

# 创建主题
kafka-topics.sh --create --topic test-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

# 发送消息
kafka-console-producer.sh --topic test-topic --bootstrap-server localhost:9092

# 接收消息（新开终端）
kafka-console-consumer.sh --topic test-topic --from-beginning --bootstrap-server localhost:9092
```

```bash
docker-compose down
```