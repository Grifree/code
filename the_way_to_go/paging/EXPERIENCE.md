# 总结(初稿20200417)

在GO版分页实现和沟通的过程中，领会到一些有用的经验，记录在此加深印象，后续亦可温故而知新。
总结内容由三方面组成：     
1. `技能点`：一些可以实现更好代码的技巧经验        
2. `问题`：过程中遇到的部分思考而未解的问题，待后续解惑    
3. `进阶体会`：罗列一些有益于找寻`技能点`的方法

## 技能点

- 掌握对（需求、功能、实现目的）的理解思路，即分析问题的能力
- 尽可能避免实现者或使用者出错的风险，是一种好的设计技巧

### panic & error & 修正

> 问题：分页传入参数`page`是负数，应该怎么做？     
(A)`panic` (B)return `error` (C)将`page`值修正为1

1. 首先，区别panic和error的意义        
error：目的是让上层捕获错误进行针对性的处理，处理通常包含需要增加日志。但会遇到错误层层传递情况，是否有这必要性？      
panic：绝对性错误，没有可处理的必要，或没有合适的处理方式     
那么，page为负数时，则应该选择panic
2. **思考业务场景**，什么时候需要主动修正page值为1？        
当请求参数缺失的情况下，例如首次打开某列表页面，没有请求参数page。这个现象很**常见**，并且有**出现合理**性，那么此时，page值会为0值，可以修正其值为1。       

> 每页显示数量`PageSize`为0值，应该报错？还是修正为一个默认值，例如10？
       
应该报错：不传PageSize是不应该的，它会影响计算分页的总页数。
不该修正默认值：原因可看后续小节`拒绝默认值`。


### 常量

尽量使用常量是个好习惯
```go
var errPageCanNotLessZero = errors.New("page can not less zero")
if Page <= 0 {
	panic(errPageCanNotLessZero)
}
```
// TODO 说明原因

每一个都应该明确定义
```go
const errMsgPrefix = "{packageName}: paging.CreatePagingData(gen) "
var errPageSizeCannotLessZero = errors.New(errMsgPrefix+"pageSize cannot less zero")
var errPageSizeCannotBeZero = errors.New(errMsgPrefix+"pageSize cannot be 0")
```

### 减法设计
> 在设计的时候，是实现更多的功能，还是明确精炼的功能？

在前端UI组件设计的时候，经常会让一个组件满足更多的情况，实现更多的功能。这不乏是前端UI需要面对多设备多终端多交互的原因成分，但更合理的是进行拆解，每一个是独立精简且明确的，通过再封装或组合等方式去满足更多的需求。        
比如一个接口，传参不传参是两种不同功能，传一个参两个参又是不同的功能。这会导致一个情况，每次使用时必须看着文档，还有实现与文档不同步的情况。      
但Go语言的出现与其业务场景，需要的更多的是明确清晰的目标。

> 那么，以下4个分页的输入参数，选择哪个？

（1）
```go
type Gen struct {
	Page int // 请求页数
	Total int // 总数
	PageSize int // 每页显示数量
	PageCount int // 总页数
}
```
（2）
```go
type Gen struct {
	Page int
	Total int
	PageSize int
}
```
（3）
```go
type Gen struct {
	Page int
	PageCount int
}
```
（4）
```go
type Gen struct {
	Page int
	Total int
	PageCount int
}
```
`Total` `PageSize`和`PageCount`在分页中的作用都是为了知道分页的总页数。      
如果选3，分页需要显示总数的时候怎么办？`Page * PageCount`可以计算一个总数的近似值，但不好。     
溯流追源，实际数据来源过程，先有`Total` `PageSize`，才有`PageCount`，而且前两者也可以计算出`PageCount`，那么它是否可省略？可以！

**这就是设计的减法思想**

那1不好的地方在于，使用者可以出现234的部分数据输入，这会使分页总页数的计算依据是哪个数据不明确，也会导致使用者可能出现不清晰的代码。这种设计是不够好的。

### 拒绝默认值
> 默认值和零值的区别？

`PageSize int`没有赋值时，会有个初始0值。        
但在客户端向服务器端请求数据时，客户端可能不传`PageSize`，服务器会有个默认每页数量为10个的情况，那10就是`PageSize`的默认值。      
这两种情况下，     
零值：可理解为数据默认值、数据初始值、零值。      
默认值：可理解为逻辑默认值。

> 基于此概念，再回头看分页的实现中，传入参数`PageSize`为0时，适合将其修正为某一个默认值吗？

// TODO 可以在补充说明不合适的原因

