# 蜜罐列表
## 新建蜜罐

| 用例编号  | 测试标题  | 重要级别   | 预置条件  |  测试步骤 | 预期结果  | 实际结果  |
| ------------ | ------------ | ------------ | ------------ | ------------ | ------------ | ------------ |
| AddHoneyPod_1  | 正常新建蜜罐  | 高级别  | 无 | a.新建蜜罐名称为test<br>b. 选择镜像<br>c. 点击新建   | 新建成功  |   |
| AddHoneyPod_2 |  新建重名蜜罐   | 高级别  | 已存在名称为test的蜜罐  | a. 新建蜜罐名称为test<br>b. 选择镜像<br>c. 点击新建   | 提示同名蜜罐已存在  |   |
| AddHoneyPod_3  | 新建蜜罐无法访问Harbor  | 高级别  | Harbor服务器无法访问  | a. 新建蜜罐名称为test<br>b. 选择镜像<br>c. 点击新建 | 新建失败  |   |
| AddHoneyPod_4  | 新建蜜罐镜像不正确  | 高级别  | 镜像无常驻进程启动失败  | a. 新建蜜罐名称为test<br>b. 选择镜像<br>c. 点击新建   | 新建失败  |   |