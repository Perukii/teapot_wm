
XINITRC=~/.xinitrc

all:
	go build -o wm ./src

set_xinitrc:
	echo "${CURDIR}/wm & xterm" > $(XINITRC)

setup_log_file:
	echo "" > ./wmlog.txt

show_log_file:
	cat ./wmlog.txt