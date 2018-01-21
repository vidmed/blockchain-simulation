package simulator

import (
	"time"

	"encoding/json"
	"os"

	"sync"

	"github.com/pkg/errors"
	"github.com/vidmed/logger"
)

// ErrNoTransactions means that there are no transactions in block
// so the block is not flushed on disk.
var ErrNoTransactions = errors.New("no transactions")

type Simulator interface {
	// Close shuts down the simulator and waits for block data to be
	// flushed. You must call this function when you don`t want to send
	// transactions any more.
	// You have to check that you don`t send transactions into Input chanel
	// after closing simulator. This may lead to writing to closed chanel
	Close()
	// Input is the input channel for the user to add transactions to the block.
	// If you closed Simulator do not send transactions into this chanel.
	Input() chan<- *Transaction
}

// implementation of Simulator interface
type simulator struct {
	flushPeriod time.Duration
	flushFile   string
	block       *Block
	input       chan *Transaction
	done        chan struct{}

	wg sync.WaitGroup
}

// NewSimulator creates new Simulator with given parameters
func NewSimulator(flushPeriod int, flushFile string) Simulator {
	s := &simulator{
		flushPeriod: time.Duration(flushPeriod) * time.Second,
		flushFile:   flushFile,
		input:       make(chan *Transaction),
		done:        make(chan struct{}),
		block:       NewBlock(""),
	}

	s.wg.Add(1)
	go s.start()

	return s
}

func (s *simulator) Input() chan<- *Transaction {
	return s.input
}

func (s *simulator) Close() {
	close(s.done)
	s.wg.Wait()
	close(s.input)
}

func (s *simulator) start() {
	sendTicker := time.NewTicker(s.flushPeriod)
	for {
		select {
		case t := <-s.input:
			logger.Get().Infof("Simulator got new Transaction: %v", t)
			s.block.Transactions = append(s.block.Transactions, t)
		case <-sendTicker.C:
			sendTicker.Stop()
			logger.Get().Infoln("Simulator ticker fired")
			s.flushBlock()
			// reset ticker
			sendTicker = time.NewTicker(s.flushPeriod)
		case <-s.done:
			sendTicker.Stop()
			logger.Get().Infoln("Stopping Simulator")
			s.flushBlock()
			s.wg.Done()
			return
		}
	}
}

func (s *simulator) flushBlock() {
	err := s.flush()
	switch err {
	case nil:
		s.block = s.block.Next()
	case ErrNoTransactions:
		// note in case if it is not allowed to save block without transaction
		logger.Get().Warningln("There are no transactions in block while flush. Block haven`t been wrote on disk")
		return
	default: // other error
		logger.Get().Errorf("Simulator flush error: %s", err.Error())
		s.block = s.block.Next()
		// todo try to save data in some way
	}
}

// Flush method encodes block data to json and writes it to file s.flushFile.
// If file doesn`t exist it will be created. The new block data will be appended to the end of the file.
func (s *simulator) flush() error {
	// if there is no transactions in block - do not flush
	if len(s.block.Transactions) == 0 {
		return ErrNoTransactions
	}
	data, err := json.Marshal(s.block)
	if err != nil {
		return errors.Wrap(err, "Error while marshalling block data")
	}
	data = append(data, "\n"...)

	f, err := os.OpenFile(s.flushFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.Wrap(err, "Error while opening file")
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return errors.Wrap(err, "Error while writing (closing) file")
}
