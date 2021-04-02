#include <X11/Xlib.h>

int c_wm_x11_get_type_of_event(XEvent);
Window c_wm_x11_query_window_from_array(Window*, int);
void c_wm_x11_send_event_destroy(Display*, Window);
void c_wm_x11_set_type_of_event(XEvent*, int);
char* c_wm_x11_get_window_title(Display*, Window);