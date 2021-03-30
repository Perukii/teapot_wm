
package main

import "log"
import "os"

func wm_debug_log_file_init(log_file *os.File){
	var path string
	{
		var err error
		path, err = os.Executable()
		if err != nil { log.Fatalf("WmError: %v", err) }
	}
	{
		var err error
		log_file, err = os.OpenFile(path + "log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil { log.Fatalf("WmError: %v", err) }
	}

	log.SetOutput(log_file)
}

func wm_debug_log(txt string){
	log.Println(txt)
}
