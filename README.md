# tmp-rank
## 问题分析
- 日活10W+，最高同时在线按5W算，月积累用户按100W算。
- 活动结束之后活动建立排行榜，不需要实时。
- 多维度排序。
- 需要查询自己名次前后的十名玩家的分数和名次。

## 实现方式
### 1.RDB存储加排序
- 如果是活动结束后的排行榜建立，可以使用rdb进行数据的存储。
- 存储涉及到分库分表。
- 活动结束后，涉及到分库分表的数据汇总和排序和排序数据再存储。
### 2.Redis-zset
- 依赖redis的zset进行数据的实时存储和排序。
- 日活10W加的单key压力还好。
- 如果依赖插入时间，可以最近未来时间-插入时间和score,进行整数编码或者拼接小数，可业务或lUA脚本实现。

### 3.自定义存储结构
#### 1.store_array
- 代码：```core/store/store_array.go```
- map存储玩家数据，活动结束后将玩家的数据搬到数组中进行排序
- 优点：结构简单，存储和获取排行效率更高 O(1)
- 缺点：排行榜不是实时更新，需要活动结束后调用才建立排序结构。
- UML类图
```UML
+-------------------------------------+
| <<interface>>                       |
| Element                             |
+-------------------------------------+
| + GetKey(): int64                   |
| + Compare(other: Element): int      |
+-------------------------------------+

+-------------------------------------+
| <<struct>>                          |
| ArraySortWrapper[T]                 |
+-------------------------------------+
| - value: T                          |
| - idx: int                          |
+-------------------------------------+
| No specific methods                 |
+-------------------------------------+

+-------------------------------------+
| <<struct>>                          |
| ArraySortStore[T]                   |
+-------------------------------------+
| - sorted: bool                      |
| - dict: map[int64]*ArraySortWrapper[T] |
| - sort: []*ArraySortWrapper[T]      |
+-------------------------------------+
| + NewArraySortStore(): *ArraySortStore[T] |
| + Add(psc: T): void                 |
| + Sort(): void                      |
| + GetSort(key: int64, rangeNum: int): ([]T, int64) |
+-------------------------------------+

```

#### 2.store_skip_list
- 代码：```core/store/store_skip_list.go```
- 跳表实现，实时维护和更新排行榜数据
- 效率：添加和获取O(logN)
- UML类图
```UML

+-------------------------------------+
| <<interface>>                       |
| Element                             |
+-------------------------------------+
| + GetKey(): int64                   |
| + Compare(other: Element): int      |
+-------------------------------------+

+-------------------------------------+
| <<struct>>                          |
| SkipListStore[T]                    |
+-------------------------------------+
| - skipList: *SkipList[T]            |
+-------------------------------------+
| + NewSkipListStore(): *SkipListStore[T] |
| + Add(t: T): void                   |
| + GetSort(key: int64, rangeNum: int): ([]T, int64) |
+-------------------------------------+

+-------------------------------------+
| <<struct>>                          |
| SkipList[T]                         |
+-------------------------------------+
| - header: *Node[T]                  |
| - tail: *Node[T]                    |
| - level: int                        |
| - dict: map[int64]*Node[T]          |
+-------------------------------------+
| + New(): *SkipList[T]               |
| + createNode(level: int, psc: T): *Node[T] |
| + insert(value: T): *Node[T]        |
| + deleteNode(x: *Node[T]): void     |
| + delete(value: T): bool            |
| + AddOrUpdate(value: T): bool       |
| + FindRangeNum(key: int64, distance: int): ([]T, int64) |
| + Look(): void                      |
| + LookTail(): void                  |
+-------------------------------------+

+-------------------------------------+
| <<struct>>                          |
| Node[T]                             |
+-------------------------------------+
| - Value: T                          |
| - backward: *Node[T]                |
| - level: []Level[T]                 |
+-------------------------------------+
| No specific methods                 |
+-------------------------------------+

+-------------------------------------+
| <<struct>>                          |
| Level[T]                            |
+-------------------------------------+
| - forward: *Node[T]                 |
| - span: int64                       |
+-------------------------------------+
| No specific methods                 |
+-------------------------------------+


```