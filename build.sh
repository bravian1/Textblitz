#!/bin/bash

# Set variables
APP_NAME="textindex"
GO_MAIN="main.go"

clean(){
    rm -f "$APP_NAME"
}

# Build the Go program
clean

if ! dpkg -l | grep -q poppler-utils; then
    echo "poppler-utils is not installed. Installing..."
    sudo apt update && sudo apt install -y poppler-utils
else
    echo "poppler-utils is already installed."
fi

echo "Building $APP_NAME..."
go build -o "$APP_NAME" "$GO_MAIN"

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful!"
else
    echo "Build failed!"
    exit 1
fi
