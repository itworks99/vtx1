;===============================================================================
; hello_world.asm - Simple Hello World example for VTX1
;===============================================================================
; This program demonstrates a basic "Hello World" program on the VTX1 architecture
; It outputs the string "Hello, VTX1!" to memory location 0x2000
;===============================================================================

        .ORG 0x1000          ; Start at address 0x1000

;===============================================================================
; Main program
;===============================================================================
main:
        ; Initialize registers
        LD T0, string        ; Load address of string
        LD T1, 0x2000        ; Destination address for output
        LD T2, 0             ; Initialize counter

loop:
        LD T3, [T0]          ; Load character from string
        BEQ T3, 0, done      ; If null terminator (0), we're done
        ST T3, [T1]          ; Store character to output

        ; Use VLIW to increment pointers and counter in parallel
        [ADD T0, T0, 1] [ADD T1, T1, 1] [ADD T2, T2, 1]

        JMP loop             ; Repeat for next character

done:
        ; Store the character count at 0x2100
        ST T2, 0x2100

        ; Signal completion (store 1 at 0x2104 to indicate program finished)
        LD T0, 1
        ST T0, 0x2104

        ; Wait for further instructions
        WFI

;===============================================================================
; Data section
;===============================================================================
string:
        .DB "Hello, VTX1!", 0  ; Null-terminated string
