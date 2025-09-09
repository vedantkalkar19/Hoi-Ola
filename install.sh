#!/bin/bash

# hoi-ola installation script

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Please install Go 1.16 or higher."
    exit 1
fi

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"

# Build the application
echo "Building hoi-ola..."
cd "$SCRIPT_DIR"
go build -o hoi-ola

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Install to /usr/local/bin (requires sudo)
echo "Installing hoi-ola to /usr/local/bin..."
sudo cp hoi-ola /usr/local/bin/

# Check if installation was successful
if [ $? -ne 0 ]; then
    echo "Installation failed!"
    exit 1
fi

echo "hoi-ola has been successfully installed!"
echo "You can now run it from anywhere by typing 'hoi-ola'"