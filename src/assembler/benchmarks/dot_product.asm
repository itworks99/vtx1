;===============================================================================
; dot_product.asm - 16-element signed dot product
;===============================================================================
; Tier A benchmark kernel: signed multiply-accumulate over two vectors.
; Expected result: -107
;===============================================================================

        .ORG 0x1000

main:
        LEA TB, vec_a        ; base of first vector
        LEA T5, vec_b        ; base of second vector
        LD  T0, length       ; element count
        LEA T1, 0            ; accumulator
        LEA T2, 0            ; index

loop:
        LD  T3, [TB+T2]      ; a[i]
        LD  T4, [T5+T2]      ; b[i]
        MUL T3, T3, T4       ; a[i] * b[i]
        ADD T1, T1, T3       ; acc += product
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
