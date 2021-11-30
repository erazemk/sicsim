# TODO

## General
- Add error reporting / handling
- Add unit tests

## Common
- Fix isType() functions to properly check for data types

## Device
- Handle writing to stdin, stdout and stderr

## Executor
- Throw errors for not implemented commands
- Implement types of executions based on nixbpe bits
- Implement device-based commands
- Remove unneeded type conversions

## Machine
- Don't export devices (implement functions for accessing them inside device.go)

## Memory
- Word() and Float() should return slices of bytes (fix other functions needing to write to specific bytes)
