#ifndef AT_UPLINK_H
#define AT_UPLINK_H

#include <asset_tracker.h>

void at_send_uplink(at_ctx_t *context);
void at_msg_sent(at_ctx_t *context);
void at_send_error(at_ctx_t *context);

#endif // AT_UPLINK_H