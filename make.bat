setlocal
  @rem x32 https://github.com/brechtsanders/winlibs_mingw/releases/download without LLVM/Clang/LLD/LLDB`
  @rem SET PATH=E:\bin\mingw32\bin;%PATH%
  @rem SET GOARCH=386

  @rem x64 scoop install mingw
  SET PATH=E:\bin\apps\mingw\current\bin;%PATH%
  SET GOARCH=amd64

  SET GOOS=windows
  SET CGO_ENABLED=1

  @rem go build -ldflags "-H=windowsgui -s -w" -o distbin/report.exe ./guiapp 
  go build -ldflags="-H=windowsgui -s -w -X 'pdfimporter/internal/entity.Mode=production'" ^
           -o distbin/findmark4z.exe ./guiapp

  @rem go build -o guiapp/findmark.exe ./guiapp 


endlocal
