package main

/**
map的元素 不是一个变量 不可以获取其地址
例如 _ = &ages["a"]  由于map增长可能会导致已有的元素被重新散列到新的位置
无序性
排序 sort包 对于key进行排序后 getkey进行输出

map 与slice 不可以进行比较 唯一能做比较 是做非空比较 如果要做比较 需要进行循环
*/

//map 比较
func compareMap(a,b map[string]string)bool{
	if len(a)!=len(b){
		return false
	}

	for k,v:=range a{
		 if vb,ok:=b[k];!ok|| vb!=v{
		 	return false
		 }
	}
	return true
}
//实现go版本set
