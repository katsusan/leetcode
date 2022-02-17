## 1. JavaScript基础

### 1.1 变量
var: 旧式的变量声明, 允许变量提升(x = 1; var x;)以及重复定义(var x = 3; var x = 4;).
let: 新式的创建变量关键字,不允许变量提升和重复定义.

类型:
- Number 
  + x.toFixed(n): 四舍五入保留n位小数转化为字符串
  + x.toString(): 转化为字符串
  + x.toExpoential(): 转化为指数表示法
  + x.toPrecision(n): 保留n位精度,包括整数部分和小数部分
- String
  + s1.indexOf(s2)
  + s1.slice(m[, n])
  + s1.replace("src", "dst")
  + s1.toLowerCase()/toUpperCase()
- Boolean
- Array
  + arr.slice(): return a shallow copy
  + arr.keys()/values(): return an iterator of keys/values. 
    * keys(): arr = ["a", , "c"]; arr.name = "rob"; Object.keys(arr) = ['0', '2', 'name']; [...arr.keys()] = [0, 1, 2];
  + arr.entries(): return an iterable object, which can be used in for...of
  + arr.concat(arr2)/split(',')/join('.')
  + arr.toString(): => arr.join(',')
  + arr.reverse()/sort([(a, b) => a - b])
  + arr.indexOf(e)/lastIndexOf(e): first/last
  + arr.push()/pop()/shift()/unshift(): 尾部插入/尾部删除/头部删除/头部插入
  + arr.splice(i, n): 删除n个元素,从i开始
  + arr.includes(e): 判断是否包含
  + arr.find(fn)/findIndex(fn): fn is `(e) = e > 0`
  + arr.flat([depth]): 展开数组, [1, 5, 67, [2, 3, 4]] => [1, 5, 67, 2, 3, 4]
  + arr.forEach(elem => console.log(elem)): for every elem, do something
  + arr.map(elem => elem * 2): for every elem, do something and returns a new array
  + arr.filter( elem => elem > 0): return a new array with elements that satisfy the condition
  + arr.every(elem => elem > 0): return true if every element satisfies the condition
  + arr.some(elem => elem > 0): return true if some element satisfies the condition
  + arr.reduce(callback[, initvalue]): [1, 5, 7].reduce((a, b) => a + b) => 13; reduceRight like reduce, starts with the last element
  + [Advanced] Array.prototype.forEach.call([1, 5, 7], (elem) => console.log(elem))
- Object
- 


`===` vs `==` | `!==` vs `!=`: 后者测试值是否等同,但类型可能不同. 但前者严格测试类型和值是否相等.   
尽量用`===`/`!==`来进行更严格的比较.


### 1.2 控制流

- if/else/else if
- switch (v) {case v1: ...; break; case v2: ...; break; default: ...;}
- for (;;) {continue/break [label]}
- for (v in object) {} // iterates over all the enumerable properties of an object
- for (v of iterable) {} // iterates over iterable objects(Array, Map, Set, arguments)
- while(cond) {}
- do {} while (cond);

### 1.3 函数

- generator function: function * f(), e.g. const function* fn() {yield 1; yield 2; yield 3;}; for {const v of fn()} {console.log(v);}
- arrow function: () => {}
- Function constructor: new Function('a', 'b', 'return a + b');   // not recommended, has security risks.

### 1.4 事件

事件对象: 传给事件处理函数的对象, 包含了事件的名称, 事件发生的时间, 事件发生的元素, 以及其他信息.

比如:
```javascript
button.onclick = function(e) {
  e.target.style.backgroundColor = 'red';
  console.log(e.target);    // 打印button的DOM元素
  console.log(e.type);      // 打印事件的类型, e.g. 点击时为'click'
  console.log(e.timeStamp); // 打印事件发生的时间戳, time (in milliseconds) at which the event was created
}
```



