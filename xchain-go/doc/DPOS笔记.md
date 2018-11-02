# DPOS 笔记

## epoch列表存储
```
key: "epoch" + epochNum
value:  RLP(map[Address]Validator)

```

## Candidates列表存储

```
key: "candidates" 
value:  RLP(map[Address]Candidate)

```



## TODO:
1. epoch列表，candidates列表存入DB
2. 新增DB后各流程修改：投票，取消投票，增加候选人，踢出候选人，选举  及其测试用例更新
2. 修改log
3. 新增出块流程
