@echo off
setlocal
cd /d "%~dp0"
where go >nul 2>nul || (echo [LOI] Chua cai Go 1.23+ & pause & exit /b 1)
where node >nul 2>nul || (echo [LOI] Chua cai Node.js 22+ & pause & exit /b 1)

if not exist "frontend\node_modules" (
  echo Dang cai thu vien frontend...
  pushd frontend
  call npm ci || exit /b 1
  popd
)

start "Caro Backend" cmd /k "cd /d %~dp0backend && go run ./cmd/server"
start "Caro Frontend" cmd /k "cd /d %~dp0frontend && set VITE_API_BASE_URL=http://localhost:5207&& npm run dev -- --host 0.0.0.0"

echo.
echo Frontend: http://localhost:5173
echo Backend : http://localhost:5207
echo Hai cua so lenh da duoc mo. Dung Ctrl+C trong tung cua so de tat.
pause
