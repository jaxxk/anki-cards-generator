#!/bin/bash

# Constants
PROGRAM_NAME="poggers"  
INSTALL_DIR="/usr/local/bin"  # Directory to install the program

# Ensure script is run with sufficient permissions
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (e.g., 'sudo ./install.sh')"
  exit 1
fi

# Step 1: Build the program
echo "Building the Go program..."
if ! go build -o "$PROGRAM_NAME"; then
  echo "Error: Failed to build the program. Make sure you are in the program directory and have Go installed."
  exit 1
fi
echo "Build successful."

# Step 2: Move the executable to the install directory
echo "Installing the program to $INSTALL_DIR..."
if ! mv "$PROGRAM_NAME" "$INSTALL_DIR/"; then
  echo "Error: Failed to move the program to $INSTALL_DIR."
  exit 1
fi

# Step 3: Verify installation
if ! command -v "$PROGRAM_NAME" &> /dev/null; then
  echo "Error: Installation failed. Make sure $INSTALL_DIR is in your PATH."
  exit 1
fi

# Success message
echo "Installation complete! You can now run '$PROGRAM_NAME' from anywhere."

# Optional: Print version if available
echo "Verifying installation..."
"$PROGRAM_NAME" || echo "$PROGRAM_NAME installed, but no version info available."
