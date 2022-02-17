doc：   https://git-scm.com/docs
        https://gitee.com/progit/

## 1. 常用命令
    git fetch <远程主机名> <分支名>     //取回分支，本地用主机名/分支来表示该分支
    git checkout -b newBranch origin/master //基于分支origin/master创建新分支newBranch
    git merge origin/master     //在当前分支上合并origin/master
    git rebase orgin/master    //以分支origin/master为基础把分离以后的commit并入该分支生成新的当前分支
    git pull <远程主机名> <远程分支>:<本地分支> //取回远程分支并与本地分支合并 (本地分支省略的话则与当前分支合并)
        git pull origin master:next = git fetch origin master && git merge origin/master
    git branch --set-upstream master origin/next    //指定master分支追踪origin/next分支
    git push <远程主机名> <本地分支>:<远程分支> //推送本地分支到远程主机上的远程分支
        git push origin :master //本地分支不指定的话相当于推送了空分支，会删除远程master分支
        git push origin --delete master    //删除master分支     

    
    git log --oneline $file     //查看file的commit历史
    git show $commitid:$file    //查看commitid提交时file的提交内容

    //修改本地和远程分支名
    git branch -m new-name      //把当前分支改为new-name
    git branch -m old-name new-name //把old-name分支改为new-name
    git push origin :old-name new-name //删掉remote的old-name分支并push本地的new-name分支，此时本地分支只有gitignore文件
    ~此处要先切换到要操作的分支，git checkout new-name
    git push origin -u new-name //把当前分支的upstream分支该问new-name

## 2. pull request
	- fork the project
	- git clone git@github.com:Katsusan/xxxx.git	// the forked project in your repo
	- git checkout -b featureA		// create new branch
	- git commit -a -m "add new featureA"
	- git push origin featureA
	- issue a pull request in you repo	// src is your featureA branch, while dst is the branch of project you want to contribute to
refer: https://www.jianshu.com/p/bd73bf2f90d2


++++++++++++++++++++++++++++++++++++++++++++
//github publickey
- ssh-keygen -t rsa
- 将生成的公钥id_rsa.pub拷贝到https://github.com/settings/keys下
- ssh git@github.com ，出现 You've successfully authenticated则表明验证ok
++++++++++++++++++++++++++++++++++++++++++++


## 3. git workflow
refer: 
    https://github.com/findxc/blog/issues/22
    http://www.ruanyifeng.com/blog/2015/12/git-workflow.html

某版本开发中：   
- 从 develop 上切出 feature-xxx 进行开发
- 开发完后提交合并请求合并回 develop 分支。在提合并请求的时候可以勾选合并后删除该分支，就不用手动去删除了。

```
git checkout -b featureA
git commit -a -m "add new featureA"
git push origin featureA


git branch -d featureA      // delete local branch
git push origin --delete featureA   // delete remote branch
```


开发完了提测：   
- 测试过程中的 bug 修复从 release-xxx 分支切 fix-xxx 出去进行修复
- 修复好了提交合并请求合并回 release-xxx 分支
- 其它本次不上线的功能开发继续在 develop 上开发
- 提测之前从 develop 分支切 release-xxx 分支作为测试用的分支


测试通过准备上线：  
- release-xxx 分支合并回 develop 分支
- release-xxx 分支合并回 master 分支
- master 分支打 tag
- 上线 master 代码


线上环境 bug 修复：
- 从 master 分支切 hotfix-xxx 出来进行修复
- 修复完后用 hotfix-xxx 作为测试分支进行测试
- 测试通过后合并回 develop 分支和 master 分支，如果存在 release-xxx 分支也要合并回 release-xxx 分支
- master 分支打 tag
- 上线 master 代码

