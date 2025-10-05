# Goulder Dash Development Plan

This document outlines the planned features and improvements for Goulder Dash. It serves as a roadmap for development, breaking down the high-level goals from the `README.md` into more concrete steps.

## 1. Core Gameplay Enhancements

### 1.1. Explosions & Crushing
- **Crushing Logic**: Implement collision detection to determine when the player or an enemy is crushed by a falling rock.
- **Explosion Animation**: When a character is crushed, trigger a multi-frame explosion animation at their location. The `game/types.go` file already defines `TileExplosion0` through `TileExplosion5`.
- **Area of Effect Damage**: When an explosion occurs, it should destroy adjacent dirt, gems, and potentially other elements.

### 1.2. Simple Enemy AI
- **Basic Movement Patterns**: Implement simple movement for enemies. For example, an enemy could patrol a horizontal or vertical path, reversing direction when hitting a wall.
- **Player Interaction**: If the player touches an enemy, the player should be killed (triggering the crushing/explosion sequence).
- **Environmental Interaction**: Enemies should be subject to the same physics as the player, meaning they can be crushed by falling rocks.

### 1.3. Improved Digging
- **Digging Animation**: Add a specific animation for the player when digging through dirt.
- **Particle Effects**: Spawn small dirt particles when the player digs to enhance the visual feedback.

## 2. Audio System

### 2.1. Sound Effects
Integrate a sound library to handle audio playback.
- **Gem Collection**: A "bling" sound when a gem is collected.
- **Digging**: A soft digging or scraping sound.
- **Rock Movement**: Sounds for a rock starting to fall and for it landing.
- **Explosion**: A boom sound for explosions.
- **Player Death**: A specific sound for when the player dies.
- **Level Start/End**: Jingles to signify the start and successful completion of a level.

### 2.2. Music
- **Background Music**: Add a looping, retro-style chiptune track that plays during gameplay.
- **Game Over Theme**: A short, distinct musical piece for the game over screen.

## 3. Level Progression

### 3.1. Multiple Levels
- **Level File Format**: Define a clear format for level files (e.g., simple text files or a structured format like JSON).
- **Level Loading System**: Implement a function to load level data from files by a level number or name.
- **Level Progression Logic**: When the player enters an open exit, transition to the next level.

### 3.2. Level Transitions
- **Transition Screen**: Display a simple screen between levels (e.g., "Level 2: Get Ready!").
- **State Management**: Ensure the game state (player position, remaining gems, etc.) is correctly reset or carried over between levels.

## 4. UI/UX Improvements

### 4.1. Main Menu
- **Menu State**: Create a new game state for the main menu.
- **Menu Options**: Implement a simple menu with options like "Start Game" and "Quit".

### 4.2. In-Game HUD
- **Score Display**: Show the player's current score on the screen.
- **Lives Counter**: Display the number of lives the player has left.
- **Gem Counter**: Show the number of gems collected and the number required to open the exit.

---

This plan will be implemented incrementally, starting with the core gameplay features.