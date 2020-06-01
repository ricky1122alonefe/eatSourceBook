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
P中 可运行队列和自由队列 运行完的G 会丢入P的自由队列中  需要的时候 从自由队列获取  除非是不够 才会创建 提高复用
*/

/*
有一个全局的G列表 runtime.allgs 新建的时候第一时间加入全局列表 集中存放当前运行时系统所有G的指针
初始化
初始化之后G被储存到本地的P runnext 字段中 存放新的G以求更早运行 如果当前P的runnext 已经有一个G 那么这个新建的G会被移动到该P 可运行队列的末尾
如果队列满， G只能追加到可运行队列中
G的状态
Gidle 新分配 但是没有进行初始化
Grunnable 可运行队列中等待运行
Grunning 正在运行
Gsyscall 执行系统调用
Gwaiting 阻塞
Gdead 闲置
Gcopystack G的栈被移动 栈的扩展或者收缩
Pdead Gdead 不同 Pdead 只能面临回收的结果 Gdead 会放入本地P或者全局自由列表，重用
 */


 /*
 PMG的容器
 全局M列表 runtime.allm  运行时系统 存放所有M的一个单向列表
 全局P列表 runtime.allp	运行时系统 存放所有P的一个数组
 全局G列表 runtime.allg	运行时系统 存放所有G的一个切片
 	任何G都会存在于全局G列表中

 调度器空间的M列表 runtime.sched.midle 调度器	空闲的M的单项连标
 调度器空间的P列表 runtime.sched.pidle 调度器 空闲的P的单项列表
 调度器可运行G队列 runtime.sched.runqhead runtime.sched.runtail 可运行G的队列 头尾

 调度器G自由列表 runtime.sched.gfreestack runtime.sched.gfreenostack  单项列表

 P可运行的G队列 runtime.p.runq  本地P 可运行的G的一个队列
 	Gdead之后优先放入本地的P的自由列表
 P的自由G列表 runtime.p.gfree 当前P中自由G的单项连标
 	如果本地的自由列表空了 运行时系统会先从调度器的自由列表转移一部分G到其中
 	当本地的自由G列表满了 会将本地的自由G列表转移给调度器的G列表
 */

/*
查找可运行的G
第一阶段
		1runtime.SetFinalizer  一个专用的G完成任务之后获取它并且状态设置为runnable 放到本地P可运行队列中
		2从本地的P的可运行G队列获取G
		3从调度器可运行G队列获取G
		4从IO netpooler 获取G
		5从其他P的可运行队列获取G 伪随机算法将全局P列表中的一个P的可运行队列的一半 转移到当前可运行P列表中
		规则1 除了本地的P还需要有非空闲的P 空闲的P的可运行的G队列必定为空
		规则2 当前M处于spinning状态 或者处于spinning的M的数量小于非空P的二分之一  spinning 意味着没有找到G来运行
		其中选取和转移会有很多次 成功则停止 如果成功会把获取的第一个G返回 否则 第一阶段结束
第二阶段
		6获取执行gc标记任务的G 正处于gc标记阶段 和本地P是否可用于GC标记任务 都为true  调度器会将本地P持有的ghc标记的专用G设置为runnable状态 返回
		7从调度器的可运行G的队列获取G 如果找不到可运行的G 则解除本地P与M的关联 并且将该P放入调度器的空闲P列表
			P的空闲列表
			当一个P不再与一个M关联的时候 系统将其放入该列表  条件为 该P下的G队列为空
			当需要一个空闲的P关联某M的时候 从该列表中取出一个
		8从全局P列表每个P的可运行G队列获取G
			遍历全局P列表的P,查找可运行的G队列 只要发现某P的可运行队列G非空,从P空闲列表里面取出一个P，判定可用后与当前M关联在一起 返回地有阶段重新搜索的G
		9 获取gc标记任务的G 判断是否处于gc阶段 以及与gc标记任务相关的全局资源是否可用
			空闲列表取一个P 如果这个p持有一个gc专用G 关联当前M P 执行6
*/