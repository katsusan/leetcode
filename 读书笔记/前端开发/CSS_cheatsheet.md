# 1. layout

## 1.1 flexbox
  refer: https://developer.mozilla.org/zh-CN/docs/Learn/CSS/CSS_layout/Flexbox

### 1.1.1 flex container
display: flex;

flex-direction: row/row-reverse/column/column-reverse;

flex-wrap: nowrap/wrap/wrap-reverse or inherit/initial/unset;   // 控制子元素的换行
- wrap：换行
- nowrap：不换行，被摆放到一行
- wrap-reverse：换行，但是换行的方向是反的
- inherit：继承父元素的属性
- initial：设置为默认值
- unset：取消某个属性

align-items: start/end/center/baseline/stretch;   // 将所有直接子节点上的align-self值设置为一个组
justify-content: center/start/end/between/around/stretch/space-around/...;  // 定义如何分配顺着弹性容器主轴(或者网格行轴) 的元素之间及其周围的空间


### 1.1.2 flex item

// flex缩写   
flex: 
- flex-grow: 规定了 flex-grow 项在 flex 容器中分配剩余空间的相对比例。
- flow-shrink: 指定了 flex 元素的收缩规则，flex 元素仅在默认宽度之和大于容器的时候才会发生收缩。(number[>0])
- flex-basis: 指定了 flex 元素在主轴方向上的初始大小。(10em/3px/auto, fill/max-content/min-content/fit-content, content)

单值:   
flex: auto/initial/none/inherit/unset;
flex: 2; // flex-grow
flex: 10em (/30px/content);   // flex-basis  

双值：   
flex: 1 30px;   // flex-grow flex-basis
flex: 1 2;      // flex-grow flex-shrink

三值：   
flex: 2 2 10%;  // flex-grow flex-shrink flex-basis

align-self: center/start/end/stretch/...;   // 对齐当前 grid 或 flex 行中的元素，并覆盖已有的 align-items 的值

order：<number>;  // 控制 flex 元素的排序顺序，越小越靠前,默认为0，相同值的按源顺序排列


## 1.2 grid
  refer: https://developer.mozilla.org/zh-CN/docs/Learn/CSS/CSS_layout/Grids

### 1.2.1 grid container

```css
display: grid;
grid-template-columns: 200px 200px 200px;   // 指定列的宽度
grid-template-columns: 2fr 1fr 1fr;     // 按空间比例划分
grid-template-columns: repeat(3, 1fr);  // 等同于gri-template-columns: 1fr 1fr 1fr;
grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));  // 浏览器会在保证200px宽的前提下尽量填充列

grid-template-areas: "a a a"
                     "a b b"
                     "a b b";   // 用命名式方式指定3行3列的网格区域， 可以用'.'留空，名字与网格单元的grid-area属性对应
grid-auto-rows: minmax(100px, auto);    // 行高最低100px，最高自适应，grid-auto-columns类似，但和grid-auto-rows不能同时使用
grid-gap: 20px;     // 指定列之间的间隔
```

显式网格与隐式网格:
- 显式网格：指定了每个网格单元的宽度和高度，并且指定了每个网格单元的内容。
  比如我们用grid-template-columns 或 grid-template-rows 属性创建的网格。
- 隐式网格：指定了每个网格单元的内容，并且没有指定每个网格单元的宽度和高度。
  比如我们用grid-template-areas 属性创建的网格。
  grid-auto-rows和grid-auto-columns用来指定隐式网格单元的高度和宽度。


简单来说，隐式网格就是为了放显式网格放不下的元素，浏览器根据已经定义的显式网格自动生成的网格部分。


### 1.2.2 grid item


```css
div {
  grid-column: <grid-line> (/ <grid-line>); // 由1或2个grid-line指定，或为auto(自动放置，默认1)，或为span(跨越多个grid-line)
                                            // Eg, 1/1、1/2是等同的，1/3代表占据1、2两列. 1/span 3类似于1/4
  grid-row: <grid-line> (/ <grid-line>); // 同上

  grid-area：a;
}
```

## 1.3 float

```CSS
float: left/right/none;   // or inline-start/inline-end or inhert/initial/revert/unset -> global values

clear: none/left/right/both;  // or inline-start/inline-end or inhert/initial/revert/unset -> global values

box-sizing: border-box;   //  tells the browser to account for any border and padding in the values you specify for an element's width and height
            content-box;  // gives you the default CSS box-sizing behavior
```

