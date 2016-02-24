var LOGIC_FUNCS = {
    and : "且",
    or  : "或",
    not : "非"
};

//func对应的显示名
var LEAF_FUNC_TITLES = {
    "="                     : "等于",
    "!="                    : "不等于",
    "\u003c"                : "小于",
    "\u003c="               : "小于等于",
    "\u003e"                : "大于",
    "\u003e="               : "大于等于",
    "ver="                  : "版本等于",
    "ver!="                 : "版本不等于",
    "ver\u003c"             : "版本小于",
    "ver\u003c="            : "版本小于等于",
    "ver\u003e"             : "版本大于",
    "ver\u003e="            : "版本大于等于",
    "str=?"                 : "str等于",
    "str!=?"                : "str不等于",
    "str-contains?"         : "str包含",
    "str-not-contains?"     : "str不包含",
    "wildcard"              : "wildcard匹配",
    "wildcard-not"          : "wildcard不匹配",
    "str-empty?"            : "str为空",
    "str-not-empty?"        : "str不为空"
};

var SYMBOL_TABLE = [
    ["APP_VERSION", "App版本号",      ["ver=","ver!=","ver\u003c","ver\u003c=","ver\u003e","ver\u003e="],                  "\d+(\.\d+)*"],
    ["LANG",        "系统语言",       ["str=?","str!=?","str-contains?","str-not-contains?","wildcard","wildcard-not"],    "^[a-z]+$"],
    ["NETWORK",     "网络",           ["str=?","str!=?","str-contains?","str-not-contains?","wildcard","wildcard-not"],    "^[A-Za-z0-9]+$"],
    ["OS_VERSION",  "操作系统版本号",  ["ver=","ver!=","ver\u003c","ver\u003c=","ver\u003e","ver\u003e="],                   "\d+(\.\d+)*"],
    ["TIMEZONE",    "时区",           ["str=?","str!=?","str-contains?","str-not-contains?","wildcard","wildcard-not"],    "^-?H(([0-1]\d)|(2[0-4])):[0-5]\d$"],
    ["IP_ADDRESS",  "IP地址",         ["str=?","str!=?","str-contains?","str-not-contains?","wildcard","wildcard-not"],    "((([1-9]?|1\d)\d|2([0-4]\d|5[0-5]))\.){3}(([1-9]?|1\d)\d|2([0-4]\d|5[0-5]))"],
    ["PROVINCE",    "省/市/自治区",    ["str=?","str!=?","str-contains?","str-not-contains?","wildcard","wildcard-not"],    "^[\u4e00-\u9fa5]*$"],
    ["CITY",        "城市",           ["str=?","str!=?","str-contains?","str-not-contains?","wildcard","wildcard-not"],    "^[\u4e00-\u9fa5]*$"],
    ["OPERATOR",    "运营商",         ["str=?","str!=?","str-contains?","str-not-contains?","wildcard","wildcard-not"],     "^[a-zA-Z0-9\u4e00-\u9fa5]+$"]
];

var SYMBOLS = [];           //Symbol列表
var SYMBOL_TITLES = {};     //Symbol对应的显示名
var SYMBOL_FUNCS  = {};
var SYMBOL_FOCUS  = {};     //Symbol可用的func列表

for(idx in SYMBOL_TABLE) {
    SYMBOLS.push(SYMBOL_TABLE[idx][0]);
    SYMBOL_TITLES[SYMBOL_TABLE[idx][0]] = SYMBOL_TABLE[idx][1];
    SYMBOL_FUNCS[SYMBOL_TABLE[idx][0]]  = SYMBOL_TABLE[idx][2];
    SYMBOL_FOCUS[SYMBOL_TABLE[idx][0]]  = SYMBOL_TABLE[idx][3];
}