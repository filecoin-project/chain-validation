# Usage

The gasmeter package is used to record and validate an implementations gas usage.

When recording, gasmeter will produce a file for each test that is ran. The file will be named after the test it is
generate for with all character but alpha numeric striped. The file will contain gas data in the form:
```text
GasUnitsUsed
GasUnitsUsed
GasUnitsUsed
...
...
```

When validating, gasmeter will search for the gas vale corresponding to the test it is running (ignoring non alpha numerica character)
in box/box.go. The order the messages are applied in correspond to the line number in the file. For example, the first
call to apply message in a test will validate against the first line in the file, second call to apply message to 
second line in file, and so on.
The same pattern is used for TipSet application.

## How To Record
1. Set the environment variable `CHAIN_VALIDATION_DATA` to the location of the chain-validation gas resources file. For 
most users this will be: `$GOPATH/chain-validation/box/resources`.
2. Set `RecordGas` in the implementations `ValidationConfig` to `true`.
3. Run tests you wish to record gas for, and verify files with names corresponding to the tests previously ran exist in
`$GOPATH/chain-validation/box/resources`
4. Run `make gen-gas` to generate `box/blob.go`.


