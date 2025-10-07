# Goulder Dash Development Plan

This document outlines the planned features and improvements for Goulder Dash. It serves as a roadmap for development, breaking down the high-level goals from the `README.md` into more concrete steps.

## 1. Core Gameplay Enhancements

### 1.1. Explosions & Crushing

#### Player Death Explosion (Simple)

- [x] **Crushing Logic**: Implement collision detection to determine when the player is crushed by a falling rock.
- [x] **Single-Tile Explosion Animation**: When the player is crushed, trigger a multi-frame explosion animation at their location only. The `game/types.go` file already defines `TileExplosion0` through `TileExplosion5`.
- [x] **No Area Effect**: Player death explosions only change the player tile itself, no damage to adjacent tiles.

#### Enemy Death Explosion (Area of Effect)

- [x] **Enemy Crushing Logic**: Implement collision detection to determine when an enemy is crushed by a falling rock/stone or gem/diamond.
- [x] **Multi-Tile Explosion Animation**: When an enemy is crushed, trigger explosion animation on:
  - [x] The enemy's tile (center of explosion)
  - [x] All 8 adjacent tiles (orthogonal and diagonal neighbors)
- [x] **Area of Effect Damage**: Enemy explosions should destroy elements in adjacent tiles:
  - [x] Destroy dirt tiles
  - [x] Destroy gems
  - [x] Destroy stones/rocks
  - [x] Potentially destroy other enemies in blast radius
  - [x] Leave walls/borders intact
  - [x] All affected tiles show explosion animation simultaneously
- [x] **Diamond Conversion**: When a diamond/gem falls onto an enemy, the enemy explodes and all 8 adjacent tiles are converted to diamonds (instead of showing explosion damage)

### 1.2. Simple Enemy AI
- [ ] **Basic Movement Patterns**: Implement simple movement for enemies. For example, an enemy could patrol a horizontal or vertical path, reversing direction when hitting a wall.
  - [ ] Create enemy struct with position and direction
  - [ ] Implement patrol logic (horizontal/vertical)
  - [ ] Add wall collision detection and direction reversal
- [ ] **Player Interaction**: If the player touches an enemy, the player should be killed (triggering the crushing/explosion sequence).
  - [ ] Implement player-enemy collision detection
  - [ ] Trigger player death on contact
- [ ] **Environmental Interaction**: Enemies should be subject to the same physics as the player, meaning they can be crushed by falling rocks.
  - [ ] Apply gravity/falling object physics to enemies
  - [ ] Trigger enemy explosion when crushed

### 1.3. Improved Digging
- [ ] **Digging Animation**: Add a specific animation for the player when digging through dirt.
  - [ ] Create digging sprite frames
  - [ ] Implement animation state machine
- [ ] **Particle Effects**: Spawn small dirt particles when the player digs to enhance the visual feedback.
  - [ ] Create particle system
  - [ ] Add dirt particle sprites
  - [ ] Spawn particles on dig events

## 2. Audio System

### 2.1. Sound Effects

- [ ] Integrate a sound library to handle audio playback
  - [ ] Research and choose audio library (e.g., beep, oto)
  - [ ] Set up audio system in project
- [ ] **Gem Collection**: A "bling" sound when a gem is collected.
  - [ ] Create/acquire gem collection sound
  - [ ] Implement playback on gem pickup
- [ ] **Digging**: A soft digging or scraping sound.
  - [ ] Create/acquire digging sound
  - [ ] Implement playback when moving through dirt
- [ ] **Rock Movement**: Sounds for a rock starting to fall and for it landing.
  - [ ] Create/acquire rock falling sound
  - [ ] Create/acquire rock landing sound
  - [ ] Implement playback in physics system
- [ ] **Explosion**: A boom sound for explosions.
  - [ ] Create/acquire explosion sound
  - [ ] Implement playback during explosion animation
- [ ] **Player Death**: A specific sound for when the player dies.
  - [ ] Create/acquire death sound
  - [ ] Implement playback on player death
- [ ] **Level Start/End**: Jingles to signify the start and successful completion of a level.
  - [ ] Create/acquire level start jingle
  - [ ] Create/acquire level complete jingle
  - [ ] Implement playback on level transitions

### 2.2. Music

- [ ] **Background Music**: Add a looping, retro-style chiptune track that plays during gameplay.
  - [ ] Create/acquire background music track
  - [ ] Implement looping music playback system
  - [ ] Add music during gameplay state
- [ ] **Game Over Theme**: A short, distinct musical piece for the game over screen.
  - [ ] Create/acquire game over music
  - [ ] Implement playback on game over state

## 3. Level Progression

### 3.1. Multiple Levels

- [ ] **Level File Format**: Define a clear format for level files (e.g., simple text files or a structured format like JSON).
  - [ ] Design level file schema
  - [ ] Document level format specification
  - [ ] Create example level files
- [ ] **Level Loading System**: Implement a function to load level data from files by a level number or name.
  - [ ] Implement level parser
  - [ ] Add level validation logic
  - [ ] Integrate with existing level system
- [ ] **Level Progression Logic**: When the player enters an open exit, transition to the next level.
  - [ ] Detect exit collision when exit is open
  - [ ] Implement level transition trigger
  - [ ] Load next level in sequence

### 3.2. Level Transitions

- [ ] **Transition Screen**: Display a simple screen between levels (e.g., "Level 2: Get Ready!").
  - [ ] Create transition screen state
  - [ ] Design transition screen UI
  - [ ] Add timer/input to continue
- [ ] **State Management**: Ensure the game state (player position, remaining gems, etc.) is correctly reset or carried over between levels.
  - [ ] Reset player position to spawn point
  - [ ] Reset collected gems counter
  - [ ] Preserve score and lives across levels
  - [ ] Clear explosion/animation states

## 4. UI/UX Improvements

### 4.1. Main Menu

- [ ] **Menu State**: Create a new game state for the main menu.
  - [ ] Define menu game state enum/constant
  - [ ] Implement menu state initialization
  - [ ] Add state switching logic
- [ ] **Menu Options**: Implement a simple menu with options like "Start Game" and "Quit".
  - [ ] Design menu layout
  - [ ] Implement menu rendering
  - [ ] Add keyboard navigation (up/down/enter)
  - [ ] Implement "Start Game" action
  - [ ] Implement "Quit" action

### 4.2. In-Game HUD

- [ ] **Score Display**: Show the player's current score on the screen.
  - [ ] Implement score tracking system
  - [ ] Add score rendering to HUD
  - [ ] Define score rules (gems, level completion, etc.)
- [ ] **Lives Counter**: Display the number of lives the player has left.
  - [ ] Implement lives system
  - [ ] Add lives rendering to HUD
  - [ ] Implement respawn logic on death
  - [ ] Add game over when lives reach zero
- [ ] **Gem Counter**: Show the number of gems collected and the number required to open the exit.
  - [ ] Add gem counter rendering to HUD
  - [ ] Show format: "Gems: X/Y" or similar
  - [ ] Update display on gem collection

---

This plan will be implemented incrementally, starting with the core gameplay features.