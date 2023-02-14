bash ./mk_data.sh &&\
rm -rf  ./ppd.app/Contents/MacOS/ppd &&\
CGO_ENABLED=1 go build -ldflags="-w" -a -o ./ppd.app/Contents/MacOS/ppd *.go && \
./ppd.app/Contents/MacOS/ppd 
