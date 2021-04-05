
XINITRC=~/.xinitrc

all:
	go build -o wm ./src

init_xinitrc:
	echo "xcompmgr & ${CURDIR}/wm & xterm & . ${CURDIR}/startup.sh" > $(XINITRC)

init_log_file:
	echo "" > ./wmlog.txt

show_log_file:
	cat ./wmlog.txt