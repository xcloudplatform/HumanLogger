

## Logger

 - Runs in background, no UI
 - Switches to *Active* state when user activity happens
 - Switches to *Non Active* state when user activity doesn't happen for 3 sec
 - Takes a screenshot every half a second in *Active* state, every 10 seconds *Non Active* state
    - compares with last screenshot taken, discards duplicates
 - Saves screenshot as *timestamp-uuid.png* file
 - Writes CSV file *timestamp-uuid.csv*


## Robotizator

 - Provides gRPC interface to mouse/keyboard
