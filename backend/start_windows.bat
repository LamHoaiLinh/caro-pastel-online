@echo off
cd /d "%~dp0"
if not exist caro-server.exe (
  echo Dang build backend...
  go build -trimpath -ldflags="-s -w" -o caro-server.exe ./cmd/server || pause
)
set "PORT=5207"
set "MATCH_DB_PATH=data\matches.json"
caro-server.exe
pause
