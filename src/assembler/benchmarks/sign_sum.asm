;===============================================================================
; sign_sum.asm - sum of sign(a[i] - b[i])
;===============================================================================
; Tier A benchmark kernel: the three-valued compare showcase. CMP produces
; the full balanced sign (-1/0/+1) as a VALUE in one instruction; a binary
; ISA needs an slt/slt/sub sequence (3 instructions) for the same result.
; Expected result: -2
;===============================================================================

        .ORG 0x1000

main:
        LEA TB, vec_a
        LEA T5, vec_b
        LD  T0, length
        LEA T1, 0            ; sum of signs
        LEA T2, 0            ; index

loop:
        LD  T3, [TB+T2]      ; a[i]
        LD  T4, [T5+T2]      ; b[i]
        CMP T3, T3, T4       ; T3 = sign(a[i] - b[i]) in ONE instruction
        ADD T1, T1, T3
        ADD T2, T2, 1
        BLT T2, T0, loop

        ST  T1, result
        WFI

        .ORG 0x1100
length: .DW 16
vec_a:  .DW 3, -5, 7, 2, -8, 4, 1, -2, 9, -3, 6, -7, 5, 0, 2, -4
vec_b:  .DW 2, 3, -1, 5, 2, -6, 8, 4, -2, 7, 1, -3, -5, 9, 0, 6

        .ORG 0x1200
result: .DW 0