## 1.4 position

```CSS
position: static;     // 默认值，不改变元素的位置
          absolute;   // 绝对定位，相对于父元素。从正常的文档布局流中移除，并且不受父元素的布局影响；top/bottom/left/right表示偏移。
                      // 定位上下文：containing block，取决于父元素的position属性，若父元素没有明确定义则为默认的static。
                      //            此时绝对定位元素会被包含在初始块容器(有着和浏览器窗口一样的尺寸，并且<html>也被包含在内),规则：
                      //  - static,relative,sticky：formed by the edge of the content box of the nearest ancestor element that is either
                      //         a block container or establishes a formatting context(table container, flex container, ...)
                      //  - absolute: containing block is formed by the edge of the padding box of the nearest ancestor element that has
                      //         a position value other than static.
                      //  - fixed: the containing block is established by the viewport or the page area.
          relative;   // 相对定位，用top/bottom/left/right来改变相对于默认位置的偏移
          fixed;      // 类似absolute，但absolute是相对于最近的定位祖先或<html>元素，fixed是相对于浏览器窗口(因此可以在创建导航栏时使用)
          sticky;     // 类似相对定位与绝对定位的混合，初始表现得像相对定位，滚动到一定的阈值后就变得固定，比如可以使导航栏随页面滚动到特定点然后固定在顶部

z-index: n;   // n为数字，表示层级，层级越高，越在上面，默认值为auto，表示不改变层级
```

## 1.5 multicol(多列布局)

```CSS
column-count: 3;
column-width: 200px;
column-gap: 20px;
column-rule: 2px solid #000;  // combines column-rule-width, column-rule-style, column-rule-color.

// column items
break-inside: avoid;  // 防止分栏内容被分栏分割
```

## 1.6 Reponsive Design(响应式设计)
RWD：Responsive Web Design，响应式网页设计，是一种在网页设计中使用设备尺寸来设计网页的方法，允许Web页面适应不同屏幕宽度因。

### 1.6.1 Media Queries(媒介查询)

example：   
```css
// 当前web页面展示为屏幕媒体(非印刷文档)且视口至少800像素宽
@media screen and (min-width: 800px) {
  .container {
    margin: 1em 2em;
  }
} 
```

使用媒体查询时的一种通用方式是，为窄屏设备（例如移动设备）创建一个简单的单栏布局，然后检查是否是大些的屏幕，   
在你知道你有足够容纳的屏幕宽度的时候，开始采用一种多栏的布局 。这经常被描述为移动优先设计。

### 1.6.2 灵活网格

早年间进行响应式设计的时候，我们唯一的实现布局的选项是使用float。灵活浮动布局是这样实现的，让每个元素都有一个   
作为宽度的百分数，而且确保整个布局的和不会超过100%。

### 1.6.3 现代布局技术

现代布局方式，例如多栏布局、伸缩盒和网格默认是响应式的。

### 1.6.4 响应式图像

max-width -> 对于小尺寸屏幕可能浪费带宽 -> <picture>|<img>的srcset/sizes特性

### 1.6.5 响应式排版

- 媒介查询: @media
- 视口单位：vw/vh/vmin/vmax，1vw等于视口宽度的1%，但这样用的后果是用户失去缩放该文本的能力。
            方法是用calc，比如`font-size: calc(1.5rem+3vw)`

### 1.6.6 视口元标签

```html
<head>
  <meta name="viewport" content="width=device-width,initial-scale=1">
</head>
```

意思是告诉移动端浏览器将视口宽度设为设备宽度，并且缩放比例为1(文档放大到预期比例的100%)。

- initial-scale：设定了页面的初始缩放，我们设定为1。
- height：特别为视口设定一个高度。
- minimum-scale：设定最小缩放级别。
- maximum-scale：设定最大缩放级别。
- user-scalable：如果设为no的话阻止缩放

## 1.7 Media Queries(媒介查询)

```css
@media media-type and (media-feature-rule) {
  /* CSS rules go here */
}
```

- 一个媒体类型，告诉浏览器这段代码是用在什么类型的媒体上的（例如印刷品或者屏幕）；
  + all
  + print
  + screen
  + speech
