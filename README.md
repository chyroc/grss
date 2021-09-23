# grss
Rss hub tool combining git and action

## 数据存储结构

数据存储需要考虑以下问题：

- 区分老数据和新数据，即让 feed api 始终返回新数据
  - json 数据是全量数据，按 utc date 切割
  - xml 只保存最新的 feed
- 能够让静态的页面支持参数，如 `/<user-id>/feed`
  - user_id 需要提前写入到 git 中，以便让 git 知道需要抓取这个页面

```shell
- json
  - 2021-01-02
      - zhihu
        - bookstore
          - newst.json
        - user
          - <user_id>
            - feed.json
- feed
  - zhihu
    - bookstore
      - newst.xml
    - user
      - <user_id>
        - feed.xml
```