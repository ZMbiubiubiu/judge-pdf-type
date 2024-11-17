
# 扫描版pdf

* 什么是扫描版pdf？

扫描版pdf，即scanned pdf
> A text document which is scanned to PDF without the text being recognised by optical character recognition (OCR) is an image, with no fonts or text properties.

扫描版pdf每一页都是图片，没有文字对象。

因为每一页都是图片，所以无法搜索、编辑文本。另外，扫描版pdf的文件大小通常会比普通pdf大，尤其是图片的resolution很高的时候。

* 为什么需要扫描版pdf？

它的目的就是保持原样，扫描版pdf通常是纸质文档的电子版，比如合同、书籍、论文等。

# 普通PDF与扫描版PDF的区别在于

扫描版PDF的每一页都是图片对象，而普通PDF的每一页几乎都很有很多文字对象。

所以我们这里的实现较为简单：

* 对于一个pdf文档，最多随机选取n页内容
* 根据每一页的内容判断是否为普通页
  * 普通页的定义：text_fragment的数量 >= 5
  * 普通页的数量为m
* 如果普通页占比选取总页数$\frac{m}{n}$>=80%$，那么推定该文档为普通pdf，否则为扫描版pdf

# 参考

* [quora - What is a scanned PDF?](https://www.quora.com/What-is-a-scan-PDF)
* [wiki - PDF](https://en.wikipedia.org/wiki/PDF)
