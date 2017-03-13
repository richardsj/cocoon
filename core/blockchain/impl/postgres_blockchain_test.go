package impl

import (
	"testing"

	"fmt"

	"strings"

	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // gorm requires it
	"github.com/ncodes/cocoon/core/types"
	"github.com/ncodes/cocoon/core/types/blockchain"
	"github.com/ncodes/cocoon/core/types/store"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPosgresBlockchain(t *testing.T) {
	Convey("PostgresBlockchain", t, func() {

		var conStr = "host=localhost user=ned dbname=cocoon-dev sslmode=disable password="
		pgChain := new(PostgresBlockchain)
		db, err := pgChain.Connect(conStr)
		fmt.Println(err)
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)

		var RestDB = func() {
			db.(*gorm.DB).DropTable(ChainTableName, BlockTableName)
		}

		Convey(".Connect", func() {
			Convey("should return error when unable to connect to a postgres server", func() {
				var conStr = "host=localhost user=wrong dbname=test sslmode=disable password=abc"
				pgBlkch := new(PostgresBlockchain)
				db, err := pgBlkch.Connect(conStr)
				So(db, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "failed to connect to blockchain backend")
			})
		})

		Convey(".Init", func() {

			Convey("when chain table does not exists", func() {

				Convey("should create chain and block table and create a global chain", func() {

					chainTableExists := db.(*gorm.DB).HasTable(ChainTableName)
					So(chainTableExists, ShouldEqual, false)

					err := pgChain.Init(blockchain.GetGlobalChainName())
					So(err, ShouldBeNil)

					chainTableExists = db.(*gorm.DB).HasTable(ChainTableName)
					So(chainTableExists, ShouldEqual, true)

					chainTableExists = db.(*gorm.DB).HasTable(ChainTableName)
					So(chainTableExists, ShouldEqual, true)

					Convey("chain table must include a global chain", func() {
						var entries []blockchain.Chain
						err := db.(*gorm.DB).Find(&entries).Error
						So(err, ShouldBeNil)
						So(len(entries), ShouldEqual, 1)
						So(entries[0].Name, ShouldEqual, blockchain.GetGlobalChainName())
					})

					Reset(func() {
						RestDB()
					})
				})
			})

			Convey("when ledger table exists", func() {
				Convey("should return nil with no effect", func() {
					err := pgChain.Init(blockchain.GetGlobalChainName())
					So(err, ShouldBeNil)

					chainTableExists := db.(*gorm.DB).HasTable(ChainTableName)
					So(chainTableExists, ShouldEqual, true)

					var chains []blockchain.Chain
					err = db.(*gorm.DB).Find(&chains).Error
					So(err, ShouldBeNil)
					So(len(chains), ShouldEqual, 1)
					So(chains[0].Name, ShouldEqual, blockchain.GetGlobalChainName())
				})

				Reset(func() {
					RestDB()
				})
			})
		})

		Convey(".MakeChainName", func() {
			Convey("Should replace namespace with empty string if provided name is equal to blockchain.GetGlobalChainName()", func() {
				name := blockchain.GetGlobalChainName()
				namespace := ""
				expected := util.Sha256(fmt.Sprintf("%s.%s", namespace, name))
				actual := pgChain.MakeChainName("namespace_will_be_ignored", name)
				So(expected, ShouldEqual, actual)
			})

			Convey("Should return expected name with namespace and name hashed together", func() {
				expected := util.Sha256(fmt.Sprintf("%s.%s", "cocooncode_1", "accounts"))
				So(expected, ShouldEqual, pgChain.MakeChainName("cocooncode_1", "accounts"))
			})
		})

		Convey(".CreateChain", func() {
			err := pgChain.Init(blockchain.GetGlobalChainName())
			So(err, ShouldBeNil)

			Convey("Should successfully create a chain", func() {
				chain, err := pgChain.CreateChain("chain1", true)
				So(err, ShouldBeNil)
				So(chain.Name, ShouldEqual, "chain1")
				So(chain.Public, ShouldEqual, true)

				Convey("Should fail when trying to create a chain with an already used name", func() {
					chain, err := pgChain.CreateChain("chain1", true)
					So(chain, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "chain with matching name already exists")
				})

				Convey(".GetChain", func() {

					Convey("Should successfully return existing chain", func() {
						chain, err := pgChain.GetChain("chain1")
						So(err, ShouldBeNil)
						So(chain.Name, ShouldEqual, "chain1")
						So(chain.Public, ShouldEqual, true)
					})

					Convey("Should return nil, nil result if chain does not exists", func() {
						chain, err := pgChain.GetChain("chain2")
						So(chain, ShouldBeNil)
						So(err, ShouldBeNil)
					})
				})
			})

			Convey(".MakeTxsHash", func() {

				Convey("Should successfully return expected sha256 hash", func() {
					txs := []*store.Transaction{
						&store.Transaction{Hash: util.Sha256("a")},
						&store.Transaction{Hash: util.Sha256("b")},
						&store.Transaction{Hash: util.Sha256("c")},
					}
					hash := MakeTxsHash(txs)
					So(len(hash), ShouldEqual, 64)
					So(hash, ShouldEqual, "1bafc323731f0953ab498414f267345bdb2afffdc07e0505b5e7cddee48dfea8")
				})
			})

			Convey(".VerifyTxs", func() {

				Convey("Should successfully verify all transactions to be accurate", func() {
					tx1 := &store.Transaction{Number: 1, Ledger: "ledger1", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123456789}
					tx1.Hash = tx1.MakeHash()
					tx2 := &store.Transaction{Number: 2, Ledger: "ledger2", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123499789}
					tx2.Hash = tx2.MakeHash()
					txs := []*store.Transaction{
						tx1,
						tx2,
					}
					failedTx, verified := VerifyTxs(txs)
					So(failedTx, ShouldBeNil)
					So(verified, ShouldEqual, true)
				})

				Convey("Should fail if at least one tx hash is invalid", func() {
					tx1 := &store.Transaction{Number: 1, Ledger: "ledger1", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123456789}
					tx1.Hash = tx1.MakeHash()
					tx2 := &store.Transaction{Number: 2, Ledger: "ledger2", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123499789}
					tx2.Hash = tx2.MakeHash()
					tx3 := &store.Transaction{Number: 3, Ledger: "ledger3", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123499789}
					tx3.Hash = "very very wrong hash"
					txs := []*store.Transaction{
						tx1,
						tx2,
						tx3,
					}
					failedTx, verified := VerifyTxs(txs)
					So(failedTx, ShouldNotBeNil)
					So(verified, ShouldEqual, false)
					So(tx3, ShouldResemble, failedTx)
				})
			})

			Convey(".CreateBlock", func() {

				chain, err := pgChain.CreateChain("chain1", true)
				So(err, ShouldBeNil)
				So(chain.Name, ShouldEqual, "chain1")
				So(chain.Public, ShouldEqual, true)

				Convey("Should return error if chain does not exist", func() {
					blk, err := pgChain.CreateBlock("unknown", nil)
					So(blk, ShouldBeNil)
					So(err, ShouldEqual, types.ErrChainNotFound)
				})

				Convey("Should return error if no transaction is provided", func() {
					blk, err := pgChain.CreateBlock("chain1", []*store.Transaction{})
					So(blk, ShouldBeNil)
					So(err, ShouldEqual, types.ErrZeroTransactions)
				})

				Convey("Should return error if a transaction hash is invalid", func() {
					tx1 := &store.Transaction{Number: 1, Ledger: "ledger1", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123456789}
					tx1.Hash = tx1.MakeHash()
					tx2 := &store.Transaction{Number: 2, Ledger: "ledger2", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123499789}
					tx2.Hash = "wrong hash"
					txs := []*store.Transaction{tx1, tx2}
					blk, err := pgChain.CreateBlock("chain1", txs)
					So(blk, ShouldBeNil)
					So(err.Error(), ShouldContainSubstring, "has an invalid hash")
				})

				Convey("Should successfully create the first block with expected block values", func() {
					tx1 := &store.Transaction{Number: 1, Ledger: "ledger1", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123456789}
					tx1.Hash = tx1.MakeHash()
					tx2 := &store.Transaction{Number: 2, Ledger: "ledger2", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123499789}
					tx2.Hash = tx2.MakeHash()
					txs := []*store.Transaction{tx1, tx2}

					blk, err := pgChain.CreateBlock("chain1", txs)
					So(blk, ShouldNotBeNil)
					So(err, ShouldBeNil)
					So(blk.ChainName, ShouldEqual, "chain1")
					So(blk.Number, ShouldEqual, 1)
					So(blk.HasRightSibling, ShouldEqual, false)
					So(blk.PrevBlockHash, ShouldEqual, strings.Repeat("0", 64))
					So(blk.Hash, ShouldEqual, MakeTxsHash(txs))
					txsBytes, _ := util.ToJSON(txs)
					So(blk.Transactions, ShouldResemble, txsBytes)

					Convey("Should successfully add another block that references the previous block", func() {
						tx1 := &store.Transaction{Number: 1, Ledger: "ledger1", ID: "some_id", Key: "key", Value: "value", CreatedAt: 123456789}
						tx1.Hash = tx1.MakeHash()
						txs := []*store.Transaction{tx1}

						blk2, err := pgChain.CreateBlock("chain1", txs)
						So(blk2, ShouldNotBeNil)
						So(err, ShouldBeNil)
						So(blk2.ChainName, ShouldEqual, "chain1")
						So(blk2.Number, ShouldEqual, 2)
						So(blk2.HasRightSibling, ShouldEqual, false)
						So(blk2.PrevBlockHash, ShouldEqual, blk.Hash)
						So(blk2.Hash, ShouldEqual, MakeTxsHash(txs))
						txsBytes, _ := util.ToJSON(txs)
						So(blk2.Transactions, ShouldResemble, txsBytes)
					})
				})

				Reset(func() {
					RestDB()
				})
			})

			Reset(func() {
				RestDB()
			})
		})

		Reset(func() {
			RestDB()
		})
	})
}
