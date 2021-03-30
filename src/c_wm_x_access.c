#include <X11/Xlib.h>

int c_wm_x11_get_type_of_event(XEvent event){
	return event.type;
}


void c_wm_x11_set_type_of_event(XEvent* event, int type){
	event->type = type;
}

void c_wm_x11_send_event(Display* display, Window window, const char* event_data){
	
	XEvent delete_event;
	delete_event.xclient.type = ClientMessage;
	delete_event.xclient.message_type = XInternAtom(display, "WM_PROTOCOLS", True);
	delete_event.xclient.format = 32;
	delete_event.xclient.data.l[0] = XInternAtom(display, event_data, True);
	delete_event.xclient.data.l[1] = CurrentTime;
	delete_event.xclient.window = window;

	// 除去イベントを送信
	XSendEvent(display, window, False, NoEventMask, &delete_event);
}


Window c_wm_x11_query_window_from_array(Window* array, int n){
	return array[n];
}