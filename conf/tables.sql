-- share_url: http://hostname?appid=abc&itemid=123&fromid=456&channel=weixin&v=1.1, 域名不一样，item不一样，fromid不一样
-- fromid : 最好是ida
-- app server post appid, fromid, share_url to sharetrace saas
-- cache share_url(share_url) share_url(id)

CREATE SEQUENCE share_url_id START 2016 NO CYCLE;                                        
CREATE TABLE share_url(                                                                  
    id BIGINT NOT NULL PRIMARY KEY,                                                 
    share_url VARCHAR(2048) NOT NULL,                                                
    fromid  VARCHAR(256) NOT NULL,                                             
    appid   VARCHAR(256) NOT NULL,                                             
    itemid  VARCHAR(256) DEFAULT NULL,                                             
    channel VARCHAR(256) DEFAULT NULL,                                             
    ver VARCHAR(256) DEFAULT NULL,                                             
    des TEXT DEFAULT NULL,                                                          
    status INT DEFAULT 0,                                                           
    created_utc INT                                                                 
);                                                                                  
CREATE INDEX uidx_su_share_url ON share_url(share_url);                                 
CREATE INDEX idx_su_fromid ON share_url(fromid);                                    
CREATE INDEX idx_su_itemid ON share_url(itemid);                                    
CREATE INDEX idx_su_appid ON share_url(appid);                                    
CREATE INDEX idx_su_channel ON share_url(channel);                                    
CREATE INDEX idx_su_ver ON share_url(ver);                                    
CREATE UNIQUE INDEX uidx_su_urlfrom ON share_url(share_url, fromid, itemid, appid);       

-- cookieid: 唯一cookieid，SAAS生成, id_shareid
-- shareid : share_url(id)
-- installid: 安装者的id，最好是ida
-- click_type: 0: st_cookie, 1: web agent type 
-- agentid: ex: md5(agentip, agent)
-- des: click user info:ip, browser, version
-- app server post appid, fromid, share_url to sharetrace saas
-- cache click_session(cookieid), click_session(id)

CREATE SEQUENCE click_session_id START 2016 NO CYCLE;                                        
CREATE TABLE click_session(
    id BIGINT NOT NULL PRIMARY KEY,                                                 
    shareid BIGINT  NOT NULL,                                                
    cookieid VARCHAR(2048) NOT NULL,                                                
    installid  VARCHAR(256) DEFAULT NULL,                                             
    click_type INT DEFAULT 0,                                                           
    agent VARCHAR(1024) DEFAULT NULL,                                             
    agentip VARCHAR(256) DEFAULT NULL,                                             
    agentid VARCHAR(1024) DEFAULT NULL,                                             
    des TEXT DEFAULT NULL,                                                          
    status INT DEFAULT 0,                                                           
    created_utc INT                                                                 
);                                                                                  
CREATE INDEX uidx_cs_cookieid ON click_session(cookieid);                                 
CREATE INDEX idx_cs_shareid ON click_session(shareid);                                 
--ALTER TABLE click_session ADD COLUMN click_type int default 0;
--ALTER TABLE click_session ADD COLUMN agent varchar(1024) default null;
--ALTER TABLE click_session ADD COLUMN agentip varchar(256) default null;
--ALTER TABLE click_session ADD COLUMN agentid varchar(1024) default null;
-- CREATE INDEX idx_cs_date ON click_session(date(to_timestamp(created_utc));
-- select count(*), date(to_timestamp(created_utc)) from click_session group by date(to_timestamp(created_utc)) order by date(to_timestamp(created_utc));

-- 账号管理
CREATE SEQUENCE user_id START 2016 NO CYCLE;
CREATE TABLE user_info (
    id BIGINT NOT NULL PRIMARY KEY,
    email VARCHAR(256) NOT NULL,
    passwd  VARCHAR(256) NOT NULL,
    name VARCHAR(256) DEFAULT NULL,
    des TEXT DEFAULT NULL,                                                          
    status INT DEFAULT 0,                                                           
    created_utc INT NOT NULL
);
CREATE INDEX uidx_ui_email ON user_info(email);                                 

-- App管理
-- 一个账号可以有多个App， 一个App只属于一个账号
CREATE SEQUENCE app_id START 2016 NO CYCLE;
CREATE TABLE app_info (
    id BIGINT NOT NULL PRIMARY KEY,
    appid VARCHAR(256) NOT NULL,
    appname  VARCHAR(256) NOT NULL,
    appicon VARCHAR(2048) DEFAULT NULL,
    userid BIGINT NOT NULL, 
    des TEXT DEFAULT NULL,                                                          
    status INT DEFAULT 0,                                                           
    created_utc INT NOT NULL
);
CREATE INDEX uidx_ai_appid ON app_info(appid);                                 


-- 同appid， 同ida，只算一个? 客户端要判断，用户是否已经安装过了
-- 用户安装了，删除了，又安装了, 怎么算
