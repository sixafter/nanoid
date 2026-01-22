# FIPS 140 Compliance

> **IMPORTANT:**  
> This package is *not* itself FIPS 140 validated or certified. See below for details on compatibility and limitations.

## Overview

This package is designed to operate within FIPS 140-2 and FIPS 140-3 validated environments, leveraging only the Go 
standard library cryptographic primitives. No custom or third-party cryptography is used.

## Compatibility with Go FIPS 140 Mode

Go provides a dedicated FIPS 140 mode, enabled with the environment variable `GODEBUG=fips140=on` or `GODEBUG=fips140=1`. 
When enabled, this mode restricts cryptographic operations to those explicitly permitted under the FIPS 140 standard, and applies additional runtime checks to enforce compliance.

This package is engineered for **full compatibility with Go’s FIPS 140 mode**:

- **All cryptographic operations** use only Go standard library algorithms.
- **No prohibited or non-standard cryptography** is invoked when FIPS mode is active.
- **No inclusion** of third-party or experimental crypto.

When Go’s FIPS 140 mode is active, any use of non-approved cryptography results in a runtime error, providing enforcement 
at the platform level.

## Important Limitations and Legal Notice

- **This package is not independently validated or certified under FIPS 140 by NIST or any authority.**  
  Operation within a FIPS-validated environment is supported, but this alone does not provide FIPS validation.
- **Go FIPS 140 mode does not grant FIPS 140 validation or certification.**  
  It only enforces approved algorithm use when run with a FIPS-validated Go runtime and platform.
- **Setting `GODEBUG=fips140=on` is enforcement, not certification.**  
  Your responsibility is to ensure the Go runtime, cryptographic libraries, and the OS are all FIPS-validated if you require formal compliance.
- **No warranty or guarantee of FIPS 140 compliance is expressed or implied.**  
  Regulatory compliance is the responsibility of the deploying organization.

## Reference: Go FIPS 140 Documentation

For further technical guidance and up-to-date status on FIPS 140 support in Go, consult the official documentation:

- [Go FIPS 140 Documentation](https://go.dev/doc/security/fips140)

## Recommendations for Deploying in FIPS-Regulated Environments

- **Deploy only on FIPS-validated platforms:**  
  Ensure the OS, Go runtime, and underlying crypto modules have valid FIPS 140 certifications.
- **Audit your deployment:**  
  Regularly review your supply chain, build, and deployment process for compliance.
- **Consult compliance and security professionals:**  
  Engage your organization’s legal and compliance teams to ensure your deployment meets all requirements.
- **Retain and review certifications:**  
  For compliance reporting, use only official documentation and certification records.

> **Disclaimer:**  
> Use of this package in FIPS-regulated environments is at your own risk. No warranty, express or implied, is provided 
> regarding FIPS 140 compliance. This package is distributed “as is” and must be reviewed and validated as part of your 
> organization’s own compliance and certification process.
