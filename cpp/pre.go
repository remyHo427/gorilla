package cpp

import "bytes"

var trigraph_map = map[byte]byte{
	'=':  '#',
	'/':  '\\',
	'\'': '^',
	'(':  '[',
	')':  ']',
	'!':  '|',
	'<':  '{',
	'>':  '}',
	'-':  '~',
}

func pre(input string) string {
	return strip(splice(trigraph(input)))
}

// replace trigraph and splice lines
func trigraph(input string) string {
	var out bytes.Buffer
	l := len(input)

	for i := 0; i < l; i++ {
		c := input[i]

		if c == '?' {
			if i+1 >= l {
				out.WriteByte(c)
				continue
			}

			i++
			nc := input[i]
			if i+1 >= l || nc != '?' {
				out.WriteByte(c)
				out.WriteByte(nc)
				continue
			}

			i++
			nnc := input[i]
			if mapped, ok := trigraph_map[nnc]; ok {
				out.WriteByte(mapped)
				continue
			}

			out.WriteByte(c)
			out.WriteByte(nc)
			out.WriteByte(nnc)
			continue
		}

		out.WriteByte(c)
	}

	return out.String()
}

func splice(input string) string {
	var out bytes.Buffer
	l := len(input)

	for i := 0; i < l; i++ {
		c := input[i]

		if c == '\\' {
			if i+1 >= l {
				out.WriteByte(c)
				continue
			}

			i++
			nc := input[i]
			if nc == '\n' {
				continue
			}

			out.WriteByte(c)
			out.WriteByte(nc)
			continue
		}

		out.WriteByte(c)
	}

	return out.String()
}

func strip(input string) string {
	var out bytes.Buffer
	l := len(input)

	for i := 0; i < l; i++ {
		c := input[i]

		if c == '/' {
			if i+1 >= l {
				out.WriteByte(c)
				continue
			}

			i++
			nc := input[i]
			if nc == '/' {
				i++
				if i+1 >= l {
					out.WriteByte(' ')
					continue
				}

				for {
					i++
					if i+1 >= l {
						out.WriteByte(' ')
						break
					} else if nnc := input[i]; nnc == '\n' {
						out.WriteByte(' ')
						break
					}
				}
				continue
			} else if nc == '*' {
				i++
				if i+1 >= l {
					out.WriteByte(' ')
					continue
				}

				for {
					i++
					if i+1 >= l {
						out.WriteByte(' ')
						break
					}

					nnc := input[i]
					i++

					if i+1 >= l {
						out.WriteByte(' ')
						break
					}

					nnnc := input[i]

					if nnc == '*' && nnnc == '/' {
						out.WriteByte(' ')
						break
					}
				}
				continue
			}

			out.WriteByte(c)
			out.WriteByte(nc)
			continue
		}

		out.WriteByte(c)
	}

	return out.String()
	// var out bytes.Buffer
	// l := len(input)

	// for i := 0; i < l; i++ {
	// 	c := input[i]

	// 	if c == '/' {
	// 		if i+1 >= l {
	// 			out.WriteByte(c)
	// 			continue
	// 		}

	// 		i++
	// 		nc := input[i]
	// 		if nc == '/' {
	// 			i++
	// 			for read := i; read < l; i++ {
	// 				if input[read] == '\n' {
	// 					out.WriteByte(' ')
	// 					break
	// 				}
	// 			}
	// 			out.WriteByte(' ')
	// 			continue
	// 		} else if nc == '*' {
	// 			i++
	// 			for read := i; read < l; i++ {
	// 				if input[read] == '*' && read+1 < l && input[read+1] == '/' {
	// 					i++
	// 					out.WriteByte(' ')
	// 					break
	// 				}
	// 			}
	// 			continue
	// 		}

	// 		out.WriteByte(c)
	// 		out.WriteByte(nc)
	// 	}

	// 	out.WriteByte(c)
	// }

	// return out.String()
}
