# 兰顿蚂蚁

以下内容[摘自维基百科](https://zh.wikipedia.org/wiki/%E5%85%B0%E9%A1%BF%E8%9A%82%E8%9A%81)

> 兰顿蚂蚁（英语：Langton's ant）是细胞自动机的例子。它由克里斯托夫·兰顿在1986年提出，它由黑白格子和一只“蚂蚁”构成，是一个二维图灵机。兰顿蚂蚁拥有非常简单的逻辑和复杂的表现。在2000年兰顿蚂蚁的图灵完备性被证明。兰顿蚂蚁的想法后来被推广，比如使用多种颜色。
>
> ## 规则
>
> 在平面上的正方形格被填上黑色或白色。在其中一格正方形有一只“蚂蚁”。它的头部朝向上下左右其中一方。
>
>* 若蚂蚁在白格，右转90度，将该格改为黑格，向前移一步；
>* 若蚂蚁在黑格，左转90度，将该格改为白格，向前移一步。

## D3JS版本

* [Online Demo](//blog.zhangyu.so/langtonant/index.html)

## Golang终端版本

* go build -o lant lant.go
* ./lant -i 200 # 默认间隔100ms

建议深色终端背景，终端能支持几个方向箭头字符，即
```bash
echo "\u25b2"  # ▲
echo "\u25bc"  # ▼
echo "\u25c0"  # ◀
echo "\u25b6"  # ▶
```

执行效果：

[![asciicast](https://asciinema.org/a/2qrizlvQHwE6ynwlxdrC8iOMY.svg)](https://asciinema.org/a/2qrizlvQHwE6ynwlxdrC8iOMY)
