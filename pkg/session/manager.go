package session

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Manager struct {
	sessions       map[string]*TimeoutSession
	kickOffChan    chan string
	sessionMax     int64
	sessionTimeout int64
	sessionLife    int64
	sync.Mutex
}

func New(max, timeout, life int64) *Manager {

	sm := &Manager{
		sessions:       make(map[string]*TimeoutSession),
		kickOffChan:    make(chan string),
		sessionMax:     max,
		sessionTimeout: timeout,
		sessionLife:    life,
	}
	return sm
}

func (m *Manager) kickoff(sid string) {
	m.Lock()
	defer m.Unlock()
	delete(m.sessions, sid)
	log.Println("kickoff session:", sid)
}

func (m *Manager) CheckSessions() {
	for {
		time.Sleep(time.Second * 2)
		for sid, ts := range m.sessions {
			if ts.Alive(m.sessionTimeout, m.sessionLife) {
				m.kickoff(sid)
			}
		}
	}
}

func (m *Manager) CheckKickoffSignal() {
	for {
		sid := <-m.kickOffChan
		log.Println("kickoff signal sid: ", sid)
		m.kickoff(sid)
	}
}

func (m *Manager) GetSession(sid string) *Session {
	ts, ok := m.sessions[sid]
	if !ok {
		log.Println("Error: Session not exist ", sid)
		return nil
	}

	ts.Update()

	return ts.session
}

func (m *Manager) Has(sid string) bool {
	_, ok := m.sessions[sid]
	return ok
}

func (m *Manager) Add(s *Session) {
	m.Lock()
	defer m.Unlock()
	nano := time.Now().UnixNano()
	s.Id = fmt.Sprint(nano)
	s.Id = "123"
	now := int64(nano / 1e9)
	m.sessions[s.Id] = &TimeoutSession{s, now, now}
}

func (m *Manager) Print() {
	fmt.Println("sessions", m.sessions)
}
