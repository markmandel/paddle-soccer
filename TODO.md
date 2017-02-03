# Unity Multiplayer Soccer - TODO List

- Implement matchmaker into Unity project (leave server option)
    - if we get an ip from /game/ (post) then connect
    - Otherwise, poll /game/{id} until we get a ip and port
- Cleanup, and write a stack of tests (especially the Go code)
- Make Kicking the ball a RPC call
- Move the ScoreController to being server side.
- Track scores once the go through the goal area (on screen)
- Show "GOAL" when a goal happens
- Make it so a player can't go into the goal area
- Put a red and blue highlight on each goal, so you can tell which is which
- Animate paddle on kick
- Does it make sense to shift creation to a ThirdPartyResource - could be more flexible. Allow people to set their own config vars, etc
- Specific Nodepools for game servers vs. everyone else.
- Autoscaling
- readiness check based on redis PING
- Go through TODOS in the code