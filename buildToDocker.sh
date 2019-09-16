sudo docker build -t go-game:1.0.0 .
# sudo docker service update --image go-game:2.0.1 go-game-service
# docker service create --replicas 1 --name go-game-service image_name
# docker images // 查看所有镜像
sudo docker rmi $(docker images -f "dangling=true" -q) // 删除构建镜像时暂时的临时的空名字镜像