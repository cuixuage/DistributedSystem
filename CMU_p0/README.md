**1.KV数据库的架构&各项功能**   
![](1.jpg)

**2.kv数据库伪代码**   
![](2.jpg)

**3.总结 interface用来统一接口**  
 ```go
package main                                                                           

import (
    "fmt"
)
//定义接口interface
type Man interface {
    name() string;
    age() int;
}
//**********************接口实现 1
type Woman struct {
}

func (woman Woman) name() string {
   return "Jin Yawei"
}
func (woman Woman) age() int {
   return 23;
}
//*********************接口实现 2
type Men struct {
}

func ( men Men) name() string {
   return "liweibin";
}
func ( men Men) age() int {
    return 27;
}

func main(){
    var man Man;            //接口变量

    man = new(Woman);       //使用women初始化接口
    fmt.Println( man.name());
    fmt.Println( man.age());
    man = new(Men);        //使用men初始化接口
    fmt.Println( man.name());
    fmt.Println( man.age());
}
```
