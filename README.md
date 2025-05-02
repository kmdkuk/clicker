# Clicker

This project is a simple incremental game built using the Ebiten game library in Go. It demonstrates how to structure a game project with separate files for game logic, input handling, and assets.

## Game Overview

Clicker is a simple incremental game where players can earn money through manual work or by purchasing buildings that generate passive income. The game features a text-based interface and is controlled via keyboard inputs.

### Key Features:
- **Manual Work**: Earn money manually by selecting the "Manual Work" option.
- **Buildings**: Purchase and upgrade buildings to generate passive income.
- **Upgrades**: Unlock and apply upgrades to enhance manual work or building efficiency.
- **Popup Messages**: Informative messages guide the player when actions cannot be performed.
- **Debug Mode**: Enable debug mode to display internal game state for testing and development.

## How to Play

1. **Navigate the menu**:
   - Use the arrow keys (`↑`, `↓`) or `W`/`S`/`J`/`K` to move the cursor.
2. **Switch pages**:
   - Use the left/right arrow keys or `A`/`D`/`H`/`L` to toggle between pages (buildings and upgrades).
3. **Perform actions**:
   - Press `Enter` or `Space` to select an option.
4. **Earn money**:
   - Select "Manual Work" to earn money manually.
5. **Purchase buildings**:
   - Use earned money to purchase buildings and generate passive income.
6. **Apply upgrades**:
   - Unlock upgrades to enhance your gameplay.
7. **Close popups**:
   - Press `Enter` to close popup messages.

For playing the game online, visit: [https://kmdk.uk/clicker/](https://kmdk.uk/clicker/)

## Project Structure

```
├── cmd/clicker       # Entry point of the application
│   └── main.go       # Main function to start the game
├── game              # Contains game core logic
│   └── game.go       # Main game logic
├── model             # Core data models
│   ├── building.go   # Building data model
│   ├── upgrade.go    # Upgrade data model
│   └── manual_work.go # Manual work data model
├── state             # Game state management
│   └── game_state.go # Handles game state updates
├── input             # Input handling logic
│   ├── handler.go    # Input handler implementation
│   └── decider.go    # Decision-making logic
├── ui                # User interface components
│   └── popup.go      # Popup message handling
├── config            # Configuration
│   └── config.go     # Game configuration
├── assets            # Contains game assets
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
