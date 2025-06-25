; Test file with multiple deliberate errors to demonstrate error reporting
; This file contains multiple errors that should be detected and reported

.org 0x1000

start:
    ; Lexical error - invalid character
    add t0, t1, @t2  ; @ is not valid in this context

    ; Syntax error - missing operand
    sub t3,          ; Missing second operand

    ; Another syntax error - extra commas
    mul t4, t5,, t6  ; Extra comma between operands

    ; Undefined symbol
    jmp undefined_label  ; Symbol not defined

duplicate_label:
    nop

    ; Invalid register
    add tx, t1, t2   ; tx is not a valid register

    ; Missing closing bracket in memory reference
    ld t0, [t1+10    ; Missing closing bracket

duplicate_label:     ; Duplicate label definition
    nop

    ; Invalid number format
    .db 0xZZ         ; Invalid hex digits

end:
    halt
