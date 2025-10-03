#!/bin/bash

set -e

# === CONFIG ===
# Set these environment variables before running the script
# export DEPLOY_SERVER="user@your-server-ip"
# export DEPLOY_KEY_PATH="$HOME/.ssh/your_key"
SERVER="${DEPLOY_SERVER:-}"
KEY_PATH="${DEPLOY_KEY_PATH:-}"
BINARY_NAME="${BINARY_NAME:-jobseeker}"
REMOTE_DIR="${REMOTE_DIR:-/root/jobseeker}"
SYSTEMD_FILE="jobseeker.service"
SYSTEMD_FILE_PATH="deploy/$SYSTEMD_FILE"
LOCAL_ENV_FILE=".env.prod"
REMOTE_ENV_FILE="$REMOTE_DIR/.env"

# === VALIDATION ===
if [ -z "$SERVER" ] || [ -z "$KEY_PATH" ]; then
    echo "❌ Error: Missing required environment variables"
    echo "Please set:"
    echo "  export DEPLOY_SERVER='user@your-server-ip'"
    echo "  export DEPLOY_KEY_PATH='\$HOME/.ssh/your_key'"
    exit 1
fi

if [ ! -f "$KEY_PATH" ]; then
    echo "❌ Error: SSH key file not found at $KEY_PATH"
    exit 1
fi

# === BUILD ===
echo "🔨 Building Go binary..."
GOOS=linux GOARCH=amd64 go build -o "$BINARY_NAME" ./cmd

echo "❌ stop the prev systemd service on server..."
ssh -i $KEY_PATH $SERVER <<EOF 
systemctl stop $SYSTEMD_FILE
EOF

echo "❌ remove the prev binary on server..."
ssh -i $KEY_PATH $SERVER <<EOF 
rm -rf $REMOTE_DIR/*
EOF

# === CREATE REMOTE DIR ===
echo "📁 Ensuring remote dir exists on server..."
ssh -i "$KEY_PATH" "$SERVER" "mkdir -p $REMOTE_DIR"

# === UPLOAD BINARY ===
echo "📤 Uploading binary to $REMOTE_DIR"
scp -i "$KEY_PATH" "$BINARY_NAME" "$SERVER:$REMOTE_DIR/$BINARY_NAME"

# === UPLOAD .env.prod AND RENAME TO .env ===
echo "📤 Uploading .env.prod as .env"
scp -i "$KEY_PATH" "$LOCAL_ENV_FILE" "$SERVER:$REMOTE_ENV_FILE"

# === UPLOAD SYSTEMD SERVICE FILE ===
echo "📤 Uploading systemd service file"
scp -i "$KEY_PATH" "$SYSTEMD_FILE_PATH" "$SERVER:/etc/systemd/system/$SYSTEMD_FILE"

# === RESTART SERVICE ===
echo "♻️ Restarting systemd service on server..."
ssh -i "$KEY_PATH" "$SERVER" <<EOF
  chmod +x $REMOTE_DIR/$BINARY_NAME
  systemctl daemon-reload
  systemctl enable $SYSTEMD_FILE
  systemctl restart $SYSTEMD_FILE
  systemctl status $SYSTEMD_FILE --no-pager
EOF

echo "✅ Deployment completed successfully."


echo "Cleaning.."
rm  jobseeker
