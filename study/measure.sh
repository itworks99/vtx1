#!/bin/sh
# Synthesize each function pair flattened and report cells + logic depth.
# Usage: ./study/measure.sh  (run from the repo root)
set -e
SRCS="src/cpu/vtx1_alu.v study/t_funcs.v study/b_funcs.v"
for m in t_neg b_neg t_abs b_abs t_sign b_sign t_cmp3 b_cmp3 t_add b_add t_mul b_mul t_scale b_scale; do
    out=$(yosys -p "read_verilog -I src/common $SRCS; synth -top $m -flatten; ltp -noff; stat" 2>/dev/null)
    cells=$(printf "%s\n" "$out" | grep -E '^ +[0-9]+ cells$' | tail -1 | awk '{print $1}')
    scope=$(printf "%s\n" "$out" | grep -F '$scopeinfo' | tail -1 | awk '{print $1}')
    depth=$(printf "%s\n" "$out" | grep "Longest topological path" | tail -1 | sed 's/.*length=\([0-9]*\).*/\1/')
    case "$scope" in ''|*[!0-9]*) scope=0;; esac
    case "$cells" in ''|*[!0-9]*) cells=0;; esac
    echo "$m cells=$((cells - scope)) depth=$depth"
done
