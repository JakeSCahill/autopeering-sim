package discover

import (
	"log"
	"testing"
	"time"

	"github.com/iotaledger/autopeering-sim/peer"
	"github.com/iotaledger/autopeering-sim/server"
	"github.com/iotaledger/autopeering-sim/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const graceTime = 100 * time.Millisecond

var logger *zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot initialize logger: %v", err)
	}
	logger = l.Sugar()

	// use much shorter intervals to speed up the tests
	reverifyInterval = 500 * time.Millisecond
	queryInterval = 1000 * time.Millisecond
}

// newTest creates a new discovery server and also returns the teardown.
func newTest(t require.TestingT, trans transport.Transport, logger *zap.SugaredLogger, masters ...*peer.Peer) (*server.Server, *Protocol, func()) {
	log := logger.Named(trans.LocalAddr().String())
	db := peer.NewMemoryDB(log.Named("db"))
	local, err := peer.NewLocal(trans.LocalAddr().Network(), trans.LocalAddr().String(), db)
	require.NoError(t, err)

	cfg := Config{
		Log:         log,
		MasterPeers: masters,
	}
	prot := New(local, cfg)
	srv := server.Listen(local, trans, log.Named("srv"), prot)
	prot.Start(srv)

	teardown := func() {
		srv.Close()
		prot.Close()
		db.Close()
	}
	return srv, prot, teardown
}

func getPeer(s *server.Server) *peer.Peer {
	return &s.Local().Peer
}

func TestProtVerifyMaster(t *testing.T) {
	p2p := transport.P2P()
	defer p2p.Close()

	srvA, _, closeA := newTest(t, p2p.A, logger)
	defer closeA()
	peerA := getPeer(srvA)

	// use peerA as masters peer
	_, protB, closeB := newTest(t, p2p.B, logger, peerA)

	time.Sleep(graceTime) // wait for the packages to ripple through the network
	closeB()              // close srvB to avoid race conditions, when asserting

	if assert.EqualValues(t, 1, len(protB.mgr.known)) {
		assert.EqualValues(t, peerA, &protB.mgr.known[0].Peer)
		assert.EqualValues(t, 1, protB.mgr.known[0].verifiedCount)
	}
}

func TestProtPingPong(t *testing.T) {
	p2p := transport.P2P()
	defer p2p.Close()

	srvA, protA, closeA := newTest(t, p2p.A, logger)
	defer closeA()
	srvB, protB, closeB := newTest(t, p2p.B, logger)
	defer closeB()

	peerA := getPeer(srvA)
	peerB := getPeer(srvB)

	// send a ping from node A to B
	t.Run("A->B", func(t *testing.T) { assert.NoError(t, protA.ping(peerB)) })
	time.Sleep(graceTime)

	// send a ping from node B to A
	t.Run("B->A", func(t *testing.T) { assert.NoError(t, protB.ping(peerA)) })
	time.Sleep(graceTime)
}

func TestProtPingTimeout(t *testing.T) {
	p2p := transport.P2P()
	defer p2p.Close()

	_, protA, closeA := newTest(t, p2p.A, logger)
	defer closeA()
	srvB, _, closeB := newTest(t, p2p.B, logger)
	closeB() // close the connection right away to prevent any replies

	peerB := getPeer(srvB)

	// send a ping from node A to B
	err := protA.ping(peerB)
	assert.EqualError(t, err, server.ErrTimeout.Error())
}

func TestProtVerifiedPeers(t *testing.T) {
	p2p := transport.P2P()
	defer p2p.Close()

	_, protA, closeA := newTest(t, p2p.A, logger)
	defer closeA()
	srvB, _, closeB := newTest(t, p2p.B, logger)
	defer closeB()

	peerB := getPeer(srvB)

	// send a ping from node A to B
	assert.NoError(t, protA.ping(peerB))
	time.Sleep(graceTime)

	// protA should have peerB as the single verified peer
	assert.ElementsMatch(t, []*peer.Peer{peerB}, protA.GetVerifiedPeers())
	for _, p := range protA.GetVerifiedPeers() {
		assert.Equal(t, p, protA.GetVerifiedPeer(p.ID(), p.Address()))
	}
}

