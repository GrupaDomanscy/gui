$ErrorActionPreference = "Stop"

if (Test-Path .\test.exe) {
    Remove-Item .\test.exe;
}

go test -c -o test.exe;
.\test.exe;