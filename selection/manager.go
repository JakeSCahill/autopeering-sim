package selection

import (
	"math/rand"
	"sync"
	"time"

	"github.com/iotaledger/autopeering-sim/peer"
	"github.com/iotaledger/autopeering-sim/salt"
	"go.uber.org/zap"
)

const (
	updateOutboundInterval = 1000 * time.Millisecond

	inboundRequestQueue = 1000
	dropQueue           = 1000

	accept = true
	reject = false
)

// A network represents the communication layer for the manager.
type network interface {
	local() *peer.Local

	RequestPeering(*peer.Peer, *salt.Salt, peer.ServiceMap) (peer.ServiceMap, error)
	DropPeer(*peer.Peer)
}

type peeringRequest struct {
	peer *peer.Peer
	salt *salt.Salt
}

type manager struct {
	net       network
	lifetime  time.Duration
	peersFunc func() []*peer.Peer

	log           *zap.SugaredLogger
	dropNeighbors bool

	inbound  *Neighborhood
	outbound *Neighborhood

	inboundRequestChan chan peeringRequest
	inboundReplyChan   chan bool
	inboundDropChan    chan peer.ID
	outboundDropChan   chan peer.ID

	rejectionFilter *Filter

	wg              sync.WaitGroup
	inboundClosing  chan struct{}
	outboundClosing chan struct{}
}

func newManager(net network, lifetime time.Duration, peersFunc func() []*peer.Peer, dropNeighbors bool, log *zap.SugaredLogger) *manager {
	m := &manager{
		net:                net,
		lifetime:           lifetime,
		peersFunc:          peersFunc,
		log:                log,
		dropNeighbors:      dropNeighbors,
		inboundClosing:     make(chan struct{}),
		outboundClosing:    make(chan struct{}),
		rejectionFilter:    NewFilter(),
		inboundRequestChan: make(chan peeringRequest, inboundRequestQueue),
		inboundReplyChan:   make(chan bool),
		inboundDropChan:    make(chan peer.ID, dropQueue),
		outboundDropChan:   make(chan peer.ID, dropQueue),
		inbound: &Neighborhood{
			Neighbors: []peer.PeerDistance{},
			Size:      4},
		outbound: &Neighborhood{
			Neighbors: []peer.PeerDistance{},
			Size:      4},
	}

	return m
}

func (m *manager) start() {
	// create valid salts
	if m.net.local().GetPublicSalt() == nil || m.net.local().GetPrivateSalt() == nil {
		m.updateSalt()
	}

	m.wg.Add(2)
	go m.loopOutbound()
	go m.loopInbound()
}

func (m *manager) close() {
	close(m.inboundClosing)
	close(m.outboundClosing)
	m.wg.Wait()
}

func (m *manager) self() peer.ID {
	return m.net.local().ID()
}

func (m *manager) loopOutbound() {
	defer m.wg.Done()

	var (
		updateOutboundDone chan struct{}
		updateOutbound     = time.NewTimer(0) // setting this to 0 will cause a trigger right away
		backoff            = 10
	)
	defer updateOutbound.Stop()

Loop:
	for {
		select {
		case <-updateOutbound.C:
			// if there is no updateOutbound, this means doUpdateOutbound is not running
			if updateOutboundDone == nil {
				updateOutboundDone = make(chan struct{})

				// check salt and update if necessary (this will drop the whole neighborhood)
				if m.net.local().GetPublicSalt().Expired() {
					m.updateSalt()
				}
				//remove potential duplicates
				dup := m.getDuplicates()
				for _, peerToDrop := range dup {
					toDrop := m.inbound.GetPeerFromID(peerToDrop)
					time.Sleep(time.Duration(rand.Intn(backoff)) * time.Millisecond)
					m.outbound.RemovePeer(peerToDrop)
					m.inbound.RemovePeer(peerToDrop)
					if toDrop != nil {
						m.dropPeer(toDrop)
					}
				}
				go m.updateOutbound(updateOutboundDone)
			}
		case <-updateOutboundDone:
			updateOutboundDone = nil
			updateOutbound.Reset(updateOutboundInterval) // updateOutbound again after the given interval
		case peerToDrop := <-m.outboundDropChan:
			if containsPeer(m.outbound.GetPeers(), peerToDrop) {
				m.outbound.RemovePeer(peerToDrop)
				m.rejectionFilter.AddPeer(peerToDrop)
				m.log.Debug("Outbound Dropped BY ", peerToDrop, " (", len(m.outbound.GetPeers()), ",", len(m.inbound.GetPeers()), ")")
			}

		// on close, exit the loop
		case <-m.outboundClosing:
			break Loop
		}
	}

	// // wait for the updateOutbound to finish
	// if updateOutboundDone != nil {
	// 	<-updateOutboundDone
	// }
}

