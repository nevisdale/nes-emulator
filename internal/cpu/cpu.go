package cpu

type ReadWriter interface {
	Read8(addr uint16) uint8
	Read16(addr uint16) uint16
	Write8(addr uint16, data uint8)
}

const (
	flagCBit = uint8(1 << 0) // Carry flag
	flagZBit = uint8(1 << 1) // Zero flag
	flagIBit = uint8(1 << 2) // Interrupt Disable flag
	flagDBit = uint8(1 << 3) // Decimal Mode flag
	flagBBit = uint8(1 << 4) // Break Command flag
	flagUBit = uint8(1 << 5) // Unused
	flagVBit = uint8(1 << 6) // Overflow flag
	flagNBit = uint8(1 << 7) // Negative flag
)

type instruction struct {
	name     string
	operate  opcodeFunc
	addrMode addrMode
	cycles   uint8
}

type CPU struct {
	regA uint8 // used to perform arithmetic and logical operations
	regX uint8 // used primarily for indexing and temporary storage
	regY uint8 // used mainly for indexing and temporary storage

	// The stack is located in the fixed memory page $0100 to $01FF.
	// The Stack Pointer holds the lower 8 bits of the address within this page
	//
	// Initialization:
	// The Stack Pointer is initialized by the system or the program at the start of execution.
	// Typically, it starts at $FF, pointing to the top of the stack.
	sp uint8

	pc     uint16 // program counter
	status uint8  // contains flags from flagXBit

	// bus to connect to RAM
	bus ReadWriter

	// Opcode matrix. see more https://www.masswerk.at/6502/6502_instruction_set.html
	//
	// Position in the slice is opcode.
	instructions []instruction

	fetched      uint8
	addrAbs      uint16
	addrRel      uint16
	opcode       uint8
	cycles       uint8
	clockCounter uint64
}

func NewCPU() (*CPU, error) {
	c := &CPU{
		sp:           0xff,
		instructions: make([]instruction, 0x100),
	}
	if err := c.parseOpcodeMatrix(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CPU) ConnectBus(bus ReadWriter) {
	c.bus = bus
}

func (c CPU) getFlag(flag uint8) bool {
	return c.status&flag > 0
}

func (c *CPU) setFlag(flag uint8, v bool) {
	if v {
		c.status |= flag
		return
	}
	c.status &= ^flag
}

func (c *CPU) Tic() {
	if c.cycles != 0 {
		c.cycles--
		return
	}

	c.opcode = c.bus.Read8(c.pc)
	c.pc++
	inst := c.instructions[c.opcode]
	cycleCount1 := c.doAddressMode(inst.addrMode)
	cycleCount2 := inst.operate()

	c.cycles = inst.cycles + cycleCount1 + cycleCount2
}

// TODO: may merge all Reset, IRQ, NMI into one function?
//
// reset the CPU to its initial state
func (c *CPU) Reset() {
}

// interrupt request signal
func (c *CPU) IRQ() {}

// non-maskable interrupt request signal
func (c *CPU) NMI() {}

func (c *CPU) fetch() uint8 {
	return 0
}
