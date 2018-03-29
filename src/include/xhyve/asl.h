#pragma once

/* Initialize ASL logger and local buffer. */
void asl_init(void);

/* Send one character to the logger: wait for full lines before actually sending. */
void asl_put(uint8_t _c);
