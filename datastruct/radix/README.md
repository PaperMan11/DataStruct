## 1 概述
&emsp;&emsp;前缀基数树（radix）又叫基数树，是前缀树的一个变种。它和前缀树不同的地方在于：前缀树是将一个 `string` 按 `char` 进行分段保存，而基数树是将多个 `char` 设为一层，然后将 `string` 进行分层保存，一般利用 **‘/’** 作为分层标识。
## 2 原理
&emsp;&emsp;本次实现的 radix 用做简单的路由匹配，每个结点存储一层路径，终点结点存储了整个路径。
- 每一层的 string 的首字符为 `:` 时为动态匹配，即该层匹配任意的 string
- 每一层的 string 的首字符为 `*` 时为动态匹配，即可匹配后方所有的内容