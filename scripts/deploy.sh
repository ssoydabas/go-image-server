#!/bin/bash

cd ~/apps/cannabox/image-server/ || { echo "Failed to change directory"; exit 1; }

# Clean untracked files and discard local changes
git clean -fd || { echo "Failed to clean untracked files"; exit 1; }
git reset --hard || { echo "Failed to reset local changes"; exit 1; }

# Fetch the latest changes and reset to the remote branch
git fetch origin main || { echo "Failed to fetch from origin"; exit 1; }
git reset --hard origin/main || { echo "Failed to reset to origin/main"; exit 1; }

# Ensure the deploy script has executable permissions for the next time regardless of any changes made to it
chmod +x ~/apps/cannabox/image-server/scripts/deploy.sh || { echo "Failed to set permissions on deploy.sh"; exit 1; }

# Create data directory if it doesn't exist and set permissions
mkdir -p ~/apps/cannabox/image-server/data || { echo "Failed to create data directory"; exit 1; }
mkdir -p ~/apps/cannabox/image-server/dev-data || { echo "Failed to create dev-data directory"; exit 1; }
chmod 777 ~/apps/cannabox/image-server/data || { echo "Failed to set permissions on data directory"; exit 1; }
chmod 777 ~/apps/cannabox/image-server/dev-data || { echo "Failed to set permissions on dev-data directory"; exit 1; }

docker compose down || { echo "Failed to stop containers"; exit 1; }

docker compose up -d --build || { echo "Docker Compose failed"; exit 1; }
