# Usage

The StateTracker is a tool used to record gas usage and state roots of an implementation when running chain-validation tests.

## Record
When set to record the state tracker will produce a file containing the results of state transitions for each test run during chain-validation. The file will look similar to:

```json
{"Receipt":{"ExitCode":0,"ReturnValue":"gkIAalUCaxVs74fCAZ1lL5JngwY2tNJPrOE=","GasUsed":2037},"Penalty":"0","Reward":"2037","Root":"bafy2bzaced3zohyakbqpbiaomxi67ymmxy7y6apwy7q2y44g6oci6c54ly75m"}
{"Receipt":{"ExitCode":0,"ReturnValue":"AA==","GasUsed":895},"Penalty":"0","Reward":"895","Root":"bafy2bzacebk6xuroukj2tcbapzyqrmjmhzsfolivoevbga6qhinb5x7ouesmy"}
{"Receipt":{"ExitCode":18,"ReturnValue":"","GasUsed":222},"Penalty":"0","Reward":"222","Root":"bafy2bzacebruj3eqg7bnmmpi2mx4zjy2ix5rb6k45vwywwq2gdzbh4juhtgeu"}
{"Receipt":{"ExitCode":18,"ReturnValue":"","GasUsed":168},"Penalty":"0","Reward":"168","Root":"bafy2bzacecu2hxgipuzhet33tbbra5ysluiuxcosjqnyjp5fmsnxm2ydekfew"}
{"Receipt":{"ExitCode":0,"ReturnValue":"","GasUsed":1027},"Penalty":"0","Reward":"1027","Root":"bafy2bzaceavryrwsjscfg3wcppogpbfqhp7e3yrjwohccxsepolhyczmotagi"}
...
...
```
Each line corresponds to an ApplyMessage or ApplyTipSetMessages call. If a test applies 4 messages its file is expected to have 4 lines, 5 messages 5 lines, etc.

### How To Record
1. Set the environment variable `CHAIN_VALIDATION_DATA` to the location of the chain-validation gas resources directory. For most users this will be: `$GOPATH/chain-validation/box/resources`.
2. Uncomment this [line](https://github.com/filecoin-project/chain-validation/blob/1f44d3090c52a1c443a2ca85c5747f3417197008/drivers/test.go#L281) to enable the statetracker `Record()` method. This will cause the statetracker to produce a file for each test as the location `CHAIN_VALIDATION_DATA`.
3. Run tests you wish to record gas for, and verify files with names corresponding to the tests exist in `CHAIN_VALIDATION_DATA`
4. Run `make gen-gas` to generate `box/blob.go` -- blob.go contains gas data as a go file and is used to populate the [resource box storage](https://github.com/filecoin-project/chain-validation/blob/f6bc23143d179bcccc9c30bfd00242a3c3398432/box/box.go#L8). Since chain-validation is a library imported by implementations storing this data in a go file is necessary.

## Validation

When set to validate the statetracker will [look up the testing values](https://github.com/filecoin-project/chain-validation/blob/f6bc23143d179bcccc9c30bfd00242a3c3398432/box/box.go#L40) for each test. If values cannot be found a warning log is displayed in the test output.
When new tests are added the Record process described above will need to be followed to generate values for them.
