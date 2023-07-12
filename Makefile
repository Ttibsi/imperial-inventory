.PHONY: access
access:
	# This requires mysql-client installed on the host
	# docker exec -i II_DB mysql -h 127.0.0.1 --port=3306 --protocol=tcp -uuser -ppassword
	#Alternatively
	docker exec -it II_DB bash
	# Inside the docker container, run mysql -uuser -ppassword

.PHONY: clean
clean:
	docker stop II_DB
	docker rm II_DB

.PHONY: db
db:
	docker build -f Dockerfile -t env .
	docker run -ti --name=II_DB -p 3306:3306 env
