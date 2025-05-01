# Clicker

This project is a simple incremental game built using the Ebiten game library in Go. It demonstrates how to structure a game project with separate files for game logic, input handling, and assets.

## Game Overview

Clicker is a simple incremental game where players can earn money through manual work or by purchasing buildings that generate passive income. The game features a text-based interface and is controlled via keyboard inputs.

### Key Features:
- **Manual Work**: Earn money manually by selecting the "Manual Work" option.
- **Buildings**: Purchase and upgrade buildings to generate passive income.
- **Popup Messages**: Informative messages guide the player when actions cannot be performed.
- **Debug Mode**: Enable debug mode to display internal game state for testing and development.

## How to Play

1. **Navigate the menu**:
   - Use the arrow keys (`↑`, `↓`) or `W`/`S` to move the cursor.
2. **Perform actions**:
   - Press `Enter` or `Space` to select an option.
3. **Earn money**:
   - Select "Manual Work" to earn money manually.
4. **Purchase buildings**:
   - Use earned money to purchase buildings and generate passive income.
5. **Close popups**:
   - Press `Enter` to close popup messages.

For play the game online, visit: [https://kmdk.uk/clicker/](https://kmdk.uk/clicker/)

## Project Structure

```
├── main.go          # Entry point of the application
├── game             # Contains game logic and input handling
│   ├── game.go      # Main game logic
│   ├── input.go     # User input handling
│   ├── building.go  # Building-related logic
│   ├── popup.go     # Popup message handling
│   ├── config.go    # Game configuration
│   └── draw.go      # Rendering logic
├── assets           # Contains game assets
│   └── README.md    # Information about game assets
├── go.mod           # Go module configuration
└── README.md        # Project documentation
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
   go run main.go
   ```

## Assets

For information regarding the images and audio files used in the game, please refer to the `assets/README.md` file.

## Debug Mode

To enable debug mode, modify the `EnableDebug` field in the `Config` struct:
```go
config := &Config{
    EnableDebug: true, // Enable debug mode
}
```
When enabled, debug information such as the player's money and cursor position will be displayed on the screen.

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
