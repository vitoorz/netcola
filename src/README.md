#note

1. engine is a scheduler, dispatch msg from bus to each service's data message
2. services receive msg and process with it, then output it to output. In most cases, output is the bus
3. 'buffer' is used to pick up message from bus quickly and append the msg to a private array. The private array is for smoothing the speed of bus and service's process speed and for avoid bus full block.
4. services get msg actually from 'buffer's private array. 