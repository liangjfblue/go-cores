/**
 *
 * @author liangjf
 * @create on 2020/9/2
 * @version 1.0
 */
package cgo

/*
#cgo CFLAGS: -I ./include
#cgo LDFLAGS: -L./lib/64 -lcal

int append(char *a, char *b, char *ret)
{
    return append(a, b, ret);
}
*/
import "C"
