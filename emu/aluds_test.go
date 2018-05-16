package emu

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/yaotsu/gcn3/insts"
)

var _ = Describe("ALU", func() {

	var (
		alu   *ALUImpl
		state *mockInstState
	)

	BeforeEach(func() {
		alu = NewALUImpl(nil)
		alu.lds = make([]byte, 4096)

		state = new(mockInstState)
		state.scratchpad = make([]byte, 4096)
	})

	It("should run DS_WRITE2_B64", func() {
		state.inst = insts.NewInst()
		state.inst.FormatType = insts.DS
		state.inst.Opcode = 78
		state.inst.Offset0 = 1
		state.inst.Offset1 = 3

		sp := state.scratchpad.AsDS()
		sp.EXEC = 0x1
		sp.ADDR[0] = 100
		sp.DATA[0] = 1
		sp.DATA[1] = 2
		sp.DATA1[0] = 3
		sp.DATA1[1] = 4

		alu.Run(state)

		lds := alu.LDS()
		Expect(insts.BytesToUint32(lds[108:])).To(Equal(uint32(1)))
		Expect(insts.BytesToUint32(lds[112:])).To(Equal(uint32(2)))
		Expect(insts.BytesToUint32(lds[124:])).To(Equal(uint32(3)))
		Expect(insts.BytesToUint32(lds[128:])).To(Equal(uint32(4)))
	})

	It("should run DS_READ2_B64", func() {
		state.inst = insts.NewInst()
		state.inst.FormatType = insts.DS
		state.inst.Opcode = 119
		state.inst.Offset0 = 1
		state.inst.Offset1 = 3

		sp := state.scratchpad.AsDS()
		sp.EXEC = 0x1
		sp.ADDR[0] = 100

		lds := alu.LDS()
		copy(lds[108:], insts.Uint32ToBytes(12))
		copy(lds[124:], insts.Uint32ToBytes(156))

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint32(12)))
		Expect(sp.DST[2]).To(Equal(uint32(156)))
	})

})