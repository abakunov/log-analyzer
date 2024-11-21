#!/bin/bash

# Variables
APP_NAME="analyzer"
INSTALL_DIR="/usr/local/bin"

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Please install it to proceed."
    exit 1
fi

# Build the application
#!/bin/bash

echo "Building the application..."
if go build -o analyzer cmd/run/main.go; then
    echo "Build successful. The application is now installed as 'analyzer'."
else
    echo "Build failed. Make sure your project compiles successfully."
fi

# Move the executable to $INSTALL_DIR
echo "Installing $APP_NAME to $INSTALL_DIR..."
sudo mv $APP_NAME $INSTALL_DIR

# Verify that the application is installed
if command -v $APP_NAME &> /dev/null
then
    echo "The application has been successfully installed! Run it with the command '$APP_NAME'."
else
    echo "Installation failed. Check permissions for $INSTALL_DIR."
    exit 1
fi
