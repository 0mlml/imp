#include "peripherals/at_timers.h"
#include "peripherals/at_button.h"
#include "peripherals/at_led.h"
#include "asset_tracker.h"

#include <zephyr/logging/log.h>
LOG_MODULE_REGISTER(at_timers, CONFIG_TRACKER_LOG_LEVEL);

static void scan_timer_cb(struct k_timer *timer_id);
static void btn_press_timer_cb(struct k_timer *timer_id);

K_TIMER_DEFINE(scan_timer, scan_timer_cb, NULL);
K_TIMER_DEFINE(btn_press_timer, btn_press_timer_cb, NULL);

bool ble_timeout = false;

static void scan_timer_cb(struct k_timer *timer_id)
{
	ARG_UNUSED(timer_id);

	at_event_send(EVENT_SCAN_SENSORS);
	at_event_send(EVENT_SEND_UPLINK);

    scan_timer_set_and_run(K_MSEC(300)); 

}

static void btn_press_timer_cb(struct k_timer *timer_id)
{
	ARG_UNUSED(timer_id);
	button_long_press = true;
	LOG_INF("Long button press...");
}

void scan_timer_set_and_run(k_timeout_t delay)
{
	k_timer_start(&scan_timer, delay, Z_TIMEOUT_NO_WAIT);
}

void btn_press_timer_set_and_run(void)
{
	button_long_press = false;
	k_timer_start(&btn_press_timer, K_MSEC(CONFIG_LONG_PRESS_PER_MS), Z_TIMEOUT_NO_WAIT);
}

void btn_press_timer_stop()
{
	k_timer_stop(&btn_press_timer);
}

