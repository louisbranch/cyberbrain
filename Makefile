.PHONY: db

swarm:
	docker swarm init

up:
	docker stack deploy -c stack.yml cyberbrain

down:
	docker stack rm cyberbrain
