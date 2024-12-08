#ifndef ASSET_TRACKER_H
#define ASSET_TRACKER_H

#define RECEIVE_TASK_STACK_SIZE (4096)
#define RECEIVE_TASK_PRIORITY (CONFIG_SIDEWALK_THREAD_PRIORITY + 1)

#define LINK_DOWN 0
#define LINK_UP 1

typedef enum uplink_msg_types {
	SERIAL_PUSH_T,
	NOLOC_T,
	WIFI_T,
	GNSS_T,
	UNSET_T,
} uplink_msg_t;

struct at_sensors {
	double temp;
	double hum;
	double max_accel_x;
	double max_accel_y;
	double max_accel_z;
	double peak_accel;
};

enum at_sidewalk_state {
	STATE_SIDEWALK_INIT,
	STATE_SIDEWALK_READY,
	STATE_SIDEWALK_NOT_READY,
	STATE_SIDEWALK_SECURE_CONNECTION,
};

enum at_state {
	AT_STATE_INIT,
	AT_STATE_RUN
};

struct at_config{
    uint16_t max_rec;
    uint8_t motion_period;
    uint8_t scan_freq_motion;
    uint8_t motion_thres;
    uint8_t scan_freq_static;
};

typedef struct at_context {
	enum at_state state;
	uplink_msg_t uplink_type;
	bool connection_request;
	bool motion;
	struct at_sensors sensors;
    struct at_config at_conf;
} at_ctx_t;

typedef enum at_events {
	EVENT_NONE,
	BUTTON_EVENT_SHORT,
	BUTTON_EVENT_LONG,
    MOTION_EVENT,
	EVENT_SEND_UPLINK,
    EVENT_SCAN_SENSORS,
	EVENT_UPLINK_COMPLETE,
} at_event_t;


void at_event_send(at_event_t event);

int at_thread_init(void);

#endif /* ASSET_TRACKER_H */