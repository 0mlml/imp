#ifndef AT_LIS3DHTR_H
#define AT_LIS3DHTR_H

int init_at_lis3dh(void);
int get_accel(struct at_sensors *sensors);

#endif /* AT_LIS3DHTR_H */