
XINITRC=~/.xinitrc

all:
	go build -o wm ./src

init_xinitrc:
	echo "xcompmgr & ${CURDIR}/wm & xterm" > $(XINITRC)

init_log_file:
	echo "" > ./wmlog.txt

create_setting_file:
	echo "\n{\n\
		\"background_pngfile_path\" : \"DEFAULT\"\n\
	}\n" > ${CURDIR}/mitewm_setting.json

show_log_file:
	cat ./wmlog.txt