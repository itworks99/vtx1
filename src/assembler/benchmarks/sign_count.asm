;===============================================================================
; sign_count.asm - three-way sign classification
;===============================================================================
; Tier A benchmark kernel playing to balanced ternary's strength: CMP yields
; the full three-valued sign (-1/0/+1) in one instruction, and the counts
; fall out of two branches.
; Expected results: 6 negative, 4 zero, 6 positive
;===============================================================================

        .ORG 0x1000

main:
        LEA TB, data
        LD  T0, length
        LEA T1, 0            ; negative count
        LEA T2, 0            ; zero count
        LEA T3, 0            ; positive count
        LEA T4, 0            ; index

loop:
        LD  T5, [TB+T4]
        CMP T6, T5, 0        ; T6 = sign(element): -1, 0, or +1
        BLT T6, 0, is_neg
        BGT T6, 0, is_pos
        ADD T2, T2, 1        ; zero
        JMP next
is_neg:
        ADD T1, T1, 1
        JMP next
is_pos:
        ADD T3, T3, 1
next:
        ADD T4, T4, 1
        BLT T4, T0, loop

        ST  T1, result_neg
        ST  T2, result_zero
        ST  T3, result_pos
        WFI

        .ORG 0x1100
length: .DW 16
data:   .DW 5, -3, 0, 7, -1, 0, 2, -9, 4, 0, -6, 1, -2, 8, 0, -5

        .ORG 0x1200
result_neg:  .DW 0
result_zero: .DW 0
result_pos:  .DW 0
