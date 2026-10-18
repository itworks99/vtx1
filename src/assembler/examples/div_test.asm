;===============================================================================
; div_test.asm - DIV/MOD semantics test for the iterative divider
;===============================================================================
; Exercises the ratified division semantics (matching vtxsim / Go):
; truncation toward zero, remainder sign follows the dividend.
; Results at 0x1200..0x1205, verified both in vtxsim and on the RTL core:
;   -7 / 2 = -3     -7 % 2 = -1
;    7 / -2 = -3     7 % -2 = 1
;    29524 / 7 = 4217   29524 % 7 = 5
;===============================================================================

        .ORG 0x1000

        LEA T0, -7
        LEA T1, 2
        DIV T2, T0, T1
        MOD T3, T0, T1
        ST T2, 0x1200
        ST T3, 0x1201

        LEA T0, 7
        LEA T1, -2
        DIV T2, T0, T1
        MOD T3, T0, T1
        ST T2, 0x1202
        ST T3, 0x1203

        LEA T0, 29524
        DIV T2, T0, 7
        MOD T3, T0, 7
        ST T2, 0x1204
        ST T3, 0x1205

        WFI

        .ORG 0x1200
results: .DW 0, 0, 0, 0, 0, 0
