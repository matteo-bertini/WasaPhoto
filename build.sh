#!/bin/bash


echo "--- Building Frontend Image ---"
docker build -t wasaphoto-frontend:latest -f Dockerfile.frontend .

echo "--- Build completata! ---"
