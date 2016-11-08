package chains

import (
	"fmt"

	"github.com/eris-ltd/eris-cli/config"
	"github.com/eris-ltd/eris-cli/definitions"
	"github.com/eris-ltd/eris-cli/log"
	"github.com/eris-ltd/eris-cli/maker"
	"github.com/eris-ltd/eris-cli/services"
	"github.com/eris-ltd/eris-cli/util"

	keys "github.com/eris-ltd/eris-keys/eris-keys"
)

// MakeChain runs the `eris-cm make` command in a Docker container.
// It returns an error. Note that if do.Known, do.AccountTypes
// or do.ChainType are not set the command will run via interactive
// shell.
//
//  do.Name          - name of the chain to be created (required)
//  do.Known         - will use the mintgen tool to parse csv's and create a genesis.json (requires do.ChainMakeVals and do.ChainMakeActs) (optional)
//  do.Output  	     - outputs the jobs_output.yaml  (default true) XXX [zr] we can probably eliminate?
//  do.ChainMakeVals - csv file to use for validators (optional)
//  do.ChainMakeActs - csv file to use for accounts (optional)
//  do.AccountTypes  - use eris-cm make account-types paradigm (example: Root:1,Participants:25,...) (optional)
//  do.ChainType     - use eris-cm make chain-types paradigm (example: simplechain) (optional)
//  do.Tarball       - instead of outputing raw files in directories, output packages of tarbals (optional)
//  do.ZipFile       - similar to do.Tarball except uses zipfiles (optional)
//  do.Verbose       - verbose output (optional)
//  do.Debug         - debug output (optional)
//
func MakeChain(do *definitions.Do) error {
	doKeys := definitions.NowDo()
	doKeys.Name = "keys"
	if err := services.EnsureRunning(doKeys); err != nil {
		return err
	}

	// loop through chains directories to make sure they exist & are appropriately populated
	//for _, d := range ChainsDirs {
	//	if _, err := os.Stat(d); os.IsNotExist(err) {
	//		os.MkdirAll(d, 0755)
	//	}
	//}
	if err := util.CheckDefaultTypes(config.AccountsTypePath, "account-types"); err != nil {
		return err
	}
	if err := util.CheckDefaultTypes(config.ChainTypePath, "chain-types"); err != nil {
		return err
	}

	// announce.
	log.Info("Hello! I'm the marmot who makes eris chains.")
	keys.DaemonAddr = "http://172.17.0.2:4767" // tmp

	if do.Known {
		do.CSV = fmt.Sprintf("%s,%s", do.ChainMakeVals, do.ChainMakeActs)
	}

	// make it
	if err := maker.MakeChain(do); err != nil {
		return err
	}

	// cm currently is not opinionated about its writers.
	if do.Tarball {
		if err := util.Tarball(do); err != nil {
			return err
		}
	} else if do.ZipFile {
		if err := util.Zip(do); err != nil {
			return err
		}
	}
	if do.Output {
		if err := util.SaveAccountResults(do); err != nil {
			return err
		}
	}

	return nil
}
