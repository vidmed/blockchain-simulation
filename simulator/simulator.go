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

// Simulator represents a blockchain simulator with adding transactions to a bloack
// and flushing it on a disk.
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
	flushPeriod     time.Duration
	maxTransactions uint
	flushFile       string
	block           *block
	input           chan *Transaction
	done            chan struct{}

	wg sync.WaitGroup
}

// NewSimulator creates new Simulator with given parameters
func NewSimulator(flushPeriod, maxTransactions uint, flushFile string) Simulator {
	s := &simulator{
		flushPeriod:     time.Duration(flushPeriod) * time.Second,
		maxTransactions: maxTransactions,
		flushFile:       flushFile,
		input:           make(chan *Transaction),
		done:            make(chan struct{}),
		block:           newBlock(""),
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
	c := uint(0)
	for {
		select {
		case t := <-s.input:
			logger.Get().Infof("Simulator got new Transaction: %v", t)
			s.block.Transactions = append(s.block.Transactions, t)
			c++
			if c == s.maxTransactions {
				sendTicker.Stop()
				logger.Get().Infof("Number of transaction reached maximum value(%d). Flushing", s.maxTransactions)
				s.flushBlock()
				c = 0
				// reset ticker
				sendTicker = time.NewTicker(s.flushPeriod)
			}
		case <-sendTicker.C:
			sendTicker.Stop()
			logger.Get().Infoln("Simulator ticker fired")
			s.flushBlock()
			c = 0
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
		s.block = s.block.next()
	case ErrNoTransactions:
		// note in case if it is not allowed to save block without transaction
		logger.Get().Warningln("There are no transactions in block while flush. Block haven`t been wrote on disk")
		return
	default: // other error
		logger.Get().Errorf("Simulator flush error: %s", err.Error())
		s.block = s.block.next()
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
