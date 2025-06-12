# Contributing to VTX1

Thank you for your interest in contributing to the VTX1 project! This document will help you get started, understand the project structure, and follow best practices for submitting your changes.

## Project Structure

- `src/` — Project source code (Icarus Verilog)
- `docs/` — Documentation sources (Asciidoc, diagrams)
- `release/` — Release artifacts (PDFs, generated diagrams)
- `Taskfile.yml` — Task automation for build, test, and documentation
- `CONTRIBUTE.MD` — This guide
- `README.MD` — Project overview

## Development Quick Start

1. **Install prerequisites:**

   ```bash
   scoop install task
   task install
   ```

2. **Build and test:**

   ```bash
   task test-all
   ```

3. **Generate documentation:**

   ```bash
   task doc
   ```

### Main Tools Used

- Icarus Verilog, vvp — Simulation
- gtkwave — Waveform viewing
- yosys — Synthesis
- mermaid-cli, wavedrom-cli, bytefield-svg — Diagrams
- Asciidoctor — Documentation

## IDE & Editor Support

- [Asciidoc IDE plugins](https://gitlab.eclipse.org/eclipse-wg/asciidoc-wg/asciidoc.org/-/blob/main/awesome-asciidoc.adoc#author)
- [Mermaid editor plugins](https://mermaid.js.org/ecosystem/integrations-community.html#editor-plugins)

## Conventions

### Code Style

- Use clear, descriptive names for modules, signals, and files.
- Keep code modular and well-commented.
- Follow Verilog best practices for simulation and synthesis.
- Place testbenches in `src/testbench/`.

### Testing

- Add or update testbenches for new features or bug fixes.
- Run `task test-all` before submitting changes.
- Use `gtkwave` to inspect simulation results as needed.

### Documentation

- Update or add Asciidoc files in `docs/` for new features or changes.
- Use diagrams (Mermaid, Wavedrom, Bytefield) where helpful.
- Regenerate documentation with `task doc`.

## Submitting Your Changes

1. **Fork the repository** and create a new branch for your feature or fix.
2. **Make your changes** with clear, atomic commits.
3. **Test your changes** locally using the provided tasks.
4. **Update documentation** if your change affects usage or design.
5. **Open a Pull Request (PR):**
   - Describe your changes and link related issues if any.
   - Ensure all checks pass (CI, tests, lint, etc.).
6. **Participate in code review** and address feedback promptly.

## Code of Conduct

Please be respectful and constructive in all project communications.

## Useful Links

- [Icarus Verilog](http://iverilog.icarus.com/)
- [Yosys](https://yosyshq.net/yosys/)
- [GTKWave](http://gtkwave.sourceforge.net/)
- [Asciidoctor](https://asciidoctor.org/)
- [Mermaid](https://mermaid.js.org/)
- [Wavedrom](https://wavedrom.com/)
