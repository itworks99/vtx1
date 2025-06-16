# VTX1 Ternary SoC

A long-standing personal project for hardware-optimized, balanced ternary logic-based System-on-Chip (SoC) designed for both FPGA prototyping and silicon fabrication. It explores the possibilities of ternary logic and implements a complete code-to-silicon development chain.

## Project Description

VTX1 leverages the computational advantages of balanced ternary logic (trits: -1, 0, +1 represented digitally using standard CMOS voltage levels) to achieve efficient computation, analog compatibility, and seamless integration with modern embedded systems.

### Key Features

- **Balanced ternary logic** internally with binary interfaces externally
- **Optimized for mathematical operations** and signal processing
- **Vector processing capabilities** for enhanced performance
- **Low power consumption** through efficient ternary operations
- **FPGA-ready design** with silicon fabrication considerations
- **Complete toolchain** from RTL to documentation generation

### Current Status

✅ **Completed:**

- Core ternary arithmetic and logic modules
- CPU pipeline with microcode sequencer
- Memory controller and cache system
- I/O controllers (UART, SPI, I2C, GPIO, DMA)
- Complete build and test automation with Taskfile
- Comprehensive test suite with multiple validation scenarios
- Synthesis flow optimized for die production (JSON output)
- Documentation framework with PDF generation

🚧 **In Progress:**

- System integration testing and optimization
- I2C state machine fixes and validation
- Physical design preparation for silicon

🔮 **Planned:**

- OpenROAD/OpenLane integration for place & route
- FPGA prototype validation on development boards
- Silicon tape-out preparation with SkyWater 130nm PDK
- Advanced debugging and monitoring features

## Implementation

This project is written with [Icarus Verilog—](https://bleyer.org/icarus/)a free compiler implementation for the IEEE-1364 Verilog hardware description language.

It uses a [Taskfile](https://taskfile.dev/) build system used for a variety of project-related tasks.

## Documentation

Comprehensive documentation is available in the `/docs` directory. The project documentation is written in asciidoc format with aim to generate documentation both in html and pdf formats:

- [Architecture Specification](docs/Architecture.adoc) - Detailed architecture description
- Complete [PDF version](release/vtx1-architecture.pdf) of the document

## Install

Clone this repository:

```git
git clone https://github.com/yourusername/vtx1.git
```

## Build, Develop and Contribute

This project welcomes contributions! See [CONTRIBUTE.MD](CONTRIBUTE.MD) for development and contribution guidelines.

## License

Copyright 2025 K.Vanyushov

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
