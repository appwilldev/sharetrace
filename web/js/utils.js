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


function isEmpty(obj) {
    // null and undefined are "empty"
    if (obj == null) return true;

    // Assume if it has a length property with a non-zero value
    // that that property is correct.
    if (obj.length > 0)    return false;
    if (obj.length === 0)  return true;

    // Otherwise, does it have any properties of its own?
    // Note that this doesn't handle
    // toString and valueOf enumeration bugs in IE < 9
    for (var key in obj) {
        if (hasOwnProperty.call(obj, key)) return false;
    }

    return true;
}


$.extend({                                                                              getUrlVars: function(){
        var vars = [], hash;
        var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&'); 
        for(var i = 0; i < hashes.length; i++){
            hash = hashes[i].split('=');
            //vars.push(hash[0]); 
            vars[hash[0]] = hash[1];
        } 
        return vars;
    },
}); 
