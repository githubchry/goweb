# minio



```
docker启动minio, 默认账密adminminio/adminminio
docker run -p 9000:9000 --name minio-chry -d -v /mnt/e/temp/minio/data:/data minio/minio server /data









```

# 命令客户端

[MinIO Client Quickstart Guide](https://docs.min.io/docs/minio-client-quickstart-guide)

```
docker pull minio/mc

docker run minio/mc ls play


```

