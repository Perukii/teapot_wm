#include <X11/Xlib.h>
#include <X11/Xutil.h>

int c_wm_x11_get_type_of_event(XEvent event){
	return event.type;
}

void c_wm_x11_set_type_of_event(XEvent* event, int type){
	event->type = type;
}

void c_wm_x11_send_event_destroy(Display* display, Window window){
	
	XEvent event;
	event.xclient.type = ClientMessage;
	event.xclient.message_type = XInternAtom(display, "WM_PROTOCOLS", True);
	event.xclient.format = 32;
	event.xclient.data.l[0] = XInternAtom(display, "WM_DELETE_WINDOW", True);
	event.xclient.data.l[1] = CurrentTime;
	event.xclient.window = window;

	XSendEvent(display, window, False, NoEventMask, &event);
}


Window c_wm_x11_query_window_from_array(Window* array, int n){
	return array[n];
}

char* c_wm_x11_get_window_title(Display* display, Window window){
	char* title;
	XTextProperty proper;
	XGetWMName(display, window, &proper);
	title = proper.value;
	return title;
}
