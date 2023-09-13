#Genera la imagen del servidor central.
docker-central:
	docker build -f Dockerfile.central . -t containerized_central:latest
	docker run --rm --name central-server --network="host" containerized_central:latest

#Descarga la imagen de la cola rabbit de dockerhub y la inicia
docker-rabbit:
	docker run -d --hostname rmq --name rabbit-server -p 8080:15672 -p 5672:5672 --network="host" rabbitmq:3-management

#Genera todas las imagenes de los servidores regionales.
docker-regional:
	docker build -f Dockerfile.america . -t containerized_america:latest
	docker build -f Dockerfile.asia . -t containerized_asia:latest
	docker build -f Dockerfile.oceania . -t containerized_oceania:latest
	docker build -f Dockerfile.europa . -t containerized_europa:latest

	docker run --rm --name america-server -p 50051:50051 --network="host" containerized_america:latest
	docker run --rm --name asia-server -p 50052:50052 --network="host" containerized_asia:latest
	docker run --rm --name oceania-server -p 50054:50054 --network="host" containerized_oceania:latest
	docker run --rm --name europa-server -p 50056:50056 --network="host" containerized_europa:latest

#Generar servidores regionales individualmente
#dist077 (central + asia)
docker-regional-asia:
	docker build -f Dockerfile.asia . -t containerized_asia:latest
	docker run --rm --name asia-server -p 50052:50052 --network="host" containerized_asia:latest

#dist078 (rabbit + europa)
docker-regional-europa:
	docker build -f Dockerfile.europa . -t containerized_europa:latest
	docker run --rm --name europa-server -p 50056:50056 --network="host" containerized_europa:latest

#dist079 (oceania)
docker-regional-oceania:
	docker build -f Dockerfile.oceania . -t containerized_oceania:latest
	docker run --rm --name oceania-server -p 50054:50054 --network="host" containerized_oceania:latest

#dist080 (america)
docker-regional-america:
	docker build -f Dockerfile.america . -t containerized_america:latest
	docker run --rm --name america-server -p 50051:50051 --network="host" containerized_america:latest

	
