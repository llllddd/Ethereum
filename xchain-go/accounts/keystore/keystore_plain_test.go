package keystore

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"xchain-go/common"
)

const (
	veryLightScryptN = 1 << 12
	veryLightScryptP = 6
)

/*
返回目录,可以keyStore类型
1.生成临时目录
2.是否加密的选项返回对应的keystore
*/

func temKeyStoreIface(t *testing.T, encrypted bool) (d string, ks keyStore) {
	//在目录/tmp生成临时目录
	dir, err := ioutil.TempDir("", "keystore-test")
	if err != nil {
		t.Fatal(err)
	}
	if encrypted == false {
		return dir, keyStorePlain{dir}
	} else {
		return dir, keyStorePassphrase{dir, veryLightScryptN, veryLightScryptP}
	}
}

/*
1.得到临时目录和keystore
2.存储密钥
3.得到密钥
4.比较
*/

func TestKeyStorePlain(t *testing.T) {
	dir, ks := temKeyStoreIface(t, false)
	defer os.RemoveAll(dir)

	pass := ""
	k1, account, err := storeNewKey(ks, rand.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}
	k2, err := ks.GetKey(k1.Address, account.URL.Path, pass)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.Address, k2.Address) {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.PrivateKey, k2.PrivateKey) {
		t.Fatal(err)
	}

}

func TestKeyStorePassphrase(t *testing.T) {
	dir, ks := temKeyStoreIface(t, true)
	defer os.RemoveAll(dir)

	pass := "test"
	k1, Account, err := storeNewKey(ks, rand.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}

	k2, err := ks.GetKey(k1.Address, Account.URL.Path, pass)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.Address, k2.Address) {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(k1.PrivateKey, k2.PrivateKey) {
		t.Fatal(err)
	}
}

/*
func TestKeyStorePassphraseDecryptFail(t *testing.T) {
	dir, ks := temKeyStoreIface(t, true)
	defer os.RemoveAll(dir)

	pass := "test"
	k1, account, err := storeNewKey(ks, rand.Reader, pass)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ks.GetKey(k1.Address, account.URL.Path, "foo")
	if err != nil {
		t.Fatalf("用户输入的密码错误:\n输入的为%q,应该为%q", err, ErrDecrypt)
	}
}
*/
type KeyStoreTest struct {
	JSON     encryptedKeyJSON
	Password string
	Priv     string
}

func Test_Scrypt_1(t *testing.T) {
	t.Parallel()
	tests := loadKeyStoreTest("testdata/v3_test_vector.json", t)
	testDecrypt(tests["wikipage_test_vector_scrypt"], t)
}

func loadKeyStoreTest(file string, t *testing.T) map[string]KeyStoreTest {
	tests := make(map[string]KeyStoreTest)
	err := common.LoadJSON(file, &tests)
	if err != nil {
		t.Fatal(err)
	}
	return tests
}

func testDecrypt(keystore KeyStoreTest, t *testing.T) {
	priv, _, err := decryptKey(&keystore.JSON, keystore.Password)
	if err != nil {
		t.Fatal(err)
	}
	privHex := hex.EncodeToString(priv)
	if keystore.Priv != privHex {
		t.Fatal(fmt.Errorf("测试密钥与期望密钥不同: %v,%v", keystore.Priv, privHex))
	}
}
