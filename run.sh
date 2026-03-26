#!/bin/bash

# Avvia il Backend in una nuova finestra
gnome-terminal --title="BACKEND LOGS" -- bash -c "docker run --rm -p 8080:8080 --name wasa-back wasaphoto-backend; exec bash"

# Avvia il Frontend in una nuova finestra
gnome-terminal --title="FRONTEND LOGS" -- bash -c "docker run --rm -p 3000:3000 --name wasa-front wasaphoto-frontend; exec bash"

echo "I container sono in esecuzione in finestre separate."
