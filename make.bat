setlocal
  SET PATH=E:\bin\apps\mingw\current\bin;%PATH%
  SET GOARCH=amd64

  SET GOOS=windows
  SET CGO_ENABLED=1

  set "EXE_VERSION=0.0.1"
  if not exist .dist mkdir .dist
  go build -ldflags="-H=windowsgui -s -w \
    -X pdfimporter/config.Mode=production \
    -X pdfimporter/config.ExeVersion=%EXE_VERSION%" ^
    -o distbin/findmark4z.exe ./guiapp
endlocal
