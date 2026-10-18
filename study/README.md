# Function-Level Study: where trits actually pay

The Tier B circuit work ended on a comparison that was true but
misaimed: 18 transistors for voltage-mode *signal* negation against 4
for its binary encoding — the microscope hammering a nail. The fair
question, and the project's original intuition, is at the **function**
level over **equal value ranges**: which math operations benefit from
balanced ternary structure, and which pay the encoding tax?

## Method

Each function is implemented twice over the same value range —
18-trit balanced ternary in the 2-bit encoding (±193,710,244) and
29-bit two's complement (the minimal covering width) — then synthesized
flattened with Yosys. Two metrics: cell count (area proxy) and longest
topological path (`ltp`, delay proxy). The ternary adder and multiplier
are the ALU's own verified structures; the fresh ternary units are
verified against value semantics on 20,000 random vectors
(`tb_funcs.v`). Binary units are idiomatic behavioral Verilog, so
synthesis gives binary its best shot. `task funcstudy` reproduces the
table.

## Results

| Function | ternary cells / depth | binary cells / depth | Verdict |
|---|---|---|---|
| NEG (negate a word) | **18 / 1** | 55 / 28 | **ternary: 3× fewer cells, 28× shallower** |
| ABS | 85 / **18** | 84 / 28 | cells parity, **ternary 1.6× shallower** |
| SCALE (÷radix, round to nearest) | **0 / 0** | 56 / 29 | **ternary: free and exact** |
| SIGN (three-way) | 68 / 34 | **29 / 6** | binary (the sign bit is free) |
| CMP3 (sign of a−b) | 498 / 116 | **157 / 14** | binary |
| ADD (mod range) | 443 / 118 | **147 / 56** | binary |
| MUL (mod range) | 3743 / 173 | **2460 / 59** | binary |

## Reading it

The line that separates the winners from the losers is **carry
propagation**:

- **Carry-free functions are ternary's home turf.** Negation is 18
  independent trit swaps — depth 1, no chain — where binary negation is
  invert-plus-increment, a full carry chain. ABS inherits this
  (conditional tritwise flip; its only depth is the sign scan). Radix
  scale-down with rounding is the purest case: dropping the least
  significant trit *is* round-to-nearest (remainders are −1/0/+1, ties
  impossible), costing **zero gates**, where binary needs a conditional
  increment — another carry chain. These are real datapath operations:
  normalization, quantization, rounding in MAC/DSP pipelines.
- **Carry-propagating arithmetic pays the encoding tax.** The hoped-for
  digit-count advantage (18 stages vs 29) is real but swamped: a
  2-bit-encoded trit full adder costs ~6.5 logic levels per stage
  against ~2 for a binary full adder, so ADD nets ~2× deeper and ~3×
  larger. CMP3 and MUL are adders at heart and follow. SIGN loses for a
  different reason: two's complement carries its sign in one free bit,
  while balanced ternary must scan for the most significant nonzero
  trit.
- **This explains the MAC accelerator's economics.** Its winning lane
  operation — multiply by a weight in {−1, 0, +1} — is select/negate,
  i.e. carry-free, squarely in the winning class; its adder tree is the
  unavoidable carry-propagating part. Ternary-leaning designs should
  maximize work done in the carry-free layer per pass through an adder.

Two layers must not be confused: at the **ISA level**, `CMP` and `ABS`
win cycles because they replace multi-instruction dependent chains —
that result stands (measured against the issue-width-matched RV32
core). At the **gate level**, the same operations cost more cells than
their binary counterparts. Both are true; one instruction that costs
1.6× the gates but replaces three dependent instructions is a good
trade, which is precisely what the matched-core cycle numbers showed.

## Answer to the design question

Can functions be built where trits benefit specific math? Yes — and the
recipe is now measured rather than theoretical: **compose from the
carry-free layer** (negate, absolute value, trit select, radix
scaling/rounding, tritwise min/max) **and amortize the carry layer**
(shared adders, reduction trees). Candidates that fit the recipe, in
rough order of promise:

1. **Quantize/normalize pipelines** — chains of scale-by-3^k with
   nearest rounding are literally free wiring between adders.
2. **Ternary-weight dot products** (built: `src/accel/`) — the lane is
   carry-free; only the reduction tree pays.
3. **Sign-domain algorithms** (median/rank filters, sign-LMS adaptive
   filters, decision trees over three-way compares) — dominated by
   NEG/ABS/select with rare full compares.
4. **CORDIC with zero-skip** — balanced digits {−1,0,+1} make the
   "no-rotation" case explicit; iterations become select/negate plus
   one add, and zero digits skip the add entirely (the todo.md idea,
   now with a measured basis).

Caveats: generic Yosys/ABC synthesis (no technology mapping beyond its
default cell set), depth is a delay proxy not timing closure, and the
ternary priority scan in SIGN/CMP3 is a linear chain that a log-tree
first-nonzero implementation would shorten (it would not change the
verdicts, only the margins).
