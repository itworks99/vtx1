;===============================================================================
; macro_demo.asm - .MACRO and .INCLUDE demonstration for VTX1
;===============================================================================
; Shows the assembler's preprocessor features: file inclusion and
; parameterized macros (including a macro invoking another macro).
;===============================================================================

        .INCLUDE "io_defs.asm"

; Copy a register (VTX1 has no MOV; ADD with immediate 0 is the idiom)
.MACRO MOV dst src
        ADD dst, src, 0
.ENDM

; Store two registers to consecutive words
.MACRO STORE_PAIR a b addr
        ST a, addr
        ST b, addr+2
.ENDM

; Signal completion through the status port
.MACRO DONE
        LEA T6, 1
        ST T6, STATUS_PORT
.ENDM

        .ORG 0x1000

main:
        LEA T0, 42           ; Load a value
        MOV T1, T0           ; Macro: expands to ADD T1, T0, 0
        STORE_PAIR T0, T1, OUTPUT_PORT
        DONE
        WFI
