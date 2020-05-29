package hl


type M struct {
	g0 *g // go运行时系统启动自动创建 执行运行时任务
	mstartfn func() //M 起始函数
	curg *g  // 当前运行时的G的指针
	p punitptr // 当前M关联的那个P
	nexp punitptr // 当前M潜在关联的P
	spinning bool //这个M是否在寻找可运行的G
	lockedg *g //具体M锁定的哪个G
}

// go运行时系统把一个M和G锁定在一起 一旦锁定 M只能运行该G
// 标准库中 LockOSThread UnLockOSThread  锁定和解锁的具体方式

//P 是G能够运行在M中的关键 运行时系统会让P与不同的M 建立关系或者断开链接，使P中的可运行的G队列及时获取运行的时机

//修改P的最大值
// 1 runtime.GOMAXPROC
// 2 环境变量获取
// 其实这个最大数量是对并发G的规模的一种限制 P的数量也就是可运行G的队列的数量
// G被启用后，追加到某P的可运行G队列中 等待调用 只有当一个P和M关联在一起的时候 才会有机会运行
/*
确定P的最大数量


P的空闲列表
	当一个P不再与一个M关联的时候 系统将其放入该列表  条件为 该P下的G队列为空
	当需要一个空闲的P关联某M的时候 从该列表中取出一个

P是有状态的
Pidle 当前P并未与任何M关联
PRunning 关联中
Psyscall P中运行的G在进行系统调用
Pgcstop 运行时系统需要停止调度 例如gc
Pdead 不会再被使用  例如 系统设置减少了P的数量
P中 可运行队列和自由队列 运行完的G 会丢入P的自由队列中  需要的时候 从自由队列获取  除非是不够 才会创建 提高服用氯
*/