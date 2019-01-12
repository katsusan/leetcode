package main

func isIsomorphic(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	m_st := make(map[byte]byte, 26)
	m_ts := make(map[byte]byte, 26)

	for i := 0; i < len(s); i++ {
		v_s, found_s := m_st[s[i]]
		if found_s {
			if v_s != t[i] {
				return false
			}
		} else {
			m_st[s[i]] = t[i]
		}

		v_t, found_t := m_ts[t[i]]
		if found_t {
			if v_t != s[i] {
				return false
			}
		} else {
			m_ts[t[i]] = s[i]
		}
	}

	return true
}
