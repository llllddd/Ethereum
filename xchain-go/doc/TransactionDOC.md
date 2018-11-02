### 合约交易相关的方法：

#### EstimateGas
计算一笔交易执行所需要花费的gaslimit
`accounts/abi/bind/backends/simulated.go`
1. 确定二分法查询时的最高和最低的gas量
2. 创建一个帮助程序来检查gas的容量是否能够执行一笔可执行交易，满足则返回true，不满足则返回false。 还涉及到了快照的读取与回滚
3. 执行二分搜索，并根据最小的执行需要的gas值lo-21000-1，以及获得的所提供的最大的gas值-hi，来确定执行交易所需要的gas
4. 如果lo+1>=hi,跳出上述的循环，即hi=最初设置的cap的值，没有改变。（目前看这种情况就是hi=21000)如果交易在最大的允许内，仍失败，则交易无效


#### IntrinsicGas
返回所对应的固定的gasprice。
`core/state_transition.go`


#### SuggestGasPrice
`accounts/abi/bind/backends/simulated.go`


#### 普通提交转账交易


#### cli命令行发起交易
> 参考internal里的方法



`PrivateAccountAPI`结构体

NewPrivateAccountAPI

> func (s *PrivateAccountAPI) SendTransaction(ctx context.Context, args SendTxArgs, passwd string) (common.Hash, error) {
使用节点上的from账户的私钥，给tx签名,并将交易提交到txpool

> func (s *PrivateAccountAPI) signTransaction(ctx context.Context, args SendTxArgs, passwd string) (*types.Transaction, error) {
// signTransactions设置默认值并对给定的事务进行签名

> func submitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
// submitTransaction 是一个帮助函数，它将tx提交给txPool并记录消息。


> func (args *SendTxArgs) setDefaults(ctx context.Context, b Backend) error {
//setDefaults 设置一些健全性默认值并在失败时终止
    
