#include <stdint.h>
#include <stdbool.h>
#include "asset_tracker.h"

#include <zephyr/kernel.h>

#include "peripherals/at_led.h"
#include "peripherals/at_button.h"
#include "peripherals/at_sht41.h"
#include "peripherals/at_lis3dh.h"
#include "peripherals/at_usb.h"

#include <zephyr/logging/log.h>
LOG_MODULE_REGISTER(main, CONFIG_TRACKER_LOG_LEVEL);

/* 1000 msec = 1 sec */
#define LED_FAST   100
#define LED_WS	   2000
#define LED_SLOW   1000

int main(void)
{
	init_at_usb();
	#if defined(CONFIG_LOG)
	k_msleep(2000);
	#endif

	init_at_led();
	init_at_button();
	init_at_sht41();
	init_at_lis3dh();

	LOG_INF("Starting Device...");

	if (at_thread_init()) {
		LOG_ERR("Failed to start at thread");
	}

	while (1) {
		if(button_long_press == true) { 
			at_led_toggle();
			k_msleep(LED_FAST);
		} else {
			at_led_toggle();
			k_msleep(LED_SLOW);
		}
	}
	return 0;
}
