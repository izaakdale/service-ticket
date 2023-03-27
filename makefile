run:
	PORT=8082 \
	HOST=localhost \
	QUEUE_URL=https://sqs.eu-west-2.amazonaws.com/735542962543/order-stored-queue \
	AWS_REGION=eu-west-2 \
	GRPC_HOST=localhost \
	GRPC_PORT=50002 \
	MAIL_HOST=localhost \
	MAIL_PORT=1025 \
	go run .

run_local:
	PORT=8082 \
	HOST=localhost \
	QUEUE_URL=http://localhost:4566/000000000000/order-stored-queue \
	AWS_REGION=eu-west-2 \
	AWS_ENDPOINT=http://localhost:4566 \
	GRPC_HOST=localhost \
	GRPC_PORT=50002 \
	MAIL_HOST=localhost \
	MAIL_PORT=1025 \
	go run .

build:
	docker build --tag izaakdale/service-ticket .
