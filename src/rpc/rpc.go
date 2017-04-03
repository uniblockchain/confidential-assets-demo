// Copyright (c) 2017 DG Lab
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

// rpc
package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ValidatedAddress struct {
	IsValid         bool   `json:"isvalid"`          // : true|false,        (boolean) If the address is valid or not. If not, this is the only property returned.
	Address         string `json:"address"`          // : "bitcoinaddress", (string) The bitcoin address validated
	ScriptPubKey    string `json:"scriptPubKey"`     // : "hex",       (string) The hex encoded scriptPubKey generated by the address
	IsMine          bool   `json:"ismine"`           // : true|false,        (boolean) If the address is yours or not
	IsWatchonly     bool   `json:"iswatchonly"`      // : true|false,   (boolean) If the address is watchonly
	IsScript        bool   `json:"isscript"`         // : true|false,      (boolean) If the key is a script
	PubKey          string `json:"pubkey"`           // : "publickeyhex",    (string) The hex value of the raw public key
	IsCompressed    bool   `json:"iscompressed"`     // : true|false,  (boolean) If the address is compressed
	Account         string `json:"account"`          // : "account"         (string) DEPRECATED. The account associated with the address, "" is the default account
	ConfidentialKey string `json:"confidential_key"` // : "pubkey" (string) The confidentiality key associated with the address, or "" if none
	Unconfidential  string `json:"unconfidential"`   // : "address"  (string) The address without confidentiality key
	Confidential    string `json:"confidential"`     // : "address"    (string) Confidential version of the address, only if it is yours and unconfidential
	Hdkeypath       string `json:"hdkeypath"`        // : "keypath"       (string, optional) The HD keypath if the key is HD and available
	Hdmasterkeyid   string `json:"hdmasterkeyid"`    // : "<hash160>" (string, optional) The Hash160 of the HD master pubkey
}

type Unspent struct {
	Txid            string `json:"txid"`            // "txid"            : "txid",    (string)  the transaction id
	Vout            int64  `json:"vout"`            // "vout"            : n,         (numeric) the vout value
	Address         string `json:"address"`         // "address"         : "address", (string)  the bitcoin address
	Account         string `json:"account"`         // "account"         : "account", (string)  DEPRECATED. The associated account, or "" for the default account
	ScriptPubKey    string `json:"scriptPubKey"`    // "scriptPubKey"    : "key",     (string)  the script key
	Amount          int64  `json:"amount"`          // "amount"          : x.xxx,     (numeric) the transaction amount in BTC
	Asset           string `json:"asset"`           // "asset"           : "hex"      (string)  the asset id for this output
	AssetCommitment string `json:"assetcommitment"` // "assetcommitment" : "hex"      (string)  the asset commitment for this output
	Confirmations   int64  `json:"confirmations"`   // "confirmations"   : n,         (numeric) The number of confirmations
	SerValue        string `json:"serValue"`        // "serValue"        : "hex",     (string)  the output's value commitment
	Blinder         string `json:"blinder"`         // "blinder"         : "blind"    (string)  The blinding factor used for a confidential output (or "")
	RedeemScript    string `json:"redeemScript"`    // "redeemScript"    : n          (string)  The redeemScript if scriptPubKey is P2SH
	Spendable       bool   `json:"spendable"`       // "spendable"       : xxx,       (bool)    Whether we have the private keys to spend this output
	Solvable        bool   `json:"solvable"`        // "solvable"        : xxx        (bool)    Whether we know how to spend this output, ignoring the lack of keys
}

type UnspentList []*Unspent

type Balance map[string]float64

type Wallet struct {
	WalletVersion      int64   `json:"walletversion"`       // : xxxxx,       (numeric) the wallet version
	Balance            Balance `json:"balance"`             // : xxxxxxx,     (numeric) the total confirmed balance of the wallet in BTC
	UnconfirmedBalance Balance `json:"unconfirmed_balance"` // : xxx,         (numeric) the total unconfirmed balance of the wallet in BTC
	ImmatureBalance    Balance `json:"immature_balance"`    // : xxxxxx,      (numeric) the total immature balance of the wallet in BTC
	TxCount            int64   `json:"txcount"`             // : xxxxxxx,     (numeric) the total number of transactions in the wallet
	KeypoolOldest      int64   `json:"keypoololdest"`       // : xxxxxx,      (numeric) the timestamp (seconds since GMT epoch) of the oldest pre-generated key in the key pool
	KeypoolSize        int64   `json:"keypoolsize"`         // : xxxx,        (numeric) how many new keys are pre-generated
	UnlockedUntil      int64   `json:"unlocked_until"`      // : ttt,         (numeric) the timestamp in seconds since epoch (midnight Jan 1 1970 GMT) that the wallet is unlocked for transfers, or 0 if the wallet is locked
	PayTxFee           int64   `json:"paytxfee"`            // : x.xxxx,      (numeric) the transaction fee configuration, set in BTC/kB
	HDMasterKeyId      string  `json:"hdmasterkeyid"`       // : "<hash160>", (string) the Hash160 of the HD master pubkey
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int64    `json:"reqSigs"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}

type Vin struct {
	Txid        string    `json:"txid"`
	Vout        int64     `json:"vout"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	txinwitness string    `json:"txinwitness"`
	sequence    int64     `json:"sequence"`
}

