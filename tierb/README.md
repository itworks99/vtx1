# Tier B: Ternary CMOS Cells on SkyWater Sky130

The circuit-level evaluation the project has pointed at since the
simulation methodology was adopted: can ternary logic in **real CMOS**
(Sky130, 130nm, stock devices, typical corner) beat the measured cost of
emulating trits on a binary substrate — the **1.64×** gate overhead of
the 2-bit-per-trit encoding measured on the MAC array (`src/accel/`)?

Everything below is measured with ngspice on the SkyWater production
device models (`task tierb`; `fetch_pdk.sh` fetches and patches the model
library). Voltage levels: 0 V = −1, 0.9 V = 0, 1.8 V = +1 (VDD = 1.8 V
plus a shared mid-rail).

## Finding 1: the three-threshold palette exists in a production PDK

Constant-current thresholds (1 µA at W = 1 µm, measured):

| Device | Vt | at the mid level (0.9 V) |
|---|---|---|
| `nfet_01v8_lvt` | 0.63 V | ON |
| `nfet_01v8` | 0.71 V | on (weak) |
| `nfet_g5v0d10v5` (thick-ox) | 0.91 V | ≈ OFF (marginal) |
| `pfet_01v8_lvt` | 0.52 V | ON |
| `pfet_01v8_hvt` | 0.92 V | ≈ off (marginal) |
| `pfet_g5v0d10v5` (thick-ox) | 1.14 V | solidly OFF |

Mixing the 1.8 V low-Vt devices with the 5 V thick-oxide devices run at
1.8 V gives exactly the low/high threshold split ternary CMOS needs — no
exotic devices, no process changes.

## Finding 2: the ternary inverter family works, with clean levels

Measured transfer characteristics (output at the three input levels):

| Cell | T | in = 0 V | in = 0.9 V | in = 1.8 V |
|---|---|---|---|---|
| NTI (detects −1) | 2 | 1.800 | 0.003 | 0.000 |
| PTI (detects ≤0) | 2 | 1.800 | 1.763 | 0.000 |
| STI (ternary negation) | 18 | 1.800 | **0.900** | 0.000 |

The STI (decoder + transmission-gate mux to the mid-rail) produces an
*exact* middle level. The mid-rail itself carries no static current.

## Finding 3: speed and energy parity, through half the wires

Identical 6-transition input sequences, 10 fF loads, the binary control
driving **both** wires of the 2-bit encoding through its own buffers:

| Metric | STI (ternary, 18T) | binary tneg (NOR2 + wire, 4T + drivers) |
|---|---|---|
| propagation delay | 228–266 ps | 221 ps |
| energy, full sequence | 211 fJ | 200 fJ |
| signal wires | **1** | 2 |

A single ternary wire moved log₂3 ≈ 1.58 bits per symbol at the same
energy and speed as the two-wire binary encoding. The **interconnect
argument survives contact with a real PDK**.

For the detector direction (ternary's natural receive operation):

| Metric | NTI (2T) | binary is(−1) = NOR2 (4T) |
|---|---|---|
| transistors | **2** | 4 |
| delay | 333–359 ps | 79 ps |
| energy, full sequence | 79 fJ | 74 fJ |

Half the transistors at energy parity; 4× slower at minimum size because
the thick-oxide PFET is weak (a sizing trade, not a wall).

## Finding 4: the costs, measured plainly

- **Static power at the mid level.** Any stage with a marginal device
  half-on at 0.9 V draws crowbar current: 5.9 µA as first drawn, reduced
  to **1.0 µA** by starving the marginal thick-ox NFET (L = 1.0 µm) —
  still ~1000× binary leakage, paid whenever an input *sits* at logic 0.
  This is voltage-mode ternary's structural weakness at this node.
- **Full ternary logic is transistor-expensive.** The STI needs 18
  transistors against 4 for its 2-bit binary equivalent. The exactness of
  the middle level costs a decoder and a mux; pass-transistor shortcuts
  stall a threshold-drop away from 0.9 V (tried, measured, rejected).
- **"Free negation" is an encoding property, not a voltage property.**
  In the 2-bit encoding, negation is a wire swap plus a NOR — in voltage
  mode it is a full 18T STI. The datapath wins VTX1 measured all day
  (tritwise negation, the MAC's select/negate lanes) belong to the
  *encoding*, and binary substrates execute that encoding well.

## Verdict

Stated the way this project states things:

1. **The architecture-level results stand on binary substrates.** The
   2-bit trit encoding that the whole VTX1 RTL uses is not a compromise
   awaiting rescue by ternary devices — at the logic-gate level it is
   *cheaper* than voltage-mode ternary CMOS in this PDK (4T vs 18T for
   negation; the 1.64× emulation overhead is the better deal).
2. **Ternary signaling earns its keep on wires, not gates.** Energy and
   speed parity per symbol through half the interconnect, with 2T
   receivers, is a real, measured win — relevant to buses, crossbars and
   memory interfaces where wiring dominates, and consistent with where
   the industry actually ships multi-level signaling (PAM-3/MLC flash).
3. **A native ternary *logic* substrate needs what Sky130 cannot give:**
   devices that are hard-off at half swing without thick-oxide weakness,
   and a mid level that costs no static power. That is a device-research
   target (and the literature's CNTFET/RRAM territory), not a design
   deficiency fixable at the circuit level here.

Limitations: typical corner only, schematic-level (no layout
parasitics), minimum-adjacent sizing, one topology per cell, inverter
and detector cells only — a ternary full adder in this style would
multiply the STI's transistor problem, which is the point of finding 4.
