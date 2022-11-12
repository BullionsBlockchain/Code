package evmc

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/evmc/v10/bindings/go/evmc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tests"
	"github.com/nsf/jsondiff"
)

type evmcVM struct {
	cap  evmc.Capability
	path string
}

type testResult struct {
	err      error
	dump     state.Dump
	dumpJSON []byte
}

var (
	// Configure the tested external EVMs
	soEVMOne = evmcVM{
		evmc.CapabilityEVM1,
		"../build/_workspace/evmone/lib/libevmone.so",
	}
	soHera = evmcVM{
		evmc.CapabilityEWASM,
		"../build/_workspace/hera/build/src/libhera.so",
	}
)

var (
	// Define fork(s) to run
	forks = []string{"Istanbul"}
)

func TestHeraVEVMOne(t *testing.T) {

	// Configure the go-ethereum logger
	verbose := false
	if verbose {
		glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
		glogger.Verbosity(log.Lvl(log.LvlDebug))
		log.Root().SetHandler(glogger)
	}

	// Define the test path(s)
	// myTestFiles :=
	myTestFiles := []string{
		// "../tests/testdata/GeneralStateTests/VMTests/vmArithmeticTest/add.json",
		"../tests/testdata/GeneralStateTests/stRandom/randomStatetest153.json",
		// "../tests/testdata/GeneralStateTests/stBadOpcode/operationDiffGas.json",
		// "../tests/testdata/GeneralStateTests/stCallCodes/callcallcallcode_001_SuicideMiddle.json",
	}

	for _, myTestFile := range myTestFiles {
		runEVMCStateTestFile(t, []evmcVM{soEVMOne, soHera}, myTestFile)
	}
}

func runEVMCStateTestFile(t *testing.T, myEVMs []evmcVM, testFile string) {
	// Load the test content from the input file
	src, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	var mytests map[string]tests.StateTest
	if err = json.Unmarshal(src, &mytests); err != nil {
		t.Fatal(err)
	}

	// Run the tests
	t.Log("====== Running tests from", testFile)
	for _, mytest := range mytests {

		evmoneResults := runEVMCStateTest(t, soEVMOne, mytest, forks)
		heraResults := runEVMCStateTest(t, soHera, mytest, forks)

		// Compare output
		for i, evmoneResult := range evmoneResults {
			didErr := false

			if evmoneResult.err != nil {
				didErr = true
				t.Errorf("EVMOne error: %v", evmoneResult.err)
			}

			heraResult := heraResults[i]
			if heraResult.err != nil {
				didErr = true
				t.Errorf("Hera error: %v", heraResult.err)
			}

			opts := jsondiff.DefaultConsoleOptions()
			diff, diffStr := jsondiff.Compare(heraResult.dumpJSON, evmoneResult.dumpJSON, &opts)

			if didErr || diff != jsondiff.FullMatch {
				t.Log("JSON diff:", diff)
				t.Log(diffStr)
			}
		}
	}
}

func runEVMCStateTest(t *testing.T, myvm evmcVM, test tests.StateTest, forks []string) []testResult {
	loggerConfig := &logger.Config{
		EnableMemory:     false,
		DisableStack:     false,
		DisableStorage:   false,
		EnableReturnData: true,
	}
	var (
		tracer   vm.EVMLogger
		debugger *logger.StructLogger
	)
	debugger = logger.NewStructLogger(loggerConfig)
	tracer = debugger

	// Iterate over all the tests, run them and aggregate the results
	cfg := vm.Config{
		Tracer: tracer,
		Debug:  true,
	}
	if myvm.cap == evmc.CapabilityEVM1 {
		cfg.EVMInterpreter = myvm.path
		vm.InitEVMCEVM(cfg.EVMInterpreter)
	} else if myvm.cap == evmc.CapabilityEWASM {
		cfg.EWASMInterpreter = myvm.path
		vm.InitEVMCEwasm(cfg.EWASMInterpreter)
	}

	results := []testResult{}
	for _, st := range test.Subtests(nil) {
		res := testResult{}
		if len(forks) > 0 {
			forkMatched := false
			for _, fork := range forks {
				if strings.Contains(fork, st.Fork) {
					forkMatched = true
					break
				}
			}
			if !forkMatched {
				continue
			}
		}

		_, s, err := test.Run(st, cfg, false)
		res.err = err
		// if err != nil {
		// 	// Test failed, mark as so and dump any state to aid debugging
		// 	// result.Pass, result.Error = false, err.Error()
		// }
		if s != nil {
			res.dump = s.RawDump(nil)
			res.dumpJSON, _ = json.MarshalIndent(res.dump, "", "  ")
		}
		results = append(results, res)
	}
	return results
}
