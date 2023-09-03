DB_URL=postgresql://admin:password123@localhost:6500/jwt-go?sslmode=disable
HOST=http://localhost:8080
dev: 
	docker-compose up -d
dev-down:
	docker-compose down
server:
	air
text:
	go run main.go -server text
http:
	go run main.go -server http
video:
	go run main.go -server video
voice:
	go run main.go -server voice
install:
	go get github.com/go-playground/validator/v10
	go get -u github.com/gofiber/fiber/v2
	go get -u github.com/golang-jwt/jwt/v4
	go get github.com/redis/go-redis/v9
	go get github.com/satori/go.uuid
	go get github.com/spf13/viper
	go get gorm.io/driver/postgres
	go get -u gorm.io/gorm
postgres:
	docker run --name auth --network chat_network -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=password123 -p 6500:5432 -d postgres:latest
createdb: 
	docker exec -it auth createdb --username admin --owner=admin auth_db
dropdb:
	docker exec -it auth dropdb auth_db --username admin
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
sqlc:
	sqlc generate

.phony:
	dev dev-down server install text http video voice postgres createdb dropdb migrateup migratedown sqlc