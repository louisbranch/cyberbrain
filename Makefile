.PHONY: db

db:
	docker swarm init
	docker stack deploy -c stack.yml postgres
