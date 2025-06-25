# VTX1 Ternary SoC

A personal project for a hardware-optimized, balanced ternary logic-based System-on-Chip (SoC) designed for both FPGA prototyping and silicon fabrication. VTX1 explores the possibilities of ternary logic and implements a complete code-to-silicon development chain.

---

## Table of Contents

- [Project Description](#project-description)
- [Key Features](#key-features)
- [Project Status](#project-status)
- [Implementation](#implementation)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
- [Build, Develop and Contribute](#build-develop-and-contribute)
- [Versioning](#versioning)
- [Changelog](#changelog)
- [License](#license)
- [Contact](#contact)

---

## Project Description

VTX1 leverages the computational advantages of balanced ternary logic (trits: -1, 0, +1, represented digitally using standard CMOS voltage levels) to achieve efficient computation, analog compatibility, and seamless integration with modern embedded systems.

### Included Components

- Verilog source files
- Assembler, compiler, and ISA description
- Microcode for the CPU
- Testbench and validation suite
- Documentation (asciidoc format)

## Key Features

- **Balanced ternary logic** internally, with binary interfaces externally
- **Optimized for mathematical operations** and signal processing
- **Vector processing capabilities** for enhanced performance
- **Low power consumption** via efficient ternary operations
- **FPGA-ready design** with silicon fabrication considerations
- **Complete toolchain** from RTL to documentation generation

## Project Status

**✅ Completed:**
- Core ternary arithmetic and logic modules
- CPU pipeline with microcode sequencer
- Memory controller and cache system
- I/O controllers (UART, SPI, I2C, GPIO, DMA)
- Complete build and test automation with Taskfile
- Comprehensive test suite with multiple validation scenarios
- Synthesis flow optimized for die production (JSON output)
- Documentation framework with PDF generation

**🚧 In Progress:**
- System integration testing and optimization
- I2C state machine fixes and validation
- Physical design preparation for silicon

**🔮 Planned:**
- OpenROAD/OpenLane integration for place & route
- FPGA prototype validation on development boards
- Silicon tape-out preparation with SkyWater 130nm PDK
- Advanced debugging and monitoring features

## Implementation

VTX1 is written in [Icarus Verilog](https://bleyer.org/icarus/), a free compiler for the IEEE-1364 Verilog hardware description language.

The project uses a [Taskfile](https://taskfile.dev/) build system for automation.

## Prerequisites

- [Icarus Verilog](https://bleyer.org/icarus/)
- [Taskfile](https://taskfile.dev/)
- Python 3.x (for scripts and tooling)
- Make (for some build steps)
- Git

## Quick Start

Clone this repository and build the project:

```sh
git clone https://github.com/yourusername/vtx1.git
cd vtx1
task build
```

To run tests:

```sh
task test
```

## Documentation

Comprehensive documentation is available in the `/docs` directory. The documentation is written in asciidoc format and can be generated in both HTML and PDF formats:

- [Architecture Specification](docs/Architecture.adoc) — Detailed architecture description
- [PDF version](release/vtx1-architecture.pdf) of the document

## Build, Develop and Contribute

Contributions are welcome! See [CONTRIBUTE.MD](CONTRIBUTE.MD) for guidelines.

## Versioning

This project uses [Semantic Versioning](https://semver.org/). Version tags are managed via Git and published on the [GitHub Releases](https://github.com/itworks99/vtx1/releases) page.

## Changelog

All notable changes are documented in [CHANGELOG.md](CHANGELOG.md).

## License

Copyright 2025 K. Vanyushov

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.

## Contact

For questions or feedback, please open an issue or contact [@itworks99](https://github.com/itworks99) on GitHub.
