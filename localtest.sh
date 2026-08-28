mkdir -p .local/data
NOTED_ADDR=127.0.0.1:8080 \
NOTED_DATA_DIR="$PWD/.local/data" \
NOTED_DB_FILE=notes.db \
go -C server run ./cmd/noted