package json

import (
	"unsafe"
)

type encodeCompileContext struct {
	typ                      *rtype
	root                     bool
	opcodeIndex              int
	ptrIndex                 int
	indent                   int
	structTypeToCompiledCode map[uintptr]*compiledCode

	parent *encodeCompileContext
}

func (c *encodeCompileContext) context() *encodeCompileContext {
	return &encodeCompileContext{
		typ:                      c.typ,
		root:                     c.root,
		opcodeIndex:              c.opcodeIndex,
		ptrIndex:                 c.ptrIndex,
		indent:                   c.indent,
		structTypeToCompiledCode: c.structTypeToCompiledCode,
		parent:                   c,
	}
}

func (c *encodeCompileContext) withType(typ *rtype) *encodeCompileContext {
	ctx := c.context()
	ctx.typ = typ
	return ctx
}

func (c *encodeCompileContext) incIndent() *encodeCompileContext {
	ctx := c.context()
	ctx.indent++
	return ctx
}

func (c *encodeCompileContext) decIndent() *encodeCompileContext {
	ctx := c.context()
	ctx.indent--
	return ctx
}

func (c *encodeCompileContext) incIndex() {
	c.incOpcodeIndex()
	c.incPtrIndex()
}

func (c *encodeCompileContext) decIndex() {
	c.decOpcodeIndex()
	c.decPtrIndex()
}

func (c *encodeCompileContext) incOpcodeIndex() {
	c.opcodeIndex++
	if c.parent != nil {
		c.parent.incOpcodeIndex()
	}
}

func (c *encodeCompileContext) decOpcodeIndex() {
	c.opcodeIndex--
	if c.parent != nil {
		c.parent.decOpcodeIndex()
	}
}

func (c *encodeCompileContext) incPtrIndex() {
	c.ptrIndex++
	if c.parent != nil {
		c.parent.incPtrIndex()
	}
}

func (c *encodeCompileContext) decPtrIndex() {
	c.ptrIndex--
	if c.parent != nil {
		c.parent.decPtrIndex()
	}
}

type encodeRuntimeContext struct {
	ptrs     []uintptr
	keepRefs []unsafe.Pointer
}

func (c *encodeRuntimeContext) init(p uintptr, codelen int) {
	if len(c.ptrs) < codelen {
		c.ptrs = make([]uintptr, codelen)
	}
	c.ptrs[0] = p
	c.keepRefs = c.keepRefs[:0]
}

func (c *encodeRuntimeContext) ptr() uintptr {
	header := (*sliceHeader)(unsafe.Pointer(&c.ptrs))
	return uintptr(header.data)
}
