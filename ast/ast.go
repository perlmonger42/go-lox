package ast

type Node interface {
	AsNode() Node // does nothing but prevent non-Nodes from looking like a Node
}
