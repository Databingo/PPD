bash ./mk_data.sh 
rsrc -manifest myapp.manifest -ico ico.ico -o myapp.syso
GOOS=windows GOARCH="386" CGO_ENABLED=1 CC="i686-w64-mingw32-gcc" go build -ldflags="-H windowsgui" -a -o cg.exe
#GOOS=windows GOARCH="386" CGO_ENABLED=1 CC="i686-w64-mingw32-gcc" go build  -a -o cg.exe
mv cg.exe ppd.exe
zip -p ppd ppd.exe
rm -rf ppd.exe
unzip ppd.zip
