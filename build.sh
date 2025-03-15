#!/bin/bash

# Set variables
APP_NAME="textindex"
GO_MAIN="main.go"

clean(){
    rm -f "$APP_NAME"
}

# Build the Go program
clean
echo "Building $APP_NAME..."
go build -o "$APP_NAME" "$GO_MAIN"

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful!"
else
    echo "Build failed!"
    exit 1
fi
