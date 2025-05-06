# Clicker

This project is a simple incremental game built using the [Ebitengine](https://ebitengine.org/) game library in Go. It demonstrates how to structure a game project with separate files for game logic, input handling, and assets.
You can play it at [https://kmdk.uk/clicker/](https://kmdk.uk/clicker/)

## Game Overview

Clicker is a simple incremental game where players can earn money through manual work or by purchasing buildings that generate passive income. The game features a text-based interface and is controlled via keyboard inputs.

### Key Features:
- **Manual Work**: Earn money manually by selecting the "Manual Work" option.
- **Buildings**: Purchase and upgrade buildings to generate passive income.
- **Upgrades**: Unlock and apply upgrades to enhance manual work or building efficiency.
- **Popup Messages**: Informative messages guide the player when actions cannot be performed.
- **Debug Mode**: Enable debug mode to display internal game state for testing and development.
- **Scrollable Lists**: Efficiently navigate long lists of buildings and upgrades.
- **Large Number Formatting**: Display large numbers in a readable format (e.g., 1K, 1M).

## How to Play

1. **Navigate the Menu**:
   - Use the arrow keys (`↑`, `↓`) or `W`/`S` to move the cursor.
2. **Switch Pages**:
   - Use the left/right arrow keys (`←`, `→`) or `A`/`D` to switch between the Buildings and Upgrades pages.
3. **Select an Option**:
   - Press `Enter` or `Space` to select an option.
4. **Earn Money**:
   - Select "Manual Work" to earn money manually.
5. **Purchase Buildings**:
   - Use earned money to purchase buildings for passive income.
6. **Apply Upgrades**:
   - Unlock upgrades to improve efficiency.
7. **Close Popups**:
   - Press `Enter` to close popup messages.

## Project Structure

```
├── cmd/clicker       # Entry point of the application
│   └── main.go       # Main function to start the game
├── game              # Contains game core logic
│   └── game.go       # Main game logic
├── application       # Application layer for use cases and DTOs
│   ├── dto/          # Data Transfer Objects
│   └── usecase/      # Use case implementations
├── domain/model      # Core data models
├── infrastructure    # Infrastructure layer for state and storage
│   ├── state/        # Game state management
│   └── storage/      # Save/load functionality
├── presentation      # Presentation layer for UI and input handling
│   ├── components/   # UI components (e.g., lists, tabs, popups)
│   ├── formatter/    # Number formatting utilities
│   └── input/        # Input handling
├── assets            # Game assets (fonts, images, etc.)
├── config            # Configuration files
├── Makefile          # Build and run commands
├── go.mod            # Go module configuration
├── .goreleaser.yaml  # Release configuration for GoReleaser
└── README.md         # Project documentation
```

## Prerequisites

Before running the game, ensure the following dependencies are installed on your system:

- **Go**: Version 1.24 or later
- **Required system libraries for Ebiten**:
  - On Ubuntu/Debian:
    ```bash
    sudo apt install libx11-dev libgl1-mesa-dev xorg-dev
    ```
  - On macOS:
    Ensure Xcode command-line tools are installed:
    ```bash
    xcode-select --install
    ```
  - On Windows:
    No additional libraries are required.

## Getting Started

1. **Clone the repository**:
   ```bash
   git clone https://github.com/kmdkuk/clicker.git
   cd clicker
   ```

2. **Run the game**:
   ```bash
   go run ./cmd/clicker/main.go
   ```

3. **Build the game**:
   ```bash
   go build -o bin/clicker ./cmd/clicker
   ```

4. **Build for WebAssembly**:
   ```bash
   make build-wasm
   ```

## Assets

For information regarding the images and audio files used in the game, please refer to the `assets/README.md` file.

## Debug Mode

To enable debug mode, use the `--debug` or `-d` flag:
```bash
go run ./cmd/clicker/main.go --debug
```

## Troubleshooting

### Common Issues

1. **Error: `fatal error: X11/Xlib.h: No such file or directory`**
   - This error occurs when the required system libraries for Ebiten are missing.
   - Solution:
     ```bash
     sudo apt install libx11-dev libgl1-mesa-dev xorg-dev
     ```

2. **Game does not start on macOS**
   - Ensure Xcode command-line tools are installed:
     ```bash
     xcode-select --install
     ```

3. **Performance issues on Windows**
   - Ensure your graphics drivers are up to date.

## Contributing

Contributions are welcome! If you'd like to contribute to this project, please fork the repository and submit a pull request.

## License

See the `LICENSE` file for details.
