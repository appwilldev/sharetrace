//注册
var Register = function (resolve, reject) {
    var template_url = 'components/register.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                data: function () {
                    return {
                        name: null,
                        pass_code: null,
                        test_pass_code: null,
                        message: null,
                        email: null,
                        email_error: null,
                        icon: null
                    };
                },
                methods: {
                    init_user: function () {
                        if (is.not.existy(this.pass_code) || this.pass_code.length < 6) {
                            this.message = '密码长度小于6个字符，请重新输入。';
                        } else {
                            if (this.test_pass_code != this.pass_code) {
                                this.message = '两次输入密码不一致，请重新输入。';
                            } else {
                                var testEmail = /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/;
                                if (!testEmail.test(this.email)) {
                                    this.email_error = '邮箱格式有误！！！';
                                } else {
                                    fetch('/op/user/init', {
                                        method: 'post',
                                        body: JSON.stringify({
                                            "name": this.name,
                                            "passwd": this.pass_code,
                                            "email":this.email 
                                        }),
                                        credentials: 'same-origin'
                                    }).then(function (response) {
                                        return response.json();
                                    }).then(function (data) {
                                        if (data.status == true) {
                                            router.go("/login");
                                        }else{
                                        	alert(data.msg);
                                        }
                                    });
                                }
                            }
                        }
                    }
                }
            });
        });
};

//登录
var Login = function (resolve, reject) {
    var template_url = 'components/login.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                data: function () {
                    return {
                        name: null,
                        pass_code: null,
                        username_error: null,
                        user_password_error: null
                    };
                },
                methods: {
                    login: function () {
                        var vm = this;
                        fetch('/op/login', {
                            method: 'post',
                            body: JSON.stringify({
                                "email": this.email,
                                "passwd": this.pass_code
                            }),
                            credentials: 'same-origin'
                        }).then(function (response) {
                            return response.json();
                        }).then(function (data) {
                            if (data.status == true) {
                                index = 2;
                                user_auth_ok = true;
                                console.log("user_auth_ok:",  user_auth_ok);
                                user_info = data;
                                //TODO if admin go to users, else go to apps
                                router.go("/apps");
                            } else {
                                if (is.startWith(data.msg, 'user not exist')) {
                                    vm.username_error = "用户不存在";
                                    vm.user_password_error = "";
                                } else if (is.startWith(data.msg, 'user passcode wrong')) {
                                    vm.username_error = "";
                                    vm.user_password_error = "密码输入错误";
                                }
                            }
                        });
                    }
                }
            });
        });
};

//用户列表
var Users = function (resolve, reject) {
    var template_url = 'components/users.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                route: {
                    data: function (transition) {
                        var page = transition.to.params.page;
                        return fetch('/op/users/' + page + '/' + COUNT_PER_PAGE, { credentials: 'same-origin' })
                            .then(function (data_resp) {
                                return data_resp.json();
                            }).then(function (data) {
                                if (data.data!= null) {
                                    var u = null;
                                    for (i in data.data) {
                                        u = data.data[i];
                                    }
                                }
                                var page_count = 0;
                                if (data.data.total_count % COUNT_PER_PAGE == 0) {
                                    page_count = data.total/ COUNT_PER_PAGE;
                                } else {
                                    page_count = Math.floor(data.total/ COUNT_PER_PAGE) + 1;
                                }
                                console.log("--- users data:", data);
                                return {
                                    data: data.data,
                                    password_error: '',
                                    email_error: '',
                                    username_already_exists_error: '',
                                    name: '',
                                    pass_code: '',
                                    confirm_pass_code: '',
                                    email: '',
                                    icon: '',
                                    aux_info: u.aux_info,
                                    page_count: page_count
                                };
                            });
                    }
                },
                methods: {
                }
            });
        });
};
