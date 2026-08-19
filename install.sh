#!/bin/sh

APP_NAME="lynway-print-service"

SRC_BIN="lynway-print-service"
SRC_SERVICE="deploy/systemd/lynway-print-service.service"

BIN_DIR="/usr/local/bin"
CONFIG_DIR="/etc/lwprint"

mkdir -p "$CONFIG_DIR"
cp .env "$CONFIG_DIR/.env"
chmod 600 "$CONFIG_DIR/.env"


cp "$SRC_BIN" "$BIN_DIR/$APP_NAME"
chmod +x "$BIN_DIR/$APP_NAME"


cp "$SRC_SERVICE" "/etc/systemd/system/$APP_NAME.service"

systemctl daemon-reload
systemctl enable "$APP_NAME.service"
systemctl start "$APP_NAME.service"

