@echo off
echo Running VTX1 Assembler Tests...
python3 test_assembler.py
if %ERRORLEVEL% EQU 0 (
    echo Tests completed successfully!
) else (
    echo Tests failed with error code %ERRORLEVEL%
)
pause
