package basic

import (
	"xchain-go/common"
)

type DerivableList interface {
	// Len() int
	// GetRlp(i int) []byte
}

func DeriveSha(list DerivableList) common.Hash {

	// keybuf := new(bytes.Buffer)
	// //创建空树
	// trie := new(trie.Trie)
	// //迭代列表中的每一项，用其更新该MPT
	// for i := 0; i <= list.Len(); i++ {
	// 	//重置keybuf
	// 	keybuf.Reset()
	// 	//使用列表中每一项的序号作为key，先对其rlp编码
	// 	rlp.Encode(keybuf, uint(i))
	// 	//update()会将key和value在MPT中联系起来
	// 	//实际会调用trie.Tryupdate
	// 	trie.Update(keybuf.Bytes(), list.GetRlp(i))
	// }
	// return trie.Hash()

	//直接将transactionlist进行hash
	return rlpHash(list)
}
