bash ./mk_data.sh &&\
CGO_ENABLED=1 go build -ldflags="-w" -a -o ./cg.app/Contents/MacOS/cg *.go && \
(umount /Volumes/cg ||\
hdiutil attach cg.dmg) &&\
rm -rf /Volumes/cg/cg &&\
cp -R cg.app /Volumes/cg &&\
cp ~/Desktop/Applications /Volumes/cg &&\
hdiutil detach /Volumes/cg &&\
rm -rf ./cg_distrbute.dmg &&\
hdiutil convert ./cg.dmg -format UDCO -o ./cg_distrbute.dmg &&\
pkill -f cg.app
./cg.app/Contents/MacOS/cg 
#open ./cg.app
