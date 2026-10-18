;===============================================================================
; tmac_matvec.asm - 16x16 ternary-weight matrix-vector product, software vs
;                   the TMAC accelerator (RTL/SoC only: uses MMIO)
;===============================================================================
; The same matvec is computed twice: a software MUL/ADD loop, then the
; ternary-weight MAC accelerator (peripheral 6; the ISA maps negative
; addresses to MMIO, so TMAC register k lives at address -(0x60000+k)).
; The program compares the two result vectors itself and stores:
;   0x1200: mismatch count (expect 0)
;   0x1201: completion flag (expect 1)
; Marker stores to 0x1210/0x1211/0x1212 let the testbench timestamp the
; phase boundaries (end of packing / end of software / end of accelerator).
;
; Expected matvec (TB cross-check of result 0):
;   18, 28, 1, -3, 14, -2, 1, 6, -15, -26, 14, 11, 17, 30, -12, 17
;===============================================================================

        .ORG 0x1000

;-------------------------------------------------------------------------------
; Phase 0: pack each weight row into one word (lane i at weight 3^i).
; Offline in real use - weights ship packed; done here so both phases
; derive from the same unpacked data.
;-------------------------------------------------------------------------------
        LEA T0, wdata        ; row element pointer base
        LEA T1, 0x1260       ; packed row destination
        LEA T2, 16           ; row counter
pack_row:
        LEA T3, 0            ; packed word
        LEA T4, 15           ; i = 15 .. 0 (most significant lane first)
pack_elem:
        LD  T5, [T0+T4]      ; w[r][i]
        LEA T6, 3
        MUL T3, T3, T6
        ADD T3, T3, T5
        DEC T4
        BGE T4, 0, pack_elem
        ST  T3, [T1]
        INC T1
        LEA T6, 16
        ADD T0, T0, T6
        DEC T2
        BGT T2, 0, pack_row

        ST  T2, 0x1210       ; marker: packing done

;-------------------------------------------------------------------------------
; Phase 1: software matvec (MUL/ADD loop)
;-------------------------------------------------------------------------------
        LEA TB, act          ; activation base (stays in TB; no stack use here)
        LEA T0, wdata
        LEA T1, 0x1220       ; software results
        LEA T2, 16           ; row counter
sw_row:
        LEA T3, 0            ; acc
        LEA T4, 0            ; i
sw_elem:
        LD  T5, [T0+T4]      ; w[r][i]
        LD  T6, [TB+T4]      ; act[i]
        MUL T5, T5, T6
        ADD T3, T3, T5
        INC T4
        LEA T6, 16
        BLT T4, T6, sw_elem
        ST  T3, [T1]
        INC T1
        LEA T6, 16
        ADD T0, T0, T6
        DEC T2
        BGT T2, 0, sw_row

        ST  T2, 0x1211       ; marker: software matvec done

;-------------------------------------------------------------------------------
; Phase 2: accelerator matvec.
; TMAC registers (negative addresses): ACT[i] at base-i, ACC at base-16,
; WEIGHT/step at base-17, with base = -0x60000.
;-------------------------------------------------------------------------------
        LD  T0, tmac_base    ; T0 = -0x60000
        LEA T4, 0
hw_act:
        LD  T5, [TB+T4]      ; act[i]
        NEG T6, T4
        ST  T5, [T0+T6]      ; ACT[i]
        INC T4
        LEA T6, 16
        BLT T4, T6, hw_act

        LEA TB, 0            ; zero source for the clear store
        LEA T1, 0x1260       ; packed weight rows
        LEA T2, 0x1240       ; accelerator results
        LEA T3, 16           ; row counter
        LEA T4, -16          ; ACC offset
        LEA T5, -17          ; WEIGHT offset
hw_row:
        LD  T6, [T1]         ; packed row
        ST  TB, [T0+T4]      ; clear accumulator
        ST  T6, [T0+T5]      ; weight word write = one 16-lane MAC step
        LD  T6, [T0+T4]      ; read accumulator
        ST  T6, [T2]
        INC T1
        INC T2
        DEC T3
        BGT T3, 0, hw_row

        ST  T3, 0x1212       ; marker: accelerator matvec done