func (m *manager) loopInbound() {
	defer m.wg.Done()

Loop:
	for {
		select {
		case req := <-m.inboundRequestChan:
			m.updateInbound(req.peer, req.salt)
		case peerToDrop := <-m.inboundDropChan:
			if containsPeer(m.inbound.GetPeers(), peerToDrop) {
				m.inbound.RemovePeer(peerToDrop)
				m.log.Debug("Inbound Dropped BY ", peerToDrop, " (", len(m.outbound.GetPeers()), ",", len(m.inbound.GetPeers()), ")")
			}

		// on close, exit the loop
		case <-m.inboundClosing:
			break Loop
		}
	}
}

// updateOutbound updates outbound neighbors.
func (m *manager) updateOutbound(done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}() // always signal, when the function returns

	// sort verified peers by distance
	distList := peer.SortBySalt(m.self().Bytes(), m.net.local().GetPublicSalt().GetBytes(), m.peersFunc())

	filter := NewFilter()
	filter.AddPeer(m.self())               //set filter for ourself
	filter.AddPeers(m.inbound.GetPeers())  // set filter for inbound neighbors
	filter.AddPeers(m.outbound.GetPeers()) // set filter for outbound neighbors

	filteredList := filter.Apply(distList)               // filter out current neighbors
	filteredList = m.rejectionFilter.Apply(filteredList) // filter out previous rejection

	// select new candidate
	candidate := m.outbound.Select(filteredList)

	if candidate.Remote != nil {
		furthest, _ := m.outbound.getFurthest()

		// send peering request
		mySalt := m.net.local().GetPublicSalt()
		myServices := m.net.local().Services()
		services, err := m.net.RequestPeering(candidate.Remote, mySalt, myServices)
		if err != nil {
			return
		}

		// add candidate to the outbound neighborhood
		if len(services) > 0 {
			//m.acceptedFilter.AddPeer(candidate.Remote.ID())
			if furthest.Remote != nil {
				m.outbound.RemovePeer(furthest.Remote.ID())
				m.dropPeer(furthest.Remote)
				m.log.Debug("Outbound furthest removed ", furthest.Remote.ID())
			}
			m.outbound.Add(candidate)
			m.log.Debug("Peering request TO ", candidate.Remote.ID(), " status ACCEPTED (", len(m.outbound.GetPeers()), ",", len(m.inbound.GetPeers()), ")")
		} else {
			m.log.Debug("Peering request TO ", candidate.Remote.ID(), " status REJECTED (", len(m.outbound.GetPeers()), ",", len(m.inbound.GetPeers()), ")")
			m.rejectionFilter.AddPeer(candidate.Remote.ID())
			//m.log.Debug("Rejection Filter ", candidate.Remote.ID())
		}

		// signal the result of the outgoing peering request
		Events.OutgoingPeering.Trigger(&PeeringEvent{Self: m.self(), Peer: candidate.Remote, Services: services})
	}
}

