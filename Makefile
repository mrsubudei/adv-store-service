name = forum
port = 8083
path=$(shell pwd)

start:  build run
build:
		@docker build -t $(name) .
		@docker image prune --filter label=stage=build -f
run:
		@docker run -p $(port):8083 --name $(name) -v $(path)/database:/app/database -d $(name) 
exec:
		@docker exec -ti $(name) sh
stop:
		@docker stop $(name)
		@docker rm $(name)
remove:
		@docker rmi $(name)
kill:   stop remove

test:
  		go test ./...