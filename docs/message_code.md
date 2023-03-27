# message code

### message code components
message code is a 6-digit number, use `ABCDEF` to present each digit
- `A`: the log level, 1-debug, 2-info, 3-warn, 4-error
- `BC`: the module number
- `D`: the submodule
- `EF`: the sequence number

### relations between code and module

| BC  | module  | D   | submodule |
|-----|---------|-----|-----------|
| 00  | message | 0   | general   |
| 01  | health  | 0   | health    |
| 02  | mysql   | 0   | config    |
| 02  | mysql   | 1   | service   |
| 03  | pmm     | 0   | config    |
