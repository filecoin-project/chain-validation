# Chain-Validation
The intent of this library is to allow any implementation of Filecoin to import it, implement a simple “driver” interface, and then run the tests provided by the testing library, passing the driver in as the parameter. It is the responsibility of the testing library to define tests in accordance with the spec so that it may aid in verifying implementations.

For a comprehensive project description refer to the [Filecoin Chain-Validation Tools Design Doc](https://docs.google.com/document/d/1o0ODvpKdWsYMK_KmK-j-uPxYei6CZAZ4n_3ilQJPn4A/edit#)
## Goal
This project aims to provide tools and a high-coverage suite of validation vectors for multiple Filecoin implementations:
- A validation library that is implementation-independent enabling validation suites to be written once and used by different Filecoin implementations.
- High-level script-like methods for constructing long and complex message sequences, and making semantic assertions about the expected state resulting from their application.
- High-level script-like methods for constructing complex blockchain structures containing those messages, and making assertions about the expected state from their evaluation.
- Validation suites with significant coverage over actor state and code paths.
- Integration with both Go-filecoin and Lotus, enabling importing and use of the validation suites.
- Incremental utility to both these implementations while they are in development (rather than requiring an implementation to be complete before validation is useful)

## Non-Goals
- Immediate integration with Filecoin implementations not written in Go (though there should be a path towards this). Other implementations will be expected to write code for their implementation to work with this tool.
- High-performance execution, if this comes at a cost of timeliness or comprehensiveness
