Feature: Player Explosion
  In order to provide feedback to the player
  As a game developer
  I want the player to explode when they die

  Scenario: Player is crushed by a falling rock
    Given the game is in a state where a rock can fall on the player
    When the physics engine updates
    Then the player should be replaced by an explosion
    And adjacent dirt and gems should be destroyed

  Scenario: Player collides with an enemy
    Given the game is in a state where an enemy can collide with the player
    When the enemy moves into the player's tile
    Then the player should be replaced by an explosion
    And adjacent dirt and gems should be destroyed