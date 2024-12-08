#include <zephyr/kernel.h>
#include <zephyr/kernel/thread_stack.h>

#if defined(CONFIG_ASSET_TRACKER_CLI)
#include "at_shell.h"
#endif

#include <zephyr/logging/log.h>
LOG_MODULE_REGISTER(asset_tracker, CONFIG_TRACKER_LOG_LEVEL);

#include <asset_tracker.h>
#include "peripherals/at_lis3dh.h"
#include "peripherals/at_sht41.h"
#include "peripherals/at_timers.h"
#include "uplink/at_uplink.h"

static struct k_thread at_thread;
K_THREAD_STACK_DEFINE(at_thread_stack, 8192);
K_MSGQ_DEFINE(at_thread_msgq, sizeof(at_event_t), 10, 4);

static at_ctx_t asset_tracker_context = {0};

static void at_app_entry(void *ctx, void *unused, void *unused2)
{
	at_ctx_t *at_ctx = (at_ctx_t *)ctx;
	ARG_UNUSED(unused);
	ARG_UNUSED(unused2);
	LOG_DBG("Starting %s ...", __FUNCTION__);

	while (true) {
		at_event_t event = EVENT_NONE;

		if (!k_msgq_get(&at_thread_msgq, &event, K_FOREVER)) {
			switch (event) {				
				case BUTTON_EVENT_SHORT:
					LOG_INF("Immediate scan and uplink triggered...");
					scan_timer_set_and_run(K_MSEC(100));
					break;
				
				case BUTTON_EVENT_LONG:
					LOG_INF("Long press...");
					break;

				case EVENT_SEND_UPLINK:
					LOG_INF("Sending uplink...");
					at_ctx->uplink_type = SERIAL_PUSH_T;
					at_send_uplink(at_ctx);
					break;

				case EVENT_SCAN_SENSORS:
					LOG_INF("Scanning sensors...");
					get_temp_hum(&at_ctx->sensors);
					get_accel(&at_ctx->sensors);
					break;
				default:
					LOG_ERR("Invalid Event received!, %d", event);
			}
		}
	}
}

void at_event_send(at_event_t event)
{
	int ret = k_msgq_put(&at_thread_msgq, (void *)&event,
			     k_is_in_isr() ? K_NO_WAIT : K_FOREVER);

	if (ret) {
		LOG_ERR("Failed to send event to asset tracker thread. err: %d", ret);
	}
}

int at_thread_init(void)
{
	asset_tracker_context.at_conf = (struct at_config) {
		.max_rec = 100,
		.motion_period = CONFIG_IN_MOTION_PER_M,
		.scan_freq_motion = CONFIG_MOTION_SCAN_PER_S,
		.motion_thres = 5,
		.scan_freq_static = CONFIG_STATIC_SCAN_PER_M,
	};

	#if defined(CONFIG_ASSET_TRACKER_CLI)
	AT_CLI_init(&asset_tracker_context);
	#endif

	(void)k_thread_create(&at_thread, at_thread_stack,
			      K_THREAD_STACK_SIZEOF(at_thread_stack), at_app_entry,
			      &asset_tracker_context, NULL, NULL, 14, 0, K_NO_WAIT);

	return 0;
}