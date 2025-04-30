# Clicker

This project is a simple text-based game built using the Ebiten game library in Go. It demonstrates how to structure a game project with separate files for game logic, input handling, and assets.

## Project Structure

```
├── main.go          # Entry point of the application
├── game             # Contains game logic and input handling
│   ├── game.go      # Main game logic
│   └── input.go     # User input handling
├── assets           # Contains game assets
│   └── README.md    # Information about game assets
├── go.mod           # Go module configuration
└── README.md        # Project documentation
```

## Prerequisites

Before running the game, ensure the following dependencies are installed on your system:

- Go (version 1.24 or later)
- Required system libraries for Ebiten:
  - On Ubuntu/Debian: `sudo apt install libx11-dev libgl1-mesa-dev xorg-dev`

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/kmdkuk/clicker.git
   cd clicker
   ```

2. Run the game:
   ```bash
   go run main.go
   ```

## Game Overview

The game features a simple text-based interface where players can interact with the game world through keyboard and mouse inputs. The main game loop handles updating the game state and rendering the graphics.

## How to Play

Instructions on how to play the game will be provided here once the game mechanics are fully implemented.

## Assets

For information regarding the images and audio files used in the game, please refer to the `assets/README.md` file.

## Troubleshooting

If you encounter the following error:
```
fatal error: X11/Xlib.h: No such file or directory
```
Install the required system libraries:
```bash
sudo apt install libx11-dev libgl1-mesa-dev xorg-dev
```
