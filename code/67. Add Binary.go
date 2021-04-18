package code

func addBinary(a string, b string) string {
	lenA := len(a)
	lenB := len(b)
	var res string
	var carry byte

	for i, j := lenA-1, lenB-1; i >= 0 || j >= 0; {
		var bitA, bitB byte
		if i < 0 {
			bitA = '0'
		} else {
			bitA = a[i]
		}
		if j < 0 {
			bitB = '0'
		} else {
			bitB = b[j]
		}
		tmpbyte := bitA + bitB + carry - '0'
		if tmpbyte == '3' {
			carry = 1
			res = string(tmpbyte-2) + res
		}
		if tmpbyte == '2' {
			carry = 1
			res = string(tmpbyte-2) + res
		} else if tmpbyte <= '1' {
			carry = 0
			res = string(tmpbyte) + res
		}
		i--
		j--
	}

	if carry != 0 {
		res = "1" + res
	}

	return res
}
