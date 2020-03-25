# Usage

The gasmeter is a tool used to record gas usage of an implementation when running chain-validation tests.

## Record
When set to record the gasmeter will produce a file containing gas values for each test run during chain-validation. The file will look similar to:

```text
GasUnitsUsed
GasUnitsUsed
GasUnitsUsed
...
...
```
Each `GasUnitsUsed` line corresponds to an ApplyMessage call. If a test applies 4 messages its gas file is expected to have 4 lines, 5 messages 5 lines, etc.

### How To Record
1. Set the environment variable `CHAIN_VALIDATION_DATA` to the location of the chain-validation gas resources directory. For most users this will be: `$GOPATH/chain-validation/box/resources`.
2. Uncomment this [line](https://github.com/filecoin-project/chain-validation/blob/1f44d3090c52a1c443a2ca85c5747f3417197008/drivers/test.go#L281) to enable the gasmeters `Record()` method. This will cause the gasmeter to produce a gas file for each test as the location `CHAIN_VALIDATION_DATA`.
3. Run tests you wish to record gas for, and verify files with names corresponding to the tests exist in `CHAIN_VALIDATION_DATA`
4. Run `make gen-gas` to generate `box/blob.go` -- blob.go contains gas data as a go file and is used to populate the [resource box storage](https://github.com/filecoin-project/chain-validation/blob/f6bc23143d179bcccc9c30bfd00242a3c3398432/box/box.go#L8). Since chain-validation is a library imported by implementations storing this data in a go file is necessary.

## Validation

When set to validate the gasmeter will [look up the gas values](https://github.com/filecoin-project/chain-validation/blob/f6bc23143d179bcccc9c30bfd00242a3c3398432/box/box.go#L40) for each test. If gas values cannot be found a warning log is displayed in the test output.
When new tests are added the Record process described above will need to be followed to generate gas values for them.
