//全局变量
var router = null;
var user_auth_ok = false;
var user_auth_not_ok = true;
var user_info = null;
var COUNT_PER_PAGE = 10;

var init_vue = function () {
    router = new VueRouter({
        //'linkActiveClass': 'active'
    });

    router.map({
        '/init_user': {
            component: Register
        },
        '/register': {
            component: Register
        },
        '/login': {
            component: Login
        },
        '/apps/all/:page': {
            component: Apps,
            //auth: true
        },
        '/users/all/:page': {
            component: Users,
            //auth: true
        },
    });

    router.alias({
        '/': '/apps/all/1',
        '/apps': '/apps/all/1',
        '/users': '/users/all/1'
    });

    router.beforeEach(function (transition) {
    	console.log("beforeEach");
        if (transition.to.auth && !user_auth_ok) {
            transition.abort();
        } else {
            window.scrollTo(0, 0);
            transition.next();
        }
    });

    router.afterEach(function (transition) {
    	console.log("afterEach");
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
            	console.log("---is_apps_active");
                return is.startWith(this.$route.path, "/app");
            },
            is_users_active: function () {
                return is.startWith(this.$route.path, "/users");
            },
            is_user_info_active: function () {
                return is.startWith(this.$route.path, "/user_info");
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
                });
            },
            apps:function(){
            	router.go("/apps");
            }
        }
    });
    router.start(ShareTrace, '#sharetrace');
};

//初始化vue
init_vue();
start_vue();
