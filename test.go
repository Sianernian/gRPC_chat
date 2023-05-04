package main

import "fmt"

// 定义链表节点
type PolyNode struct {
	coef int       // 系数
	exp  int       // 指数
	next *PolyNode // 下一个节点指针
}

// 定义链表
type PolyList struct {
	head *PolyNode // 头节点指针
	tail *PolyNode // 尾节点指针
}

// 定义链表的加法运算方法
func (pl *PolyList) Add(poly *PolyList) {
	p1 := pl.head
	p2 := poly.head
	p3 := &PolyNode{}
	pl.head = p3
	for p1 != nil && p2 != nil {
		if p1.exp > p2.exp {
			p3.next = p1
			p1 = p1.next
		} else if p1.exp == p2.exp {
			coef := p1.coef + p2.coef
			if coef != 0 {
				p3.next = &PolyNode{coef: coef, exp: p1.exp}
				p3 = p3.next
			}
			p1 = p1.next
			p2 = p2.next
		} else {
			p3.next = p2
			p2 = p2.next
		}
		p3 = p3.next
	}
	if p1 != nil {
		p3.next = p1
	}
	if p2 != nil {
		p3.next = p2
	}
	pl.head = pl.head.next
}

// 定义链表的乘法运算方法
func (pl *PolyList) Multiply(poly *PolyList) {
	p1 := pl.head
	p2 := poly.head
	p3 := &PolyNode{}
	pl.head = p3
	for p1 != nil {
		for p2 != nil {
			coef := p1.coef * p2.coef
			exp := p1.exp + p2.exp
			p := &PolyNode{coef: coef, exp: exp}
			p3.next = p
			p3 = p3.next
			p2 = p2.next
		}
		p2 = poly.head
		p1 = p1.next
	}
	pl.head = pl.head.next
}

// 定义链表的插入方法
func (pl *PolyList) Insert(coef int, exp int) {
	node := &PolyNode{coef: coef, exp: exp}
	if pl.head == nil {
		pl.head = node
		pl.tail = node
	} else {
		pl.tail.next = node
		pl.tail = node
	}
}

// 定义链表的输出方法
func (pl *PolyList) Print() {
	p := pl.head
	for p != nil {
		fmt.Printf("%dX^%d", p.coef, p.exp)
		if p.next != nil {
			fmt.Print("+")
		}
		p = p.next
	}
	fmt.Println()
}
