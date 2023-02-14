bash ./mk_data.sh &&\
CGO_ENABLED=1 go build -ldflags="-w" -a -o ./ppd.app/Contents/MacOS/ppd *.go && \
(umount /Volumes/ppd_pre ||\
hdiutil attach ppd_pre.dmg) &&\
rm -rf /Volumes/ppd_pre/ppd &&\
cp -R ppd.app /Volumes/ppd_pre &&\
cp ~/Desktop/Applications /Volumes/ppd_pre &&\
hdiutil detach /Volumes/ppd_pre &&\
rm -rf ./ppd.dmg &&\
hdiutil convert ./ppd_pre.dmg -format UDCO -o ./ppd.dmg &&\
pkill -f ppd.app
./ppd.app/Contents/MacOS/ppd 
#open ./ppd.app
