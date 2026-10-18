# PicoRV32 (vendored)

`picorv32.v` is the PicoRV32 RV32I(M) CPU core by Claire Xenia Wolf,
vendored unmodified from https://github.com/YosysHQ/picorv32 (ISC license,
see the header of the file). It is the binary control core for the
cycle-level Tier A comparison: the same benchmark kernels run on the VTX1
RTL and on this core under stated memory models, and both testbenches
assert the computed results.
