# This repo is deprecated 🔚

See https://github.com/filecoin-project/test-vectors/ instead.

---

# Chain-Validation
[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![CircleCI](https://circleci.com/gh/filecoin-project/chain-validation.svg?style=svg)](https://circleci.com/gh/filecoin-project/chain-validation)

This library provides tools for validating the correctness of a Filecoin implementation according to the [specification](https://github.com/filecoin-project/specs). 

To maintain consensus, all Filecoin implementations must produce identical state transformations for any (state, message) pair. Further, they must implement the same block reward and chain selection logic. Validating correctness in this respect requires extensive coverage over (state, message) pairs, message sequences, and blockchain structures, and is important in maintaining the security and integrity of the network.

This library designed to allow any implementation of Filecoin to import it, implement a simple “driver” interface, and then run the tests provided by the testing library, passing the driver in as the parameter. 

For a comprehensive project description refer to the [Filecoin Chain-Validation Tools Design Doc](https://docs.google.com/document/d/1o0ODvpKdWsYMK_KmK-j-uPxYei6CZAZ4n_3ilQJPn4A/edit#).

## Goals
- A validation library that is implementation-independent enabling validation suites to be written once and used by different Filecoin implementations.
- High-level script-like methods for constructing long and complex message sequences, and making semantic assertions about the expected state resulting from their application.
- High-level script-like methods for constructing complex blockchain structures containing those messages, and making assertions about the expected state from their evaluation.
- Validation suites with significant coverage over actor state and code paths.
- Integration with both Go-filecoin and Lotus, enabling importing and use of the validation suites.
- Incremental utility to both these implementations while they are in development (rather than requiring an implementation to be complete before validation is useful)

## Non-Goals
- Immediate integration with Filecoin implementations not written in Go (though there should be a path towards this). Other implementations will be expected to write code for their implementation to work with this tool.
- High-performance execution, if this comes at a cost of timeliness or comprehensiveness

## Usage
// TODO
