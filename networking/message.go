package networking

import (
// "common"
// "sync"
)

func (m *Message) IsInTransport(rid ID) bool {
	m.Lock()
	defer m.Unlock()
	return m.recipients[rid].status == PROCESSING
}

func (m *Message) status() MessageStatus {
	// one of three to ocheck whether the message
	// has been attempted to be diistributed
	m.RLock()
	defer m.RUnlock()
	for _, rs := range m.recipients {
		if rs.status == PROCESSING || rs.status == PENDING {
			return rs.status
		}
	}
	return COMPLETED
}

func (m *Message) Delivered(rid ID) MessageStatus {
	m.Lock()
	r := m.recipients[rid]
	r.updateStatus(DELIVERED, nil)
	m.Unlock()
	return m.status()
}

func (m *Message) Failed(rid ID, err error) MessageStatus {
	m.Lock()
	r := m.recipients[rid]
	r.updateStatus(FAILED, err)
	m.Unlock()
	return m.status()
}

func (m *Message) SetInTransport(rid ID) {
	m.Lock()
	defer m.Unlock()
	r := m.recipients[rid]
	r.updateStatus(PROCESSING, nil)
}

func (m *Message) AddRecipient(rid ID) {
	m.Lock()
	defer m.Unlock()
	m.recipients[rid] = Status{status: PENDING}
}

func (m *Message) AddRecipients(ids []ID) {
	m.Lock()
	defer m.Unlock()
	for _, id := range ids {
		m.recipients[id] = Status{status: PENDING}
	}
}

func (m *Message) RecipientIDs() []ID {
	m.RLock()
	defer m.RUnlock()
	ids := make([]ID, 0, len(m.recipients))
	for id := range m.recipients {
		ids = append(ids, id)
	}
	return ids
}

func (m Message) SenderID() ID {
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
