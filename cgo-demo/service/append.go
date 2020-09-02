/**
 *
 * @author liangjf
 * @create on 2020/9/2
 * @version 1.0
 */
package service

/**
int append(char *a, char *b, char *ret);
*/
import "C"
import "unsafe"

///Append 追加字符串
//export Append
func Append(a, b string) string {
	var ret *C.char = (*C.char)(C.malloc(C.sizeof(a)+C.sizeof(b)) + 1)
	defer C.free(unsafe.Pointer(ret))

	C.append(C.char(a), C.char(b), ret)
	return C.GOString(ret)
}