func TestProtVerifiedPeer(t *testing.T) {
	p2p := transport.P2P()
	defer p2p.Close()

	srvA, protA, closeA := newTest(t, p2p.A, logger)
	defer closeA()
	srvB, _, closeB := newTest(t, p2p.B, logger)
	defer closeB()

	peerA := getPeer(srvA)
	peerB := getPeer(srvB)

	// send a ping from node A to B
	assert.NoError(t, protA.ping(peerB))
	time.Sleep(graceTime)

	// we should have peerB as a verified peer
	assert.Equal(t, peerB, protA.GetVerifiedPeer(peerB.ID(), peerB.Address()))
	// we should not have ourself as a verified peer
	assert.Nil(t, protA.GetVerifiedPeer(peerA.ID(), peerA.Address()))
	// the address of peerB should match
	assert.Nil(t, protA.GetVerifiedPeer(peerB.ID(), ""))
}

func TestProtDiscoveryRequest(t *testing.T) {
	p2p := transport.P2P()
	defer p2p.Close()

	srvA, protA, closeA := newTest(t, p2p.A, logger)
	defer closeA()
	srvB, protB, closeB := newTest(t, p2p.B, logger)
	defer closeB()

	peerA := getPeer(srvA)
	peerB := getPeer(srvB)

	// request peers from node A
	t.Run("A->B", func(t *testing.T) {
		if ps, err := protA.discoveryRequest(peerB); assert.NoError(t, err) {
			assert.ElementsMatch(t, []*peer.Peer{peerA}, ps)
		}
	})
	// request peers from node B
	t.Run("B->A", func(t *testing.T) {
		if ps, err := protB.discoveryRequest(peerA); assert.NoError(t, err) {
			assert.ElementsMatch(t, []*peer.Peer{peerB}, ps)
		}
	})
}

func TestProtDiscovery(t *testing.T) {
	net := transport.NewNetwork("M", "A", "B", "C")
	defer net.Close()

	srvM, protM, closeM := newTest(t, net.GetTransport("M"), logger)
	defer closeM()
	time.Sleep(graceTime) // wait for the master to initialize

	srvA, protA, closeA := newTest(t, net.GetTransport("A"), logger, getPeer(srvM))
	defer closeA()
	srvB, protB, closeB := newTest(t, net.GetTransport("B"), logger, getPeer(srvM))
	defer closeB()
	srvC, protC, closeC := newTest(t, net.GetTransport("C"), logger, getPeer(srvM))
	defer closeC()

	time.Sleep(queryInterval + graceTime)    // wait for the next discovery cycle
	time.Sleep(reverifyInterval + graceTime) // wait for the next verification cycle

	// now the full network should be discovered
	assert.ElementsMatch(t, []*peer.Peer{getPeer(srvA), getPeer(srvB), getPeer(srvC)}, protM.GetVerifiedPeers())
	assert.ElementsMatch(t, []*peer.Peer{getPeer(srvM), getPeer(srvB), getPeer(srvC)}, protA.GetVerifiedPeers())
	assert.ElementsMatch(t, []*peer.Peer{getPeer(srvM), getPeer(srvA), getPeer(srvC)}, protB.GetVerifiedPeers())
	assert.ElementsMatch(t, []*peer.Peer{getPeer(srvM), getPeer(srvA), getPeer(srvB)}, protC.GetVerifiedPeers())
}

func BenchmarkPingPong(b *testing.B) {
	p2p := transport.P2P()
	defer p2p.Close()
	log := zap.NewNop().Sugar() // disable logging

	_, protA, closeA := newTest(b, p2p.A, log)
	defer closeA()
	srvB, _, closeB := newTest(b, p2p.B, log)
	defer closeB()

	peerB := getPeer(srvB)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		// send a ping from node A to B
		_ = protA.ping(peerB)
	}

	b.StopTimer()
}

func BenchmarkDiscoveryRequest(b *testing.B) {
	p2p := transport.P2P()
	defer p2p.Close()
	log := zap.NewNop().Sugar() // disable logging

	_, protA, closeA := newTest(b, p2p.A, log)
	defer closeA()
	srvB, _, closeB := newTest(b, p2p.B, log)
	defer closeB()

	peerB := getPeer(srvB)

	// send initial request to ensure that every peer is verified
	_, err := protA.discoveryRequest(peerB)
	require.NoError(b, err)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = protA.discoveryRequest(peerB)
	}

	b.StopTimer()
}
