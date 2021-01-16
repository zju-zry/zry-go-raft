/**
 * @note: event对象，在server接收到消息的时候使用这个event将信息传递到循环中
 * @author: zhangruiyuan
 * @date:2021/1/16
**/
package common

type Ev struct {
	// 这个内容存放传入的信息
	Target       interface{}
	// 这个内容存放经过处理后的信息，
	// 并设置这个数据为chan类型，
	// 这个数据由循环进行处理，在处理结束之后，
	// 由server端获取相应的内容，并进行返回
	ReturnC      chan bool
	ReturnValue  interface{}
	/**
	 * 启动的server服务中可能由于收到错误的服务申请
	 * 而产生错误，这个错误可以通过这个通道进行传递
	 * 如在follow状态下收到了其他状态的信息
	 */
	C           chan error
}
