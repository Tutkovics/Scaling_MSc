# service-graph-simulator
Source code for docker image

[DockerHub link](https://hub.docker.com/repository/docker/tuti/service-graph-simulator "DH: tuti/service-graph-simulator")

## Run
```
# Until develop:
$ go run main.go -name Backend -delay 90 -port 9999 -cpu 90 -memory 900 -endpoint-url /read -endpoint-cpu 99 -endpoint-delay 192 -endpoint-url /index -endpoint-cpu 22 -endpoint-delay 111 -endpoint-call='"back-end:9898/staus__front-end:9876/health"' -endpoint-call "database:1234/asd?q=user1" 

# Docker: 
$ docker run -p 9090:9090 --rm my-onlab /app/main [CONFIG]
```

## Create and push image:
```
$ docker login --username=tuti

$ docker build -t tuti/service-graph-simulator:latest .
Successfully built <<ID>
Successfully tagged tuti/service-graph-simulator:latest

$ docker push tuti/service-graph-simulator 
```
