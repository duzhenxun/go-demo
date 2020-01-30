//➜  生成文档 git:(master) ✗ go doc
//➜  生成文档 git:(master) ✗ go doc Demo
//➜  生成文档 git:(master) ✗ go doc Demo.Test1

package demo
//这里的Demo的注释
//再注释一行
type Demo struct {

}
//测试1 注释
func (t *Demo) Test1(v int) int {
	return v
}

// 测试2 注释
// 写好注释
func (t *Demo) Test2() int  {
	return 999
}
// 专业写法 godoc使用
//		e.g. t.Test3(123)
func(t *Demo) Test3(v int)int{
	return v
}
