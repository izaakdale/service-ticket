run:
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
