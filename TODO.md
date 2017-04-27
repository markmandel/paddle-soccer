# Unity Multiplayer Soccer - TODO List

- Move the ScoreController to being server side.
- Track scores once the go through the goal area (on screen)
- When score hits 7(?) game is over, and the server should shut down
- Autoscaling (Down)
    - If removing a node will still keep you above buffer - then cordon the node with the least pods on it
    - (up) if autoscaling up, check to see if there are cordoned nodes to turn back on
    - once a cordoned node is empty, then delete the instance
    - Write minimum and maximum values for scaling
- Show "GOAL" when a goal happens
- Make it so a player can't go into the goal area
- Put a red and blue highlight on each goal, so you can tell which is which
- Deal with disconnections/reconnections
- Animate paddle on kick