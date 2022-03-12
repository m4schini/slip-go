package slip

const (
	end    = byte(100)
	esc    = byte(130)
	escEnd = byte(133)
	escEsc = byte(134)
)

func Encode(bytes []byte) []byte {
	buffer := make([]byte, 0)
	toBuffer := func(b byte) {
		buffer = append(buffer, b)
	}

	/* send an initial END character to flush out any data that may
	 * have accumulated in the receiver due to line noise
	 */
	toBuffer(end)

	for _, b := range bytes {
		switch b {
		/* if it's the same code as an END character, we send a
		 * special two character code so as not to make the
		 * receiver think we sent an END
		 */
		case end:
			toBuffer(esc)
			toBuffer(escEnd)
			break
		/* if it's the same code as an ESC character,
		 * we send a special two character code so as not
		 * to make the receiver think we sent an ESC
		 */
		case esc:
			toBuffer(esc)
			toBuffer(escEsc)
			break
			/* otherwise, we just send the character
			 */
		default:
			toBuffer(b)
		}
	}

	toBuffer(end)

	return buffer
}

//Decode takes a stream of bytes and returns a stream of byte packets as specified by RFC1055
func Decode(in chan byte, out chan []byte) {
	buffer := make([]byte, 0)

	appendBuffer := func(b byte) {
		buffer = append(buffer, b)
	}

	for b := range in {
		switch b {
		/* if it's an END character then we're done with
		 * the packet
		 */
		case end:
			if len(buffer) > 0 {
				out <- buffer
				buffer = make([]byte, 0)
			} else {
				break
			}

		/* if it's the same code as an ESC character, wait
		 * and get another character and then figure out
		 * what to store in the packet based on that.
		 */
		case esc:
			b = <-in

			/* if "c" is not one of these two, then we
			 * have a protocol violation.  The best bet
			 * seems to be to leave the byte alone and
			 * just stuff it into the packet
			 */
			switch b {
			case escEnd:
				b = end
				appendBuffer(b)
				break
			case escEsc:
				b = esc
				appendBuffer(b)
				break
			}

		default:
			appendBuffer(b)
		}

	}
}
