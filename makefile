
dev:
	docker-compose up

killdev:
	@kill -9 $(shell top -l 1 | grep "air\s" | cut -d " " -f1)
