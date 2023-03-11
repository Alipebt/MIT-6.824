# 笔记

## 1 MapReaduce

### 1.1 思路 

`mr/worker.go`中的`Worker()`接收了`Map`和`Readuce`函数，然后通过`Call()`来向协调器发送`RPC` 请求任务，获取`filenames`。然后在`Worker()`中处理`filenames`，将每个文件都打开读取其中内容，并且生成键值对。将所有键值对添加到中间件`intermediate`(也是键值对)中。此时不可以将生成的键值对作为整体传入`intermediate`，所以需要用到`...`来依次添加。由于`append()`返回的是一个没有经过显示定义的切片(添加的切片类型可能与原类型不同)，所以在之后的操作应该将`intermediate`类型装换为键值对。