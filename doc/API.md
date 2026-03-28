## 规范

### 参数和返回值

* 端点：`http://127.0.0.1:31943`
* 均是 POST 方法
* 需要带参的接口，参数为 JSON 字符串，放置到 body 里，标头 Content-Type 为 `application/json`
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": {}
  }
  ```
    * `code`：非 0 为异常情况
    * `msg`：正常情况下是空字符串，异常情况下会返回错误信息
    * `data`：可能为 `{}`、`[]` 或者 `NULL`，根据不同接口而不同

## 1 应用

#### 1.1 退出后端服务

* > /api/v1/app/exit
* 参数
    * 无请求参数
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

## 2 账本

#### 2.1 创建指定名称的账本

* > /api/v1/ledger/create-one
* 参数
  ```json
  {
    "name": "test-name"
  }
  ```
    * `name`：账本名称
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": "0197f526-8160-70e6-84f7-16f6f5c8417e"
  }
  ```
    * `data`：账本ID

#### 2.2 删除指定ID的账本

* > /api/v1/ledger/delete-one
* 参数
  ```json
  {
    "id": "账本id"
  }
  ```
    * `id`：账本ID
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

#### 2.3 修改指定ID的账本的名称

* > /api/v1/ledger/modify-name
* 参数
  ```json
  {
    "id": "0197fa7b-9b5b-73a4-b131-f022cabd1cf5",
    "name": "test"
  }
  ```
    * `id`：账本ID
    * `name`：账本新名称
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

#### 2.4 查询所有账本

* > /api/v1/ledger/query-all
* 参数
  ```json
  {
    "id": "all"
  }
  ```
    * `id`：账本ID
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
        {
            "id": "0197fa7b-9b5b-73a4-b131-f022cabd1cf5",
            "name": "test-name",
            "created_at": 1760251611,
            "updated_at": 1760251611
        }
    ]
  }
  ```
    * `data`：账本列表

## 3 消费记录

#### 3.1 创建消费记录

* > /api/v1/tr/create-one
* 参数
  ```json
  {
    "ledger_id": "0197fa7b-9b5b-73a4-b131-f022cabd1cf5",
    "price": 3433,
    "transaction_type": "expense",
    "category": "购物消费",
    "description": "test des",
    "tags": ["tags1","tags2"],
    "transaction_at": 1760251611
  }
  ```
    * `ledger_id`：账本ID
    * `price`：消费价格
    * `transaction_type`：交易类型，取值范围['expense','income','transfer']，分别表示支出，收入，转账
    * `category`：消费分类
    * `description`：消费描述
    * `tags`：消费标签列表
    * `transaction_at`：消费时间，秒级时间戳
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": "0199d72c-0868-7d34-90ad-31ecb4f59fc9"
  }
  ```
    * `data`：消费记录id

#### 3.2 删除指定ID的消费记录

* > /api/v1/tr/delete-by-id
* 参数
  ```json
  {
    "trId": "0199d72c-0868-7d34-90ad-31ecb4f59fc9"
  }
  ```
    * `trId`：消费记录ID
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

#### 3.3 条件查询消费记录

* > /api/v1/tr/query
* 参数
  ```json
  {
    "ledgerId": "0197fa7b-9b5b-73a4-b131-f022cabd1cf5",
    "offset": 10,
    "limit": 10,
    "tsRange": [1760251611,1760251611],
    "transactionTypes": ["expense","income"],
    "categoryTags":{
      "购物消费": ["衣物","家具"]
    }
  }
  ```
    * `ledgerId`：账本ID
    * `offset`：偏移量
    * `limit`：查询条数
    * `tsRange`：查询时间戳范围内的记录，秒级时间戳
    * `transactionTypes`：交易类型列表
    * `categoryTags`：需要匹配的分类标签条件
* 返回值
  ```json
  {
    "code": 0,
    "data": {
        "items": [
            {
                "category": "餐饮美食",
                "created_at": 1760251611,
                "description": "早餐",
                "ledger_id": "0199d72c-0868-7d34-90ad-31ecb4f59fc9",
                "price": 4.5,
                "tags": [
                    "三餐"
                ],
                "transaction_at": 1760251611,
                "transaction_id": "01981e51-c6df-716f-953f-dc8771370af8",
                "transaction_type": "expense",
                "updated_at": 1760251611
            }
        ],
        "total": 10,
        "trStatistics": {
            "expense": 10,
            "income": 0,
            "transfer": 0
        }
    },
    "msg": ""
  }
  ```
    * `data`：查询结果
    * `items`：消费记录列表
    * `total`：所有符合条件的记录数量
    * `trStatistics`：所有符合条件的记录的统计数据，包括总支出，总收入，总转账

## 4 消费分类

#### 4.1 查询指定交易类型下的分类

* > /api/v1/category/query/:type
* 路径参数
    * `type`：交易类型
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": ["购物消费","生活缴费"]
  }
  ```
    * `data`：分类列表

## 5 消费标签

#### 5.1 查询指定分类下的标签

* > /api/v1/tag/query/:category
* 路径参数
    * `category`：分类名称
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": ["衣物","家居"]
  }
  ```
    * `data`：标签列表

## 6 工作空间

#### 6.1 指定目录打开工具空间

* > /api/v1/workspace/open
* 参数
  ```json
  {
    "workspaceDir": "path/to/workspace"
  }
  ```
    * `workspaceDir`：工作空间目录路径，目录为空或目录下存在合法数据
* 返回值
  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```
  