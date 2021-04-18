package system_test

import (
	"os"
	"syscall"
	"testing"

	"{{[ .Project ]}}/pkg/system"

	"go.uber.org/zap"
)

const (
	testSignal                               = syscall.SIGUSR2
	customSignalType       system.SignalType = 777
	customSignalTypeString                   = "777"
)

// testHandling implement simplest Operator interface.
type testHandling struct {
	ch chan system.SignalType
}

// Reload implementation.
func (th testHandling) Reload() error {
	th.ch <- system.Reload

	return nil
}

// Maintenance implementation.
func (th testHandling) Maintenance() error {
	th.ch <- system.Maintenance

	return nil
}

// Shutdown implementation.
func (th testHandling) Shutdown() []error {
	var errs []error
	th.ch <- system.Shutdown

	return errs
}

func TestSignals(t *testing.T) {
	// Setup logger
	logger := zap.NewExample()

	defer func(*zap.Logger) {
		if err := logger.Sync(); err != nil {
			// Usually here are stdout/stderr errors for sync operations which are unsupported for it
			logger.Debug(err.Error())
		}
	}(logger)

	pid := os.Getpid()
	proc, err := os.FindProcess(pid)

	if err != nil {
		t.Error("Finding process:", err)
	}

	signals := system.NewSignals()

	shutdownSignals := signals.Get(system.Shutdown)
	verifySignal(t, syscall.SIGTERM, shutdownSignals, system.Shutdown)
	verifySignal(t, syscall.SIGINT, shutdownSignals, system.Shutdown)

	reloadSignals := signals.Get(system.Reload)
	verifySignal(t, syscall.SIGHUP, reloadSignals, system.Reload)

	maintenanceSignals := signals.Get(system.Maintenance)
	verifySignal(t, syscall.SIGUSR1, maintenanceSignals, system.Maintenance)

	ignoredSignals := signals.Get(system.Ignore)
	verifySignal(t, syscall.SIGURG, ignoredSignals, system.Ignore)

	handling := &testHandling{ch: make(chan system.SignalType, 1)}

	go func() {
		if err := signals.Wait(logger, handling); err != nil {
			t.Error("Waiting signal:", err)
		}
	}()

	// Prepare and send reload signal.
	signals.Add(testSignal, system.Reload)
	sendSignal(t, handling.ch, proc, system.Reload)
	signals.Remove(testSignal, system.Reload)

	// Prepare and send maintenance signal.
	signals.Add(testSignal, system.Maintenance)
	sendSignal(t, handling.ch, proc, system.Maintenance)
	signals.Remove(testSignal, system.Maintenance)

	// Prepare and send shutdown signal.
	signals.Add(testSignal, system.Shutdown)
	sendSignal(t, handling.ch, proc, system.Shutdown)
	signals.Remove(testSignal, system.Shutdown)
}

func sendSignal(t *testing.T, ch <-chan system.SignalType, proc *os.Process, signal system.SignalType) {
	err := proc.Signal(testSignal)
	if err != nil {
		t.Error("Sending signal:", err)

		return
	}

	if sig := <-ch; sig != signal {
		t.Error("Expected signal:", signal, "got", sig)
	}
}

func verifySignal(t *testing.T, signal os.Signal, signals []os.Signal, sigType system.SignalType) {
	if !isSignalAvailable(signal, signals) {
		t.Error("Absent of the signal:", signal, "among", sigType, "signal type")
	}
}

func TestSignalStringer(t *testing.T) {
	s := system.Shutdown

	if s.String() != "SHUTDOWN" {
		t.Error("Expected signal type SHUTDOWN, got", s.String())
	}

	s = system.Reload

	if s.String() != "RELOAD" {
		t.Error("Expected signal type RELOAD, got", s.String())
	}

	s = system.Maintenance

	if s.String() != "MAINTENANCE" {
		t.Error("Expected signal type MAINTENANCE, got", s.String())
	}

	s = system.Ignore

	if s.String() != "IGNORE" {
		t.Error("Expected signal type IGNORE, got", s.String())
	}

	s = customSignalType

	if s.String() != customSignalTypeString {
		t.Error("Expected signal type ", customSignalTypeString, "got", s.String())
	}
}

func TestRemoveNotExistingSignal(t *testing.T) {
	s := system.NewSignals()
	count := len(s.Get(system.Maintenance))

	s.Remove(syscall.SIGUSR2, system.Maintenance)

	if len(s.Get(system.Maintenance)) != count {
		t.Error("Expected count of signals", count, "got", len(s.Get(system.Maintenance)))
	}
}

// Checks if a signal is available in the specified list.
func isSignalAvailable(signal os.Signal, list []os.Signal) bool {
	for _, s := range list {
		if s == signal {
			return true
		}
	}

	return false
}
