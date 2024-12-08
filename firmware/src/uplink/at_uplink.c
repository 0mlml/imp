#include <zephyr/kernel.h>
#include <zephyr/logging/log.h>
LOG_MODULE_REGISTER(at_uplink, CONFIG_TRACKER_LOG_LEVEL);

#include <asset_tracker.h>
#include <uplink/at_uplink.h>


void at_send_uplink(at_ctx_t *context) 
{

	at_ctx_t *at_ctx = (at_ctx_t *)context;

	uint8_t uptype = (uint8_t) at_ctx->uplink_type;
	
	switch (uptype) {
		case SERIAL_PUSH_T: 
			LOG_INF("SERIALPUSH>temp:%f,hum:%f,motx:%f,moty:%f,motz:%f,peak:%f<SERIALPUSH", at_ctx->sensors.temp, at_ctx->sensors.hum, at_ctx->sensors.max_accel_x, at_ctx->sensors.max_accel_y, at_ctx->sensors.max_accel_z, at_ctx->sensors.peak_accel);
			break;

		default:
			LOG_ERR("Invalid uplink type!");
	}
}

void at_msg_sent(at_ctx_t *context) {
	at_event_send(EVENT_SEND_UPLINK);
}

void at_send_error(at_ctx_t *context) {
	// nop
}

