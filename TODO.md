# Unity Multiplayer Soccer - TODO List

- Test and see if it works on Minikube
- Autoscaling
- Make Kicking the ball a RPC call
- Move the ScoreController to being server side.
- Track scores once the go through the goal area (on screen)
- Show "GOAL" when a goal happens
- Make it so a player can't go into the goal area
- Put a red and blue highlight on each goal, so you can tell which is which
- Animate paddle on kick
- readiness check based on redis PING
- http health check
- Specific Nodepools for game servers vs. everyone else.