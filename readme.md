# file
## 说明
文件操作框架 用户基本文件io或者状态监控

## file.open
- 打开写的文件句柄 不存在会默认创建
- ud = file.open{name , path , delim}
- 满足lua.writer接口
```lua
    local ud = file.open{
        name = "demo",
        path = "/var/logs/a.log",
        delim = "\n",
    }
    
    start(ud)
```

### 接口函数
- [ud.backup()]()
- [ud.push()]()
```lua
    local ud = file.open("x-yyyy-MM-dd.hh-mm")
    ud.backup() --执行当前目录的事件 根据当前时间
    ud.push("xxx %s" , "helo") --写入
```

## file.stat
- 获取文件状态
- stat = file.stat(path)
### 基础接口
- [stat.ok]()
- [stat.name]()
- [stat.ext]()
- [stat.mtime]()
- [stat.ctime]()
- [stat.atime]()
- [stat.path]()
- [stat.dir]()
```lua
    local st = file.stat("a.txt")
    print(st.ok)
    print(st.name)
    print(st.ext)
    -- etc
```

## file.dir
- 打开文件目录
- d = file.dir(path)
#### 基础接口
- [d.ok]()
- [d.err]()
- [d.count]()
- [d.grep()]()
- [d.ipairs()]()

```lua
    local d = file.dir("/var/logs")
    print(d.err)

    d.grep("*" , function(stat)
        print(stat.ok)        
        --todo
    end)

    d.ipairs(function(stat)
        print(stat.ok)
        --todo
    end)
```

## file.walk
- 新建一个文件遍历器
- walk = file.walk(name)
### 基础接口
- [walk.deep(int)](walk) 
- [walk.dir(bool)](walk)
- [walk.pipe(object)](walk)
- [walk.filter(string...)](walk)
- [walk.ignore(string...)](walk)
- [walk.run()](walk)

### 内部变量 info
- info.path
- info.ok
- info.name
- info.ext
- info.mtime
- info.ctime
- info.atime
- info.dir

### walk

```lua
    local wk = file.walk("/var/log" , "/usr/local")
    wk.deep(10)                  --深度设置
    wk.dir(true)                 --目录处理
    wk.filter("*.log" , "*.jar") --匹配字符串
    wk.ignore("*.zip")           --忽略
    wk.pipe(function(info)
        print(info.path)
    end)
    wk.run()
```

## file.glob
- 文件模糊模式搜索 
- glob= file.glob(string)

### 基础接口
- [glob.pipe(object)](glob)
- [glob.run()](glob)
- [glob.wrap()](glob)
- [glob.result](glob)

### glob
```lua
    local glob = file.glob("/var/log/*/*.log")
    glob.pipe(function(info)
        print(info)
    end)

    print(glob.result) --slice数据返回
```

## file.dir
- 打开文件目录
- dir = file.dir(string)

### 基础接口
- [dir.pipe(object)](dir)
- [dir.filter(string...)](dir)
- [dir.err](dir)
- [dir.count](dir)
- [dir.ok](dir)
- [dir.result](dir)

### dir
```lua
    local dir = file.dir("/var/log")
    dir.filter("*.log")
    dir.pipe(function(info)
        print(info)
    end)

    print(dir.result) --slice数据返回
    print(dir.count)  --数量
    print(dir.ok)     --是否正常
    print(dir.err)    --错误信息
```

## file.stat
- 查看文件状态 
- info = file.stat(string) [info](内部变量 info)
- 
```lua
    local stat = file.stat("/var/log")
    print(stat.path)
    print(stat.name)
    print(stat.ok)
    print(stat.err)
```
## file.cat , file.read_all
```lua
    local v = file.cat("/var/log/*")
```