- 一个媒体表达式，是一个被包含的CSS生效所需的规则或者测试；
  + min-width/max-width/width; width | height;
  + orientation(朝向): portrait(竖放) | landscape(横放);
  + @ media (hover: hover) : 非指点设备，能触发hover事件的设备。
  + @ media (pointer: none/coarse/fine) 
    * none → does not include a pointing device
    * coarse →  includes a pointing device of limited accuracy
    * fine → includes a pointing device of limited accuracy
  
  带逻辑的媒体查询：
    - And:  @ media screen and (min-width: 400px) and (orientation: landscape)
    - Or:   @ media screen and (min-width: 400px), screen and (orientation: landscape)  // separated with comma
    - Not:  @ media not all and (orientation: landscape)  // 只在竖着的时候生效
- 一组CSS规则，会在测试通过且媒体类型正确的时候应用。

## 1.8 Legacy Layout(传统布局)

## 1.9 Supporting Old Browsers

特性查询:   
@supports (display: grid) {
  /* CSS rules go here */
}


# 2 CSS basis

## 2.1 CSS box model(盒模型)

### 2.1.1 盒子分类

- Block box(块级盒子)：
  + 盒子会在内联的方向上扩展并占据父容器在该方向上的所有可用空间，在绝大数情况下意味着盒子会和父容器一样宽
  + 每个盒子都会换行
  + width 和 height 属性可以发挥作用
  + 内边距（padding）, 外边距（margin） 和 边框（border） 会将其他元素从当前盒子周围“推开”
  + example: 除非特殊指定，诸如标题(`<h1>`等)和段落(`<p>`)默认情况下都是块级的盒子。

- Inline box(内联盒子)
  + 盒子不会产生换行。
  + width 和 height 属性将不起作用。
  + 垂直方向的内边距、外边距以及边框会被应用但是不会把其他处于 inline 状态的盒子推开。
  + 水平方向的内边距、外边距以及边框会被应用且会把其他处于 inline 状态的盒子推开。
  + example: 用做链接的 <a> 元素、 <span>、 <em> 以及 <strong> 都是默认处于 inline 状态的。

设置display属性为inline或block可以改变盒子的外部显示类型。

**内部/外部显示类型**

css的box模型有一个外部显示类型，来决定盒子是块级还是内联。

同样盒模型还有内部显示类型，它决定了盒子内部元素是如何布局的。默认情况下是按照 正常文档流 布局，   
也意味着它们和其他块元素以及内联元素一样(如上所述).

在一个元素上，外部显示类型是block，但是内部显示类型修改为display:flex。 该盒子的所有直接子元素都会成为flex元素，   
会根据 弹性盒子（Flexbox ）规则进行布局.

### 2.1.2 盒子组成

- **Content box**: 这个区域是用来显示内容，大小可以通过设置 width 和 height.
- **Padding box**: 包围在内容区域外部的空白区域； 大小通过 padding 相关属性设置。
- **Border box**: 边框盒包裹内容和内边距。大小通过 border 相关属性设置。
- **Margin box**: 这是最外面的区域，是盒子和其他元素之间的空白区域。大小通过 margin 相关属性设置。

```CSS
.box {
  width: 350px;
  height: 150px;
  margin: 25px;
  padding: 25px;
  border: 5px solid black;
}
```

则盒子宽度为350+25x2+5x2=410px, 高度为150+25x2+5x2=210px,


注: margin 不计入实际大小 —— 当然，它会影响盒子在页面所占空间，但是影响的是盒子外部空间。   
盒子的范围到边框为止 —— 不会延伸到margin。


### 2.1.3 替代(IE)盒模型

该模型宽度为可见宽度,内容宽度为width减去padding和border的宽度.  

默认浏览器会使用标准模型。如果需要使用替代模型，可以通过为其设置 box-sizing: border-box.

希望所有元素都使用替代模式:   

```CSS
html {
  box-sizing: border-box;
}
*, *::before, *::after {
  box-sizing: inherit;
}
```

### 2.1.4 外边距margin/内边距padding/边框border

```CSS
margin: margin-top margin-right margin-bottom margin-left;

border, border-top/right/bottom/left: width style color;
border-X-Y: value; // X: top, right, bottom, left, Y: width, style, color

padding, padding-top/bottom/left/right: length;
```

外边距折叠:   
理解外边距的一个关键是外边距折叠的概念。如果你有两个外边距相接的元素，这些外边距将合并为一个外边距，   
即最大的单个外边距的大小。


### 2.1.5 display:inline-block

一个元素使用 display: inline-block，实现我们需要的块级的部分效果(不会切换到新行,但width和height有效)：

- 设置width 和height 属性会生效。
- padding, margin, 以及border 会推开其他元素。





