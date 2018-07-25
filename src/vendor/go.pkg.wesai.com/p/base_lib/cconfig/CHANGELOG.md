Centralized Configuration (client) - Change log 1.x
==============================

## 1.1

### 功能：
#### 1.【变更】本库使用的etcd client
由于在etcd 2.2的release note中已不推荐使用go-etcd库，并建议用etcd/client库替代，所以本库改用etcd/client。据称go-etcd不会再随着etcd版本更新了。