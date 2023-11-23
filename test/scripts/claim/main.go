package main

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-bridge-service/etherman"
	clientUtils "github.com/0xPolygonHermez/zkevm-bridge-service/test/client"
	"github.com/0xPolygonHermez/zkevm-bridge-service/utils"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
)

const (
	l2BridgeAddr = "0x10B65c586f795aF3eCCEe594fE4E38E1F059F780"

	l2AccHexAddress    = "0x2ecf31ece36ccac2d3222a303b1409233ecbb225"
	l2AccHexPrivateKey = "0xdfd01798f92667dbf91df722434e8fbe96af0211d4d1b82bbbbc8f1def7a814f"
	l2NetworkURL       = "http://localhost:8123"
	bridgeURL          = "http://localhost:8080"

	mtHeight = 32
)

func main() {
	ctx := context.Background()
	c, err := utils.NewClient(ctx, l2NetworkURL, common.HexToAddress(l2BridgeAddr))
	if err != nil {
		log.Fatal("Error: ", err)
	}
	auth, err := c.GetSigner(ctx, l2AccHexPrivateKey)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Get Claim data
	cfg := clientUtils.Config{
		L1NodeURL:    l2NetworkURL,
		L2NodeURL:    l2NetworkURL,
		BridgeURL:    bridgeURL,
		L2BridgeAddr: common.HexToAddress(l2BridgeAddr),
	}
	client, err := clientUtils.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	deposits, _, err := client.GetBridges(l2AccHexAddress, 0, 10) //nolint
	if err != nil {
		log.Fatal("Error: ", err)
	}
	bridgeData := deposits[0]
	proof, err := client.GetMerkleProof(deposits[0].NetworkId, deposits[0].DepositCnt)
	if err != nil {
		log.Fatal("error: ", err)
	}
	log.Debug("bridge: ", bridgeData)
	log.Debug("mainnetExitRoot: ", proof.MainExitRoot)
	log.Debug("rollupExitRoot: ", proof.RollupExitRoot)

	var smt [mtHeight][32]byte
	for i := 0; i < len(proof.MerkleProof); i++ {
		log.Debug("smt: ", proof.MerkleProof[i])
		smt[i] = common.HexToHash(proof.MerkleProof[i])
	}
	globalExitRoot := &etherman.GlobalExitRoot{
		ExitRoots: []common.Hash{common.HexToHash(proof.MainExitRoot), common.HexToHash(proof.RollupExitRoot)},
	}
	log.Info("Sending claim tx...")
	err = c.SendClaimAndWait(ctx, bridgeData, smt, globalExitRoot, auth)
	if err != nil {
		log.Fatal("error: ", err)
	}
	log.Info("Success!")
	balance, err := c.Client.BalanceAt(ctx, common.HexToAddress(l2AccHexAddress), nil)
	if err != nil {
		log.Fatal("error getting balance: ", err)
	}
	log.Info("L2 balance: ", balance)
}
