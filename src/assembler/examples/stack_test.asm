;===============================================================================
; stack_test.asm - the ratified stack ABI, in vtxsim and on the RTL core
;===============================================================================
; Stack ABI: TB (base pointer) is the stack pointer; full-descending,
; word-addressed. PUSH r: TB <- TB-1, mem[TB] <- r (old r).
; POP r: r <- mem[TB], TB <- TB+1. For POP TB the loaded value wins.
;
; Results at 0x1200..0x1205: 33, 22, 11, 5, 7936 (balanced TB = 0x1F00),
; 4242 (the POP TB edge case).
;===============================================================================

        .ORG 0x1000

        LEA TB, 0x1F00       ; initialize the stack pointer

        ; LIFO order
        LEA T0, 11
        LEA T1, 22
        LEA T2, 33
        PUSH T0
        PUSH T1
        PUSH T2
        POP T3               ; 33
        POP T4               ; 22
        POP T5               ; 11
        ST T3, 0x1200
        ST T4, 0x1201
        ST T5, 0x1202

        ; Nested calls: f1 saves the link register around its call to f2
        LEA T0, 0
        CALL f1
        ST T0, 0x1203        ; 5 (f2 adds 2, f1 adds 3)
        ST TB, 0x1204        ; stack balanced again: 0x1F00 = 7936

        ; POP TB: the loaded value supersedes the increment
        LEA T4, 4242
        PUSH T4
        POP TB               ; TB <- 4242
        ST TB, 0x1205

        WFI

f1:
        PUSH TA              ; save the link register
        CALL f2
        ADD T0, T0, 3
        POP TA               ; restore it
        RET

f2:
        ADD T0, T0, 2
        RET

        .ORG 0x1200
results: .DW 0, 0, 0, 0, 0, 0
