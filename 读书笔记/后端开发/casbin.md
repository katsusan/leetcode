参考:   https://www.cnblogs.com/wang_yb/archive/2018/11/20/9987397.html
        https://studygolang.com/topics/6999


basic:
    一个模型conf至少有四个部分:
        [request_definition], 
        [policy_definition], 
        [policy_effect], 
        [matchers],
        [role_definition], //如果模型使用RBAC，应该添加这个

[request_definition]:
    r = sub, obj, act
        - accessing entity (Subject)
        - accessed resource (Object)
        - the access method (Action)

[policy_definition]
    p = sub, obj, act   =>  sub对指定的obj都可以执行act
    p2 = sub, act   =>  sub对所有obj都可以执行act

[policy_effect]
    e = some(where (p.eft == allow))    =>  表示p有任意一条结果满足则为allow ??

[matchers]
    m = r.sub == p.sub && r.obj == p.obj && r.act == p.act  => 定义request和policy的匹配规则，决定p.eft为allow还是deny

[role_definition]
    g = _, _   
    g2 = _, _
    g3 = _, _, _
    => _, _ 表示用户，角色
    => _, _, _ 表示用户，角色，域


step1:
    //从MySQL中创建enforcer
    import (
        "github.com/casbin/gorm-adapter/v2"
        _ "github.com/go-sql-driver/mysql"
    )
    mysqladpt := mysqladapter.NewDBAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
    e := casbin.NewEnforcer("path/to/basic_model.conf", mysqladpt)

    //从文件中创建enforcer
    e := casbin.NewEnforcer("path/to/basic_model.conf", "path/to/basic_policy.csv")


step2:
    //验证请求是否符合规则
    //{role, req.URL.Path, req.Method}的参数格式由conf中的[request_definition]决定
    ok, err := e.Enforce(role, req.URL.Path, req.Method)

