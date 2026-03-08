@echo off
echo Clearing GitHub credentials...
echo.

REM Create a temporary file with credential info
echo protocol=https> temp_cred.txt
echo host=github.com>> temp_cred.txt
echo.>> temp_cred.txt

REM Erase the credential
type temp_cred.txt | git credential-manager erase

REM Clean up
del temp_cred.txt

echo.
echo GitHub credentials cleared!
echo.
echo Now you can push again and enter PrachiPatel2105 credentials.
echo.
pause