### 数据逻辑 & 渲染逻辑

> [数据](https://github.com/onface/paging/blob/master/README.md) 中分为两部分，大致罗列如下，它们有什么区别？

（1）
```
page pageCount pageSize dataTotal prevBatch nextBatch prevSomePage nextSomePage
```
（2）
```
hasPaging isFirstPage isLastPage prevBatch nextBatch prevPage nextPage prevHasMorePage nextHasMorePage prevSomePage nextSomePage
```
上部分数据是用作进行数据逻辑的，根据1中的数据进行逻辑计算，获得第2部分渲染分页所需的数据。      
下部分数据完全是分页渲染时，用于控制具体对应部分的显示或渲染判断。

这个区分性质的概念很重要，它会在设计接口或约定数据时，使最终实现更好。

> 在分页中实现`url`，应不应该？好不好？

不好不应该！     
  
1. 在js版中url是必须实现的吗？      
我认为是可以提供便捷的获取跳转链接url方法, 但不是必须提供的。因为这个方法是辅助渲染用的，不是渲染不可缺。      
2. go版中是不应该实现的吗？     
go和js不同在于，js更偏向展示出更多组合样式的分页，而go在于明确给出分页所需数据。两版分页实现的目的是不同的。


### 明确清晰的数据结构

数据的命名、归类，都体现了实现者的思想认知！

1、名字的定义

> 对于输入数据，当前页的相邻页码显示数量和点击"..."显示更多向前/后页码的码数，如下两种设计，哪种更好？

（A）
```go
type Gen struct {
	...
	PrevSomePage int // 点击"..."，需显示更多向前的哪一页
	PrevBatchPage int // 显示当前页向前相邻的几个页码
	NextSomePage int // 点击"..."，需显示更多向后的哪一页
	NextBatchPage int // 显示当前页向后相邻的几个页码
}
```
（B）
```go
type Gen struct {
	...
	PrevJumpBatchPage int // 点击"..."，需显示跳至向前的哪一页
	PrevClosestPagesCount int // 显示当前页向前相邻的几个页码
	NextJumpBatchPage int // 点击"..."，需显示跳至向后的哪一页
	NextClosestPagesCount int // 显示当前页向后相邻的几个页码
}
```

2、数据结构归类

> 同样的问题，上面B和下面C，两种设计，哪种更好？

（C）
```go
type Gen struct {
	...
	PrevPage struct {
		ClosestCount int
		JumpPage int
	}
	NextPage struct {
		ClosestCount int
		JumpPage int
	}
}
```
C优于B的地方在于，3有同类归纳的意识。        

3、识别度

C不好的地方在于`PrevPage`或`NextPage`下的属性时，因为命名一样，`PrevPage.JumpPage` `NextPage.JumpPage`很容易疏忽导致混淆。     
尽可能调用属性时，也提高命名差异，降低出错风险。

> 那将C修改成下面D，怎么样？

（D）
```go
type Gen struct {
	...
	ClosestPages struct {
		PrevCount int
		NextCount int
	}
	JumpBatchPage struct {
		PrevPage int
		NextPage int
	}
}
```

！！到D，很多时候就结束了，但是是不是应该提个问题`还有没有更好的`？这个提问的想法意识很重要
     
再看看第E种

（E）
```go
type Gen struct {
	...
	ClosestPages struct {
		PrevLength int
		NextLength int
	}
	JumpBatchPage struct {
		PrevInterval int
		NextInterval int
	}
}
```


### 避免重复数据

在分页一开始实现时，我定义的输出数据大致是这样的

（1）
```go
type PagingData struct {
	Page int
	PageSize int
	Total int

	HasPaging bool
	LastPage int
	IsFirstPage bool
	IsLastPage bool
	ClosestPages struct{
		Prev []int
		Next []int
	}
	JumpBatchPage struct{
		HasPrev bool
		PrevPage int
		HasNext bool
		NextPage int
	}
}
```
在自我意识到`数据逻辑 & 渲染逻辑`的区别后，我大致改成了这样

（2）
```go
type PagingData struct {
	Page int
	PageSize int
	Total int

	RenderData struct {
		HasPaging bool
		LastPage int
		IsFirstPage bool
		IsLastPage bool
		ClosestPages struct{
			Prev []int
			Next []int
		}
		JumpBatchPage struct{
			HasPrev bool
			PrevPage int
			HasNext bool
			NextPage int
		}
	}
}
```
可是，在实现逻辑过程中，遇到了一个场景如下所示，我修改了page值，但后面我需要判断page值做一些操作。
> 那么请问（？）处，应该写`Gen.Page`还是`PagingData.Page`？

```go
// 输入参数 Gen
if Gen.Page == 0 {
	PagingData.Page = 1
}

// do very very many things 

if ?.page == 1 {
	// do some thing
}
```
这里的问题是，一开始（？）处可能很明确写了`PagingData.Page`，但是如果有其他人来维护这段代码呢？
他会没有注意到很多行代码之前page值被修改过，而两个page都可被直接使用时，他是不是两个page都有可能使用在后面的代码中？        
所以这里有一个**隐患**。

隐患的原因是出现了两个page，它是重复数据。只要使只有一个page可以使用即可。       
在第2种的基础上，将输出参数改成了如下
```go
func CreateData(gen Gen) ( RenderData, Gen) {}

type Gen struct {
	Page int
	PageSize int
	Total int

	ClosestPages struct {
		PrevLength int
		NextLength int
	}
	JumpBatchPage struct {
		PrevInterval int
		NextInterval int
	}
}
type RenderData struct {
	HasPaging bool
	LastPage int
	IsFirstPage bool
	IsLastPage bool
	ClosestPages struct{
		Prev []int
		Next []int
	}
	JumpBatchPage struct{
		HasPrev bool
		PrevPage int
		HasNext bool
		NextPage int
	}
}
```

在一开始我没有意识到重复数据的情况，是因为在js中Object对象作为参数，它是传址的，存在被修改的安全隐患，数据不可信，导致开始的设计（2）中PagingData是内部一个独立的完整数据。     
js语言环境导致我设计的限制。这一点中可知，go的多返回数据和显性指针的优点，它可以更明确更清晰。


### 统一数据，减少转换

通篇使用了的每页显示数量是`PageSize`，但服务端语言大多官方支持的是`PerPage`，那么尽量前后端都同时使用`PerPage`。


### 杜绝祈祷式编程

分析本质规律，写算法，写逻辑代码。不依靠一步步测试走通为目的。

尤其在多层if，循环语句用常见。

// TODO 补充示例说明


### 单元测试

1. 编写时间     
测试代码应与逻辑代码实现，同步进行。每写一段逻辑代码，就应当对这段代码进行对应单元测试代码的编写，以尽早发现逻辑错误或疏漏，也避免单元测试代码未尽可能多地覆盖逻辑分支和值。
2. 覆盖率      
语句覆盖率100% < 逻辑分支覆盖率100% < 测试100%
3. defer recover捕捉panic

## 遗留问题

1. math.ceil 为什么支持float 不支持int

2. js实现版的数据分析思路, 是如何设计提炼出来的
[数据](https://github.com/onface/paging/blob/master/README.md)

## 进阶

### 总结 & 输出 
灵感总是一闪而过，领悟时常随风而散。      
自己写下来的积累更持久，脑子记不住全部的东西。     
经过自己吸收，再输出的东西，才是经过消化和理解的。写是输出的一种方式。

### 阶梯式的学习
在学习和教学的过程中，应该一步步深入，达成完整的学习路线。例如：
1. 完成基础的实现
2. 完善实现
3. 优化实现
4. 摸索场景边界、跳出思考
5. 提炼、精简
6. 最优的设计

### 思索的过程
对于编程生涯中，如何不断提升自己的编码能力、技巧和意识，可以层层深入的思考这样的问题。
1. 这样写对吗？
2. 这样写好吗？
3. 还有别的写法吗？
4. 不同的写法，本质区别是什么？
5. 哪种更好？为什么？
6. 最后选择哪种？

当发现问题后，分两步依次分析问题。
1. 解决这个问题
    - 问题的原因是什么
    - 怎么意识到问题的？怎么发现有问题的？
    - 问题的解决办法是什么？
    - 还有别的解决方式吗？
    - 哪种解决办法更好？
2. 解决同类问题
    - 问题的本质是什么？
    - 有没有同类的问题？
    - 如何解决它们？以防再发生

```
例如：之前`重复数据`的问题中，也是这样的思考过程。        
Q：隐患的本质是什么？     
A：原因是出现了两个page，它是重复数据       
Q：如何避免出现重复数据？
A：只要使只有一个page可以使用即可。        
Q：如何保证只有一个page可使用方法？        
A：需要想并试验各种不同的方法。// TODO 补充更多过程示例        
Q：多种方法对比，哪种方法更好？        
Q：还有没有更好的？      
Q：怎么意识这个重复数据的隐患的？       
A：起源意识很重要！有不明确的可能就是隐患风险。
```