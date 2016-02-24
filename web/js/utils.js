//复制数据
function clone(cond) {
    var copy;
    if (cond instanceof Array) {
        copy = [];  //创建一个空的数组 
        var i = cond.length;
        while (i--) {
            copy[i] = clone(cond[i]);
        }
        return copy;
    } else if (cond instanceof Object) {
        copy = {};  //创建一个空对象 
        for (var k in cond) {  //为这个对象添加新的属性 
            copy[k] = clone(cond[k]);
        }
        return copy;
    } else {
        return cond;
    }
}