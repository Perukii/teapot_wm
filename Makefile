
XINITRC=~/.xinitrc

all:
	go build -o wm ./src

init_xinitrc:
	echo "xcompmgr & ${CURDIR}/wm ${CURDIR}/wm_setting.json & xterm" > $(XINITRC)

init_log_file:
	echo "" > ./wmlog.txt

create_setting_file:
	echo "\n{\n\
		\"UserSetting\" : {\n\
			\"BackgroundPngFilePath\" : \"${CURDIR}/resources/background/Default.png\"\n\
		}\n\
	}\n" > ${CURDIR}/wm_setting.json

show_log_file:
	cat ./wmlog.txt