func (m *manager) updateInbound(requester *peer.Peer, salt *salt.Salt) {
	// TODO: check request legitimacy
	//m.log.Debug("Evaluating peering request FROM ", requester.ID())
	reqDistance := peer.NewPeerDistance(m.self().Bytes(), m.net.local().GetPrivateSalt().GetBytes(), requester)

	candidateList := []peer.PeerDistance{reqDistance}

	filter := NewFilter()
	filter.AddPeers(m.outbound.GetPeers())      // set filter for outbound neighbors
	filteredList := filter.Apply(candidateList) // filter out current neighbors

	// make decision
	toAccept := m.inbound.Select(filteredList)

	// reject request
	if toAccept.Remote == nil {
		m.log.Debug("Peering request FROM ", requester.ID(), " status REJECTED (", len(m.outbound.GetPeers()), ",", len(m.inbound.GetPeers()), ")")
		m.inboundReplyChan <- reject
		return
	}
	// accept request
	m.inboundReplyChan <- accept
	furthest, _ := m.inbound.getFurthest()
	// drop furthest neighbor
	if furthest.Remote != nil {
		m.inbound.RemovePeer(furthest.Remote.ID())
		m.dropPeer(furthest.Remote)
		m.log.Debug("Inbound furthest removed ", furthest.Remote.ID())
	}
	// update inbound neighborhood
	m.inbound.Add(toAccept)
	m.log.Debug("Peering request FROM ", toAccept.Remote.ID(), " status ACCEPTED (", len(m.outbound.GetPeers()), ",", len(m.inbound.GetPeers()), ")")
}

func (m *manager) updateSalt() (*salt.Salt, *salt.Salt) {
	pubSalt, _ := salt.NewSalt(m.lifetime)
	m.net.local().SetPublicSalt(pubSalt)

	privSalt, _ := salt.NewSalt(m.lifetime)
	m.net.local().SetPrivateSalt(privSalt)

	m.rejectionFilter.Clean()

	if !m.dropNeighbors { // update distance without dropping neighbors
		m.outbound.UpdateDistance(m.self().Bytes(), m.net.local().GetPublicSalt().GetBytes())
		m.inbound.UpdateDistance(m.self().Bytes(), m.net.local().GetPrivateSalt().GetBytes())
	} else { // drop all the neighbors
		m.dropNeighborhood(m.inbound)
		m.dropNeighborhood(m.outbound)
	}

	return pubSalt, privSalt
}

func (m *manager) dropNeighbor(peerToDrop peer.ID) {
	m.inboundDropChan <- peerToDrop
	m.outboundDropChan <- peerToDrop

	// signal the dropped peer
	Events.Dropped.Trigger(&DroppedEvent{Self: m.self(), DroppedID: peerToDrop})
}

// containsPeer returns true if a peer with the given ID is in the list.
func containsPeer(list []*peer.Peer, id peer.ID) bool {
	for _, p := range list {
		if p.ID() == id {
			return true
		}
	}
	return false
}

func (m *manager) acceptRequest(p *peer.Peer, s *salt.Salt, services peer.ServiceMap) peer.ServiceMap {
	m.inboundRequestChan <- peeringRequest{p, s}
	status := <-m.inboundReplyChan
	if status {
		// signal the received request
		Events.IncomingPeering.Trigger(&PeeringEvent{Self: m.self(), Peer: p, Services: services})
	}

	return m.net.local().Services()
}

func (m *manager) getNeighbors() []*peer.Peer {
	var neighbors []*peer.Peer

	neighbors = append(neighbors, m.inbound.GetPeers()...)
	neighbors = append(neighbors, m.outbound.GetPeers()...)

	return neighbors
}

func (m *manager) getIncomingNeighbors() []*peer.Peer {
	var neighbors []*peer.Peer

	neighbors = append(neighbors, m.inbound.GetPeers()...)

	return neighbors
}

func (m *manager) getOutgoingNeighbors() []*peer.Peer {
	var neighbors []*peer.Peer

	neighbors = append(neighbors, m.outbound.GetPeers()...)

	return neighbors
}

func (m *manager) getDuplicates() []peer.ID {
	var d []peer.ID

	for _, p := range m.inbound.GetPeers() {
		if containsPeer(m.outbound.GetPeers(), p.ID()) {
			d = append(d, p.ID())
		}
	}
	return d
}

func (m *manager) dropNeighborhood(nh *Neighborhood) {
	for _, p := range nh.GetPeers() {
		nh.RemovePeer(p.ID())
		m.dropPeer(p)
	}
}

func (m *manager) dropPeer(p *peer.Peer) {
	// send the drop request over the network
	m.net.DropPeer(p)
	// signal the drop
	Events.Dropped.Trigger(&DroppedEvent{Self: m.self(), DroppedID: p.ID()})
}
