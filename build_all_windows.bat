@echo off
setlocal EnableExtensions
cd /d "%~dp0"

where go >nul 2>nul || (echo [LOI] Chua cai Go 1.23+ & pause & exit /b 1)
where node >nul 2>nul || (echo [LOI] Chua cai Node.js 22+ & pause & exit /b 1)

if "%VITE_API_BASE_URL%"=="" set "VITE_API_BASE_URL=http://localhost:5207"

if not exist dist mkdir dist

echo [1/3] Build backend Windows...
pushd backend
go build -trimpath -ldflags="-s -w" -o "..\dist\caro-server.exe" ./cmd/server || (popd & pause & exit /b 1)
popd

echo [2/3] Cai thu vien frontend...
pushd frontend
call npm ci || (popd & pause & exit /b 1)

echo [3/3] Build frontend tinh...
set "BASE_PATH="
call npm run build || (popd & pause & exit /b 1)
popd

if exist "dist\frontend" rmdir /s /q "dist\frontend"
xcopy "frontend\build" "dist\frontend\" /E /I /Y >nul

echo.
echo HOAN TAT:
echo - Backend: dist\caro-server.exe
echo - Frontend: dist\frontend
pause