;-------------------------------------------------------------------------------
; Phase 2b: DMA-fed matvec. One 4-register descriptor (activation base,
; packed-weight base, result base, row count) and a status poll - the
; accelerator streams everything from the data scratchpad itself.
;-------------------------------------------------------------------------------
        LD  T0, tmac_base    ; T0 = -0x60000
        LEA T1, 0x1100       ; ACT_BASE: the activation vector
        LEA T4, -18          ; register 0x12
        ST  T1, [T0+T4]
        LEA T1, 0x1260       ; W_BASE: the packed weight rows
        LEA T4, -19          ; register 0x13
        ST  T1, [T0+T4]
        LEA T1, 0x1280       ; R_BASE: DMA results
        LEA T4, -20          ; register 0x14
        ST  T1, [T0+T4]
        LEA T1, 16           ; 16 rows + START
        LEA T4, -21          ; register 0x15
        ST  T1, [T0+T4]
dma_poll:
        LD  T5, [T0+T4]      ; status: +1 while busy
        BNE T5, 0, dma_poll

        LEA T3, 0
        ST  T3, 0x1213       ; marker: DMA matvec done

;-------------------------------------------------------------------------------
; Phase 3: compare both accelerator result vectors against software
;-------------------------------------------------------------------------------
        LEA T0, 0x1220
        LEA T1, 0x1240
        LEA T6, 0x1280
        LEA T2, 16
        LEA T3, 0            ; mismatch count
cmp_loop:
        LD  T4, [T0]
        LD  T5, [T1]
        BEQ T4, T5, cmp_mmio_ok
        INC T3
cmp_mmio_ok:
        LD  T5, [T6]
        BEQ T4, T5, cmp_dma_ok
        INC T3
cmp_dma_ok:
        INC T0
        INC T1
        INC T6
        DEC T2
        BGT T2, 0, cmp_loop

        ST  T3, 0x1200       ; mismatches (expect 0)
        LEA T4, 1
        ST  T4, 0x1201       ; done flag
        WFI

;-------------------------------------------------------------------------------
        .ORG 0x1100
act:    .DW 3, -5, 7, 2, -8, 4, 1, -2, 9, -3, 6, -7, 5, 0, 2, -4
tmac_base: .DW -393216

        ; the 256-word weight matrix gets its own region, clear of the
        ; result/marker page at 0x1200
        .ORG 0x1400
wdata:
        .DW 1, -1, -1, 1, 0, -1, -1, -1, 1, -1, 1, 1, 1, -1, 1, 0
        .DW -1, -1, -1, -1, -1, 1, 1, -1, 1, -1, 1, 1, 1, 1, 0, -1
        .DW 0, 1, 0, -1, -1, 1, 0, 0, 0, -1, -1, 0, -1, -1, 0, -1
        .DW 0, 0, 1, 0, -1, 1, 0, 1, -1, 0, -1, 1, 0, 1, 1, 0
        .DW 1, -1, 1, -1, -1, 1, -1, 0, -1, -1, -1, 0, 0, 0, 1, 0
        .DW -1, 0, 0, -1, 1, 0, 1, 1, 1, -1, 1, 1, -1, 1, 1, -1
        .DW -1, 0, 0, 0, 1, 1, 1, -1, 1, 0, -1, -1, -1, 0, 0, 0
        .DW -1, -1, 1, 1, 0, -1, 1, 0, 0, 1, 0, -1, 0, -1, -1, 1
        .DW 1, 1, 0, 1, 1, 0, 1, 0, 0, -1, -1, 1, 0, -1, -1, -1
        .DW -1, 1, -1, 1, 0, 1, -1, 0, 0, 1, 0, 1, 0, 1, -1, 1
        .DW 1, -1, 1, 1, 0, 1, 0, -1, 0, 0, -1, 0, -1, 1, 1, 0
        .DW 1, -1, 1, -1, 1, 0, 1, 1, 1, -1, -1, 0, -1, 1, 1, -1
        .DW 1, 0, 0, -1, -1, 0, 0, -1, -1, -1, 1, -1, -1, 1, 0, -1
        .DW 1, -1, -1, 1, 0, 1, -1, 0, 1, 1, 0, -1, 1, 1, 1, -1
        .DW 1, 0, 0, 1, 1, 0, 0, 1, 0, -1, -1, -1, -1, 0, -1, 1
        .DW 1, -1, 1, -1, -1, -1, 1, 1, -1, -1, -1, -1, 0, -1, 1, -1

        .ORG 0x1200
results: .DW 0, 0
