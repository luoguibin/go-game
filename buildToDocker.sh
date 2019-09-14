sudo docker build -t go-game:2.0.1 .
# sudo docker service update --image go-game:2.0.1 go-game-service

# docker images // 查看所有镜像
sudo docker rmi $(docker images -f "dangling=true" -q) // 删除构建镜像时暂时的临时的空名字镜像