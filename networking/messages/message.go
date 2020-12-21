package messages

import (
	// "common"
	// "sync"
	ln "latticesphere/networking"
)

func (m *Message) IsInTransport(rid ln.ID) bool {
	m.Lock()
	defer m.Unlock()
	return m.recipients[rid].status == ln.PROCESSING
}

func (m *Message) status() ln.MessageStatus {
	// one of three to ocheck whether the message
	// has been attempted to be diistributed
	m.RLock()
	defer m.RUnlock()
	for _, rs := range m.recipients {
		if rs.status == ln.PROCESSING || rs.status == ln.PENDING {
			return rs.status
		}
	}
	return ln.COMPLETED
}

func (m *Message) Delivered(rid ln.ID) ln.MessageStatus {
	m.Lock()
	r := m.recipients[rid]
	r.updateStatus(ln.DELIVERED, nil)
	m.Unlock()
	return m.status()
}

func (m *Message) Failed(rid ln.ID, err error) ln.MessageStatus {
	m.Lock()
	r := m.recipients[rid]
	r.updateStatus(ln.FAILED, err)
	m.Unlock()
	return m.status()
}

func (m *Message) SetInTransport(rid ln.ID) {
	m.Lock()
	defer m.Unlock()
	r := m.recipients[rid]
	r.updateStatus(ln.PROCESSING, nil)
}

func (m *Message) AddRecipient(rid ln.ID) {
	m.Lock()
	defer m.Unlock()
	m.recipients[rid] = Status{status: ln.PENDING}
}

func (m *Message) AddRecipients(ids []ln.ID) {
	m.Lock()
	defer m.Unlock()
	for _, id := range ids {
		m.recipients[id] = Status{status: ln.PENDING}
	}
}

func (m *Message) RecipientIDs() []ln.ID {
	m.RLock()
	defer m.RUnlock()
	ids := make([]ln.ID, 0, len(m.recipients))
	for id := range m.recipients {
		ids = append(ids, id)
	}
	return ids
}

func (m Message) SenderID() ln.ID {
	m.RLock()
	defer m.RUnlock()
	return m.senderID
}

func (m *Message) String() string {
	m.RLock()
	defer m.RUnlock()
	return m.content
}

func (m *Message) Bytes() []byte {
	return []byte(m.String())
}
