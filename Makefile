DB_URL=postgres://marouaneboufarouj:@localhost:5432/gator

migrate-up:
	cd sql/schema && goose postgres "$(DB_URL)" up

migrate-down:
	cd sql/schema && goose postgres "$(DB_URL)" down

migrate-down-all:
	cd sql/schema && goose postgres "$(DB_URL)" down-to 0

migrate-reset:
	cd sql/schema && goose postgres "$(DB_URL)" down-to 0 && goose postgres "$(DB_URL)" up
