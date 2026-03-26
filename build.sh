#!/bin/bash

echo "--- Building Backend Image ---"
docker build -t wasaphoto-backend:latest -f Dockerfile.backend .

echo "--- Building Frontend Image ---"
docker build -t wasaphoto-frontend:latest -f Dockerfile.frontend .

echo "--- Build completata! ---"
