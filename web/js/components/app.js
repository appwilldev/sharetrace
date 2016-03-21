//应用列表
var Apps = function (resolve, reject) {
    var template_url = 'components/apps.html';
    return fetch(template_url)
        .then(function (template_resp) {
            return template_resp.text();
        }).then(function (template) {
            resolve({
                template: template,
                route: {
                    data: function (transition) {
                        var page = transition.to.params.page;
                        var path;
                        if (is.startWith(transition.to.path, '/apps/all')) {
                            path = '/op/apps/all/' + transition.to.params.page + '/' + COUNT_PER_PAGE;
                        }
                        return fetch(path, { method: 'GET', credentials: 'same-origin' })
                            .then(function (data_resp) {
                                return data_resp.json();
                            }).then(function (data) {
                                if (data && data.data != null) {
                                    var app;
                                    var page_count = null;
                                    var app_data = data.data;
                                    if (path == '/op/apps/all/' + page + '/' + COUNT_PER_PAGE) {
                                        for (i in data.data) {
                                            page_count = 0;
                                            app = data.data[i];
                                            if (data.total % COUNT_PER_PAGE == 0) {
                                                page_count = data.total/ COUNT_PER_PAGE;
                                            } else {
                                                page_count = Math.floor(data.total/ COUNT_PER_PAGE) + 1;
                                            }
                                        }
                                    } else {
                                    	alert("path error");
                                    }
                                }
                                return {
                                    data: app_data,
                                    app: null,
                                    operater: '',
                                    modal_title: '',
                                    judge_value: 1,
                                    page_count: page_count,
                                    err_msg: null,
                                    keyword: null,
                                }
                            })
                    }
                },
                methods: {
                    //保存App数据
                    save_data: function () {
                        var vm = this;
                        if (vm.judge_value == 1) {
                            //增加APP
                            fetch('/op/app', {
                                method: 'POST',
                                body: JSON.stringify({
                                    "appid": vm.app.appid,
                                    "appname": vm.app.appname,
                                    "appschema": vm.app.appschema,
                                    "apphost": vm.app.apphost,
                                    "status": parseInt(vm.app.status),
                                    //"share_click_money": parseInt(vm.app.share_click_money),
                                    "share_install_money": parseFloat(vm.app.share_install_money),
                                    "install_money": parseFloat(vm.app.install_money),
                                    //"appicon": vm.app.icon
                                }),
                                credentials: 'same-origin'
                            }).then(function (response) {
                                return response.json();
                            }).then(function (data) {
                                if (data.status == true) {
                                    $('#myModal').modal('hide');
                                } else {
                                	alert(data.msg);
                                }
                            });
                        } else if (vm.judge_value == 2) {
                            //修改APP
                            fetch('/op/app', {
                                method: 'PUT',
                                body: JSON.stringify({
                                    "id": vm.app.id,
                                    "appid": vm.app.appid,
                                    "appname": vm.app.appname,
                                    "appschema": vm.app.appschema,
                                    "apphost": vm.app.apphost,
                                    "status": parseInt(vm.app.status),
                                    //"share_click_money": parseInt(vm.app.share_click_money),
                                    "share_install_money": parseFloat(vm.app.share_install_money),
                                    "install_money": parseFloat(vm.app.install_money),
                                    //"icon": vm.app.icon
                                }),
                                credentials: 'same-origin'
                            }).then(function (response) {
                                return response.json();
                            }).then(function (data) {
                                if (data.status == true) {
                                    $('#myModal').modal('hide');
                                } else {
                                	alert(data.msg);
                                }
                            });
                        }
                    },
                    //新增App
                    show_add_app: function () {
                        this.err_msg = null;
                        this.judge_value = 1;
                        this.app = {
                            "appid": "",
                            "appname": "",
                            "appschema": "",
                            "apphost": "",
                            "status": "",
                            //"share_click_money": "",
                            "share_install_money": "",
                            "install_money": "",
                            "icon": "",
                        };
                        this.modal_title = '增加应用';
                        this.operater = '增加';
                    },
                    //修改App
                    show_edit_app: function (app) {
                        this.err_msg = null;
                        this.judge_value = 2;
                        this.app = clone(app);
                        this.modal_title = '修改应用';
                        this.operater = '保存';
                    },
                    all_apps: function () {
                        router.go('/apps')
                    },
                    show_app_data: function (app) {
                    	var url = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/web/components/stats.html?appid=" + app.appid + "&appname=" + app.appname;
                    	console.log(url);
                    	window.open(url);

                    },
                    show_host_data: function (app) {
                    	var url = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/web/components/hoststats.html?host=" + app.apphost;
                    	console.log(url);
                    	window.open(url);

                    },
                    show_app_money: function (app) {
                    	var url = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/web/components/moneystats.html?appid=" + app.appid + "&appname=" + app.appname;
                    	console.log(url);
                    	window.open(url);

                    },
                }
            });
        });
};