type Vout struct {
	Value        float64      `json:"value"`
	N            int64        `json:"n"`
	Asset        string       `json:"asset"`
	Assettag     string       `json:"assettag"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type RawTransaction struct {
	Txid     string  `json:"txid"`
	Hash     string  `json:"hash"`
	Size     int64   `json:"size"`
	Vsize    int64   `json:"vsize"`
	Version  int64   `json:"version"`
	LockTime int64   `json:"locktime"`
	Fee      float64 `json:"fee"`
	Vin      []Vin   `json:"vin"`
	Vout     []Vout  `json:"vout"`
}

type SignedTransaction struct {
	Hex      string `json:"hex"`
	Complete bool   `json:"complete"`
}

type Rpc struct {
	Url  string
	User string
	Pass string
	View bool
}

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc,"`
	Id      string        `json:"id,"`
	Method  string        `json:"method,"`
	Params  []interface{} `json:"params,"`
}

type RpcResponse struct {
	Result interface{} `json:"result,"`
	Error  interface{} `json:"error,"`
	Id     string      `json:"id,"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (res *RpcResponse) UnmarshalError() (RpcError, error) {
	var rerr RpcError
	if res.Error == nil {
		return rerr, fmt.Errorf("RpcResponse Error is nil.")
	}
	data, ok := res.Error.(map[string]interface{})
	if !ok {
		return rerr, fmt.Errorf("RpcResponse Error is not map[string]interface{}")
	}
	bs, _ := json.Marshal(data)
	json.Unmarshal(bs, &rerr)
	return rerr, nil
}

func (res *RpcResponse) UnmarshalResult(result interface{}) error {
	if res.Result == nil {
		return fmt.Errorf("RpcResponse Result is nil.")
	}
	var bs []byte
	m, ok := res.Result.(map[string]interface{})
	if !ok {
		arr, ok := res.Result.([]interface{})
		if !ok {
			return fmt.Errorf("RpcResponse Result is neither map[string]interface{} nor []interface{}")
		} else {
			bs, _ = json.Marshal(arr)
		}
	} else {
		bs, _ = json.Marshal(m)
	}
	err := json.Unmarshal(bs, result)
	if err != nil {
		return err
	}
	return nil
}

func NewRpc(url, user, pass string) *Rpc {
	rpc := new(Rpc)
	rpc.Url = url
	rpc.User = user
	rpc.Pass = pass
	return rpc
}

func (rpc *Rpc) Request(method string, params ...interface{}) (RpcResponse, error) {
	var res RpcResponse
	if len(params) == 0 {
		params = []interface{}{}
	}
	id := fmt.Sprintf("%d", time.Now().Unix())
	req := &RpcRequest{"1.0", id, method, params}
	bs, _ := json.Marshal(req)
	if rpc.View {
		fmt.Printf("%s\n", bs)
	}
	client := &http.Client{}
	hreq, _ := http.NewRequest("POST", rpc.Url, bytes.NewBuffer(bs))
	hreq.SetBasicAuth(rpc.User, rpc.Pass)
	hres, err := client.Do(hreq)
	if err != nil {
		return res, err
	}
	defer hres.Body.Close()
	body, _ := ioutil.ReadAll(hres.Body)
	if rpc.View {
		fmt.Printf("%d, %s\n", hres.StatusCode, body)
	}
	err = json.Unmarshal(body, &res)
	if err != nil || hres.StatusCode != http.StatusOK || res.Id != id {
		return res, fmt.Errorf("status:%v, error:%v, body:%s reqid:%v, resid:%v", hres.Status, err, body, id, res.Id)
	}
	return res, nil
}

func (rpc *Rpc) RequestAndUnmarshalResult(result interface{}, method string, params ...interface{}) (RpcResponse, error) {
	res, err := rpc.Request(method, params...)
	if err != nil {
		return res, err
	}
	err = res.UnmarshalResult(result)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (rpc *Rpc) RequestAndCastNumber(method string, params ...interface{}) (float64, RpcResponse, error) {
	var num float64
	res, err := rpc.Request(method, params...)
	if err != nil {
		return num, res, err
	}
	num, ok := res.Result.(float64)
	if !ok {
		return num, res, fmt.Errorf("RpcResponse Result cast error:%+v", res.Result)
	}
	return num, res, nil
}

func (rpc *Rpc) RequestAndCastString(method string, params ...interface{}) (string, RpcResponse, error) {
	var str string
	res, err := rpc.Request(method, params...)
	if err != nil {
		return str, res, err
	}
	str, ok := res.Result.(string)
	if !ok {
		return str, res, fmt.Errorf("RpcResponse Result cast error:%+v", res.Result)
	}
	return str, res, nil
}

func (rpc *Rpc) RequestAndCastBool(method string, params ...interface{}) (bool, RpcResponse, error) {
	var b bool
	res, err := rpc.Request(method, params...)
	if err != nil {
		return b, res, err
	}
	b, ok := res.Result.(bool)
	if !ok {
		return b, res, fmt.Errorf("RpcResponse Result cast error:%+v", res.Result)
	}
	return b, res, nil
}
