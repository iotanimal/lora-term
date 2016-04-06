# LoRa Terminal

This is a simple terminal interface to write/read command from Microchip RN2903. This include and simple dump util and sub channel configuration util.

## Usage:

  Usage of ./lora-term:
    -baud uint
      	Baud rate (default 57600)
    -databits uint
      	Data bits (default 8)
    -minread uint
      	Minimum read count (default 1)
    -port string
      	serial port to test (/dev/tty.usbmodem1421, /dev/ttyUSB0) (default "/dev/tty.usbmodem1421")
    -stopbits uint
      	Stop bits (default 1)
    -sub-band int
      	Set specific sub band channels on (default -1)
    -term
      	Terminal Emu mode


## Example

  ./termgo
  Version             : invalid_param
  Dev Addr            : 001A5A5E
  Dev EUI             : 0004A30B001A5A5E
  APP EUI             : 0000000000000000
  ---------------------------------------------------------
    computing sub channels \
    SUB CHAN |   0   |   1   |   2   |   3   |   4   |   5   |   6   |   7   | 500MHZ
  +----------+-------+-------+-------+-------+-------+-------+-------+-------+--------+
           1 | off   | off   | off   | off   | off   | off   | off   | off   | off
           2 | on    | on    | on    | on    | on    | on    | on    | on    | on
           3 | off   | off   | off   | off   | off   | off   | off   | off   | off
           4 | off   | off   | off   | off   | off   | off   | off   | off   | off
           5 | off   | off   | off   | off   | off   | off   | off   | off   | off
           6 | off   | off   | off   | off   | off   | off   | off   | off   | off
           7 | off   | off   | off   | off   | off   | off   | off   | off   | off
           8 | off   | off   | off   | off   | off   | off   | off   | off   | off

