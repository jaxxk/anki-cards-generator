#!/bin/zsh

PROGRAM_NAME="poggers"
INSTALL_DIR="$HOME/bin"

# Step 1: Build the program
echo "Building the Go program..."
if ! go build -o "$PROGRAM_NAME"; then
    echo "Error: Failed to build the program. Ensure Go is installed and properly set up."
    exit 1
fi
echo "Build successful."

# Step 2: Create install directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR"
fi

# Step 3: Move the executable to the install directory (overwrite if exists)
echo "Installing the program to $INSTALL_DIR..."
if ! mv -f "$PROGRAM_NAME" "$INSTALL_DIR/"; then
    echo "Error: Failed to move the program to $INSTALL_DIR."
    exit 1
fi
echo "Installation complete: $PROGRAM_NAME successfully updated in $INSTALL_DIR."

# Step 4: Add install directory to PATH (if not already added)
echo "Ensuring $INSTALL_DIR is in PATH..."
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$HOME/.zshrc"
    echo "Added $INSTALL_DIR to PATH. Restart your terminal or source ~/.zshrc to apply changes."
else
    echo "$INSTALL_DIR is already in PATH."
fi

# Step 5: Set environment variables for logging if not already set
echo "Checking environment variables for logging..."
if [ "$LOG_LEVEL" = "INFO" ] && [ "$LOG_MODE" = "production" ]; then
    echo "Logging environment variables are already set."
else
    echo "Setting environment variables for production logging..."
    export LOG_LEVEL="INFO"
    export LOG_MODE="production"
    echo "export LOG_LEVEL=INFO" >> "$HOME/.zshrc"
    echo "export LOG_MODE=production" >> "$HOME/.zshrc"
    echo "Environment variables set successfully. Restart your terminal or source ~/.zshrc to apply changes."
fi

echo "All steps complete! You can now run $PROGRAM_NAME from anywhere."
