BINARY_NAME = lynway-print-service
SERVICE_NAME = lynway-print-service.service

.PHONY: build install start stop restart status log

build:
	@echo "Building ($(BINARY_NAME))..."
	go build -o $(BINARY_NAME) cmd/lynway-print-service/main.go

# Запуск скрипта установки
install:
	@echo "Запуск скрипта установки..."
	chmod +x install.sh
	sudo ./install.sh

# Запуск службы в systemd
start:
	@echo "Запуск службы $(SERVICE_NAME)..."
	sudo systemctl start $(SERVICE_NAME)

# Остановка службы
stop:
	@echo "Остановка службы $(SERVICE_NAME)..."
	sudo systemctl stop $(SERVICE_NAME)

# Перезапуск службы (например, после обновления бинарника)
restart:
	@echo "Перезапуск службы $(SERVICE_NAME)..."
	sudo systemctl restart $(SERVICE_NAME)

# Полезный бонус: просмотр текущего статуса службы
status:
	@sudo systemctl status $(SERVICE_NAME)

log:
	journalctl -xeu $(SERVICE_NAME)