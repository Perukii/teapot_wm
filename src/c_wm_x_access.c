#include <X11/Xlib.h>

int c_wm_x11_get_type_of_event(XEvent event){
	return event.type;
}

Window c_wm_x11_query_window_from_array(Window* array, int n){
	return array[n];
}