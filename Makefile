docker:
	docker build --rm -t cyberbrain/server:latest .

push:
	docker push cyberbrain/server:latest
