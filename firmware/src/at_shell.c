#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

#include <zephyr/kernel.h>
#include <zephyr/shell/shell.h>
#include "at_shell.h"

#include <zephyr/logging/log.h>
LOG_MODULE_REGISTER(at_shell, CONFIG_TRACKER_LOG_LEVEL);

#define AT_STACK_SIZE 5120
#define AT_WORKER_PRIO 5

struct k_work_q at_work_q;
K_THREAD_STACK_DEFINE(at_work_q_stack, AT_STACK_SIZE);

static at_ctx_t *atcontext;

DECLARE_RESULT(size_t, StrToHexRet, MAC_INVALID);
DECLARE_RESULT(uint8_t, HexToByteRet, NOT_HEX);
DECLARE_RESULT(uint8_t, StrToDecRet, OUT_RANGE);

void AT_CLI_init(at_ctx_t *at_context)
{
	atcontext = at_context;

	k_work_queue_init(&at_work_q);

	k_work_queue_start(&at_work_q, at_work_q_stack,
			   K_THREAD_STACK_SIZEOF(at_work_q_stack), AT_WORKER_PRIO, NULL);
}
