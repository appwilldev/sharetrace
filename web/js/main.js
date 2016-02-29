//全局变量
var router = null;
var user_auth_ok = false;
var user_info = null;
var COUNT_PER_PAGE = 10;

var init_vue = function () {
    router = new VueRouter({
        //'linkActiveClass': 'active'
        'hashbang':false // TOASK
    });

    router.map({
        '/register': {
            component: Register
        },
        '/login': {
            component: Login
        },
        '/apps/all/:page': {
            component: Apps,
            auth: true
        },
        '/users/all/:page': {
            component: Users,
            auth: true
        },
    });

    router.alias({
        '/apps': '/apps/all/1',
        '/users': '/users/all/1'
    });

    router.beforeEach(function (transition) {
        if (transition.to.auth && !user_auth_ok) {
            transition.abort();
        } else {
            window.scrollTo(0, 0);
            transition.next();
        }
    });

    router.afterEach(function (transition) {
        if (is.startWith(transition.to.path, "/users")) {
            setTimeout(function () {
                $('.icon').initial({ charCount: 1, width: 30, height: 30, fontSize: 18 });
            }, 100);
        } else if (is.startWith(transition.to.path, "/apps")) {
            setTimeout(function () {
                $('.icon').initial({ charCount: 1, width: 40, height: 40, fontSize: 24 });
            }, 100);
        }
    });

    Vue.filter('datetime', function (value) {
        if (!value || value == "") {
            return "";
        }
        var d = moment(value, "X")
        return d.format("YYYY-MM-DD HH:mm:ss");
        //return d.format("YYYY-MM-DD HH:mm:ss") + "<br/>(" + d.fromNow() + ")";
    });
};

var start_vue = function () {
    var ShareTrace= Vue.extend({
        data: function () {
            return {
                user_info: user_info
            };
        },
        methods: {
            is_apps_active: function () {
                return is.startWith(this.$route.path, "/app");
            },
            is_users_active: function () {
                return is.startWith(this.$route.path, "/users");
            },
            register: function () {
                router.go("/register");
            },
            login: function () {
                router.go("/login");
            },
            logout: function () {
                var vm = this;
                fetch('/op/logout', {
                    method: 'post',
                    credentials: 'same-origin'
                }).then(function (response) {
                    return response.json();
                }).then(function (data) {
                    user_auth_ok = false;
                    user_info = null;
                    vm.user_info = null;
                    router.go("/login");
                    $.removeCookie("userinfo");
                    //console.log($.cookie("email"));
                });
            },
            apps:function(){
            	router.go("/apps");
            },
        }
    });
    router.start(ShareTrace, '#sharetrace');
};

//初始化vue
init_vue();

//确认是否有登录
var userinfo= $.cookie("userinfo");
if (userinfo== undefined){
    start_vue();
    router.go("/login");
}else{
    user_auth_ok = true;
    start_vue();
    router.app.user_info = userinfo;
}

//fetch('/op/users/1/1', { credentials: 'same-origin' })
//    .then(function (response) {
//        return response.json();
//    }).then(function (data) {
//        if (data.status == false && data.code) {
//            start_vue();
//            if (data.code == 'user_not_init') {
//                router.go("/init_user");
//            } else if (data.code == 'not_login') {
//                router.go("/login");
//            }
//        } else {
//            user_auth_ok = true;
//            fetch('/op/user/info', { credentials: 'same-origin' })
//                .then(function (response) {
//                    return response.json();
//                }).then(function (data) {
//                    user_info = data;
//                    start_vue();
//                });
//        }
//    });
