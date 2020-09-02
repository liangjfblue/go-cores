#ifndef _CAL_H_
#define _CAL_H_

/*相加数字的范围*/
#define CAL_MAX                                1<<6
#define CAL_MIN                                0

int append(char *a, char *b, char *sum);

#ifdef _WIN32
typedef __int64 int64_t;
typedef unsigned __int64 uint64_t;
#else
#include <stdint.h>
#endif


#endif /* _CAL_H_ */
