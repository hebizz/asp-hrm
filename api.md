[toc]

## asp-hrm

### 新建部门

**url**: `/api/v1/department/create`

**method**: `post`

**url params**: None

**request body**: 
```
{
     "title": xx           # 部门名字      [必填]
     "pid":  xx           # 上级部门id    [必填]
}
```

**success response**:

```
{
    "c": "200",
    "msg": "success",
    "data": null
}
```



---

#### 编辑部门

**url**: `/api/v1/department/update`

**method**: `POST`

**url params**: None

**request body**:
```
{
    "id": xxxx              # 部门id    [必填]
    "title": xxx             # 部门名称   [必填]
    "pid": xxx              # 父部门id   [必填]
}
```

**success response**
```
{
    "c": "200",
    "msg": "success",
    "data": null
}   
```



---
#### 删除部门

**url**: `/api/v1/department/remove`

**method**: `POST`

**url params**: None

**request body**:
```
{
    "id": xxxx              # 部门id    [必填]
}
```

**success response**
```
{
    "c": "200",
    "msg": "success",
    "data": null
}   
```



---

#### 人员注册

**url**: `/api/v1/human/register`

**method**: `POST`

**url params**: None

**request body**:
form-data
```
     "title" =  xx            # 人员姓名 [必填]
     "did"  =   xx            # 部门id  [必填]
     "dname" =  xx           # 部门名称 [必填]
     "image" = xx            # 图片     [必填]
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```



---

#### 删除人员

**url**: `/api/v1/human/remove`

**method**: `POST`

**url params**: 

```
{
     "ids": [xx, xx]                # 人员id列表     [必填]
}
```

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```



---

#### 编辑人员

**url**: `/api/v1/human/update`

**method**: `POST`

**url params**: None

**request body**:
form-data
```
     "id" =  xx                # 人员id     [必填]
     "title" =  xx                # 人员姓名  
     "image" =  xx               # 图片      
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```


---

#### 获取人员信息

**url**: `/api/v1/human/query`

**method**: `GET`

**url params**: 
`   id = xx    用户id   [可选],
    did = xx   部门id    [可选],
    offset = 0 偏移量    [必选],
    limit = xx 限制量    [必选]
`

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": {
        "data":[{
        "title": xx,  用户姓名
        "id": xx,    用户id
        "did": xx,   部门id
        "dname": xx, 部门姓名
        "path": xx, 图片路径
        "timestamp": xx 录入时间}]
    }
}
```
---

---

#### 获取部门组

**url**: `/api/v1/department/query`

**method**: `GET`

**url params**: None

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": [{"title": 产品部,  部门名称
              "id": xx,       部门id
              "type": "department", 部门 
              "hList": [{"title": xx, "id": xx, "type":"human"}, {"title": xx, "id": xx, "type":"human"}, 人员列表
              "children": [{"title":xx, 
                         "id": xx, 
                         hList:[{"title": xx, "id": xx, "type":"human"}, {"title": xx, "id": xx, "type":"human"}], 部门列表, 如果为空列表,则无数据
              }]
              }]
}
```
---

#### 获取所有部门信息

**url**: `/api/v1/departmentList/query`

**method**: `GET`

**url params**: None

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": [
        {
            "id": "1",
            "pid": "",
            "title": "未分组"
        },
        {
            "id": "60617ebea5070846e2fc02f8",
            "pid": "",
            "title": "体育部"
        },
        {
            "id": "60617ec5a5070846e2fc02f9",
            "pid": "",
            "title": "文艺部"
        },
        {
            "id": "60617ee9a5070846e2fc02fb",
            "pid": "60617ebea5070846e2fc02f8",
            "title": "商业部"
        }
    ]
}
```
---

#### 识别日志

**url**: `/api/v1/identify/log`

**method**: `GET`

**url params**:
`  
humanType = xx    人员类型    [可选]  0: 系统人员, 1:陌生人, 不传为所有人员
title = xx         人员姓名    [可选]
starttime = xx    开始时间    [必选]
endtime = xx      结束时间    [必选]
offset = 0        偏移量      [必选]
limit = xx        限制量      [必选]
`

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": [
        {
            "id": "",
            "humanType":   ,  // 人员类型,  0: 系统人员, 1:陌生人
            "title": "",       // 人员姓名
            "dname": "",      // 所属部门
            "timestamp": ,    // 识别时间
            "device": "",     // 设备名称
            "position": "",   //安装位置
            "rawPath": "",    //原始照片
            "aiPath": "",     //ai识别照片
        },
    ]
}
```
---