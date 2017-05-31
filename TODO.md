# Unity Multiplayer Soccer - TODO List

- Autoscaling (Down)
    - (up) if autoscaling up, check to see if there are cordoned nodes to turn back on (will need to
    ensure capacity calculation doesn't include unscheduled nodes)
    - once a cordoned node is empty, then delete the instance (after a given time)
    - Write minimum and maximum values for scaling
- Make it so a player can't go into the goal area
- Show "GOAL" when a goal happens
- Put a highlight on the goal you are aiming for, so you can tell which is which
- Deal with disconnections/reconnections
- Animate paddle on kick    