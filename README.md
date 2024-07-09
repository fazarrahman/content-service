# content-service

docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=pass12345 \
  -v postgres_data:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:latest
