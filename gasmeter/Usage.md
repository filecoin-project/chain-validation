# Usage

The gasmeter package is used to record and validate an implementations gas usage.

When recording, gasmeter will produce a file for each test that is ran. The file will be named after the test it is
generate for with all character but alpha numeric striped. The file will contain a set of new line delimited comma
separated values of the form:
```text
oldStateRootCID,MessageCID,newStateRootCID,GasUnitsUsed
...
...
```

When validating, gasmeter will search for the file corresponding to the test it is running (ignoring non alpha numerica character)
and validate the gas used by an apply message call matches the value in the file. The order the messages are applied in 
correspond to the line number in the file. For example, the first call to apply message in a test will validate against
the first line in the file, sencond call to apply message to second line in file, and so on.

If gasmeter fails to find the file for the test it is running, the test fails.

## How To Record
1. Set the environment variable `CHAIN_VALIDATION_DATA` to the destination you wish gasmeter to write gas files to. Usually
this is `$GOPATH/chain-validation/gasmeter/gas_files`.
2. Set `RecordGas` in the implementations `ValidationConfig` to `true`.
3. Run tests you wish to record gas for

## How To Validate
1. Set the environment variable `CHAIN_VALIDATION_DATA` to the location you wish gasmeter to read gas files from. Usually
this is `$GOPATH/chain-validation/gasmeter/gas_files`.
2. Run the tests you wish to validate.

