.PHONY: run dev build clean

# Разработка: hot reload + Genkit Dev UI (http://localhost:4000)
# Просто запусти `air` в терминале
dev:
	air

# Простой запуск без hot reload и Dev UI
run:
	go run ./cmd/bot

# Build
build:
	go build -o ./tmp/bot ./cmd/bot

clean:
	rm -rf ./tmp
